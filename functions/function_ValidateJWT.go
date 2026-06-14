package functions

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/integrations"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserId      string   `json:"UserId"`
	Email       string   `json:"email"`
	Role        string   `json:"Role"`
	Permissions []string `json:"permissions"`
	Scope       string   `json:"scope"`
	jwt.RegisteredClaims
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var jwksCache = struct {
	sync.RWMutex
	keys      map[string]*rsa.PublicKey
	expiresAt time.Time
}{keys: map[string]*rsa.PublicKey{}}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if integrations.Auth0Domain == "" {
		return nil, errors.New("AUTH0_DOMAIN is required")
	}

	claims := &JWTClaims{}
	parserOptions := []jwt.ParserOption{
		jwt.WithIssuer(auth0Issuer()),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	}
	if integrations.Auth0Audience != "" {
		parserOptions = append(parserOptions, jwt.WithAudience(integrations.Auth0Audience))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, auth0KeyFunc, parserOptions...)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fiber.ErrUnauthorized
	}
	if claims.UserId == "" {
		claims.UserId = claims.Subject
	}
	if claims.Role == "" {
		claims.Role = claimRole(tokenString, integrations.Auth0RolesClaim)
	}
	if claims.Role == "" {
		claims.Role = dto.RoleUser
	}

	return claims, nil
}

func Auth0Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Authorization header is required")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		claims, err := ValidateJWT(parts[1])
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		userId := claims.UserId
		if userId == claims.Subject {
			userId = ""
		}
		email := claims.Email
		role := claims.Role

		if email != "" {
			user, err := dao.FindUserByEmail(email)
			if err != nil {
				// User not in DB yet — provision them from Auth0 claims
				now := time.Now()
				newUser := dto.User{
					UserId:    uuid.New().String(),
					FullName:  email,
					Email:     email,
					GoogleId:  claims.Subject,
					Provider:  "auth0",
					Role:      dto.RoleUser,
					IsActive:  true,
					CreatedAt: now,
					UpdatedAt: now,
				}
				if createErr := dao.CreateUser(newUser); createErr != nil {
					log.Printf("Auth0Middleware: failed to create user %s: %v", email, createErr)
				} else {
					userId = newUser.UserId
					role = newUser.Role
				}
			} else {
				userId = user.UserId
				role = user.Role
			}
		}

		if userId == "" && claims.Subject != "" {
			if user, err := dao.FindUserByGoogleId(claims.Subject); err == nil {
				userId = user.UserId
				role = user.Role
				if email == "" {
					email = user.Email
				}
			}
		}

		if userId == "" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Authenticated user is not synced")
		}

		c.Locals("userId", userId)
		c.Locals("email", email)
		c.Locals("role", role)
		c.Locals("auth0Sub", claims.Subject)
		c.Locals("permissions", claims.Permissions)
		return c.Next()
	}
}

func JWTMiddleware() fiber.Handler {
	return Auth0Middleware()
}

func auth0KeyFunc(token *jwt.Token) (interface{}, error) {
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return nil, errors.New("missing token kid")
	}

	if key := cachedJWK(kid); key != nil {
		return key, nil
	}
	if err := refreshJWKS(); err != nil {
		return nil, err
	}
	if key := cachedJWK(kid); key != nil {
		return key, nil
	}
	return nil, errors.New("matching Auth0 signing key not found")
}

func cachedJWK(kid string) *rsa.PublicKey {
	jwksCache.RLock()
	defer jwksCache.RUnlock()
	if time.Now().After(jwksCache.expiresAt) {
		return nil
	}
	return jwksCache.keys[kid]
}

func refreshJWKS() error {
	jwksCache.Lock()
	defer jwksCache.Unlock()

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(auth0Issuer() + ".well-known/jwks.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("failed to load Auth0 JWKS")
	}

	var payload jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey, len(payload.Keys))
	for _, key := range payload.Keys {
		if key.Kty != "RSA" || key.N == "" || key.E == "" {
			continue
		}
		publicKey, err := rsaPublicKey(key)
		if err != nil {
			continue
		}
		keys[key.Kid] = publicKey
	}
	jwksCache.keys = keys
	jwksCache.expiresAt = time.Now().Add(10 * time.Minute)
	return nil
}

func rsaPublicKey(key jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 + int(b)
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: exponent,
	}, nil
}

func auth0Issuer() string {
	domain := strings.TrimSpace(integrations.Auth0Domain)
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	return "https://" + strings.TrimSuffix(domain, "/") + "/"
}

func claimRole(tokenString, rolesClaim string) string {
	mapClaims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenString, mapClaims)
	if err != nil {
		return ""
	}
	candidateKeys := []string{"Role", "role", "roles"}
	if rolesClaim != "" {
		candidateKeys = append([]string{rolesClaim}, candidateKeys...)
	}
	for _, key := range candidateKeys {
		switch value := mapClaims[key].(type) {
		case string:
			return value
		case []interface{}:
			if len(value) > 0 {
				if role, ok := value[0].(string); ok {
					return role
				}
			}
		case []string:
			if len(value) > 0 {
				return value[0]
			}
		}
	}
	return ""
}
