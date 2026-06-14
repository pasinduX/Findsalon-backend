package integrations

import (
	"log"
	"os"
	"strconv"
	"strings"

	"findsalon-backend/dbConfig"
)

// Auth
var JwtSecret string
var JwtExpiryHours string
var RefreshTokenExpiryDays string
var Auth0Domain string
var Auth0Audience string
var Auth0ClientId string
var Auth0RolesClaim string
var GoogleClientId string
var GoogleClientSecret string
var GoogleRedirectUrl string
var GoogleAuthStateSecret string
var FrontendUrl string

// Server
var ServerPort string
var ServiceName string

// Database
var DatabaseUrl string
var DatabaseName string

// Storage (local + S3)
var StorageBasePath string
var MaxImageSizeMb int
var AwsRegion string
var AwsAccessKeyId string
var AwsSecretAccessKey string
var AwsS3Bucket string
var AwsS3BaseUrl string

// SMTP
var SmtpHost string
var SmtpPort int
var SmtpUsername string
var SmtpPassword string
var SmtpFromName string
var SmtpFromEmail string
var SmtpEnabled bool

// Pagination
var DefaultPageSize int

func SetEnvironmentVariables() {
	DatabaseUrl = getEnv("DATABASE_URL", "mongodb://localhost:27017")
	DatabaseName = getEnv("DATABASE_NAME", "findsalon_db")
	ServerPort = getEnv("PORT", getEnv("SERVER_PORT", "8888"))
	ServiceName = getEnv("SERVICE_NAME", "findsalon-backend")

	JwtSecret = getEnv("JWT_SECRET", "your_super_secret_jwt_key_here")
	JwtExpiryHours = getEnv("JWT_EXPIRY_HOURS", "24")
	RefreshTokenExpiryDays = getEnv("REFRESH_TOKEN_EXPIRY_DAYS", "7")
	Auth0Domain = getEnv("AUTH0_DOMAIN", "")
	Auth0Audience = getEnv("AUTH0_AUDIENCE", "")
	Auth0ClientId = getEnv("AUTH0_CLIENT_ID", "")
	Auth0RolesClaim = getEnv("AUTH0_ROLES_CLAIM", "")

	GoogleClientId = getEnv("GOOGLE_CLIENT_ID", "")
	GoogleClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
	GoogleRedirectUrl = getEnv("GOOGLE_REDIRECT_URL", "")
	GoogleAuthStateSecret = getEnv("GOOGLE_AUTH_STATE_SECRET", "")
	FrontendUrl = getEnv("FRONTEND_URL", "http://localhost:4000")

	StorageBasePath = getEnv("STORAGE_BASE_PATH", "./uploads")
	MaxImageSizeMb = getEnvAsInt("MAX_IMAGE_SIZE_MB", 3)

	AwsRegion = getEnv("AWS_REGION", "")
	AwsAccessKeyId = getEnv("AWS_ACCESS_KEY_ID", "")
	AwsSecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")
	AwsS3Bucket = getEnv("AWS_S3_BUCKET", "")
	AwsS3BaseUrl = getEnv("AWS_S3_BASE_URL", "")

	SmtpHost = getEnv("SMTP_HOST", "")
	SmtpPort = getEnvAsInt("SMTP_PORT", 587)
	SmtpUsername = getEnv("SMTP_USERNAME", "")
	SmtpPassword = getEnv("SMTP_PASSWORD", "")
	SmtpFromName = getEnv("SMTP_FROM_NAME", "FindSalon")
	SmtpFromEmail = getEnv("SMTP_FROM_EMAIL", "noreply@findsalon.lk")
	SmtpEnabled = getEnv("SMTP_ENABLED", "false") == "true"

	DefaultPageSize = getEnvAsInt("DEFAULT_PAGE_SIZE", 20)

	if DatabaseUrl == "" || DatabaseName == "" {
		log.Fatal("DATABASE_URL and DATABASE_NAME must be set")
	}
	if !isMongoURI(DatabaseUrl) {
		log.Fatal("DATABASE_URL must start with mongodb:// or mongodb+srv://")
	}

	dbConfig.DATABASE_URL = DatabaseUrl
	dbConfig.DATABASE_NAME = DatabaseName
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func isMongoURI(value string) bool {
	uri := strings.TrimSpace(value)
	return strings.HasPrefix(uri, "mongodb://") || strings.HasPrefix(uri, "mongodb+srv://")
}
