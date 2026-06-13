package functions

import (
	"strconv"
	"time"

	"findsalon-backend/integrations"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId, email, role string) (string, error) {
	expiryHours, err := strconv.Atoi(integrations.JwtExpiryHours)
	if err != nil {
		expiryHours = 24
	}
	claims := JWTClaims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(integrations.JwtSecret))
}

func GenerateRefreshToken(userId string) (string, error) {
	expiryDays, err := strconv.Atoi(integrations.RefreshTokenExpiryDays)
	if err != nil {
		expiryDays = 7
	}
	claims := JWTClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(integrations.JwtSecret))
}
