package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	AppEnv     string
	AppPort    string
	DBHost     string
	DBPassword string
	DBPort     string
	DBUser     string
	DBName     string
	DBCharset  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Computed TTLs (in minutes for access, days for refresh)
	AccessTTLMin  int
	RefreshTTLDay int

	UploadDir       string
	MaxUploadSizeMB int64

	// Scanner service
	PDFStorageDir string
	ScannerURL    string
}


var App *Config

func Load() error {
	//Load .env file if present (ingored in production)
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	maxUpload, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)
	gin.SetMode(gin.ReleaseMode)

	App = &Config{
		AppName: getEnv("APP_NAME", "The_Scan"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8008"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "psn"),
		DBPassword: getEnv("DB_PASSWORD", "160125Neath"),
		DBName:     getEnv("DB_NAME", "the_scan"),
		DBCharset:  getEnv("DB_CHARSET", "utf8mb4"),

		JWTSecret:      getEnv("JWT_SECRET", "Scc_secret"),
		JWTExpiryHours: jwtExpiry,

		// Compute TTLs from JWTExpiryHours
		// Access token TTL = JWTExpiryHours (in minutes)
		// Refresh token TTL = JWTExpiryHours * 30 (in days, default 30 days per refresh)
		AccessTTLMin:  jwtExpiry * 60,
		RefreshTTLDay: jwtExpiry / 24,

		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB: maxUpload,

		// Scanner service config
		PDFStorageDir: getEnv("PDF_STORAGE_DIR", "./pdfs"),
		ScannerURL:    getEnv("SCANNER_URL", "http://localhost:8008"),
	}
	return nil

}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// getEnv returns the value of an environment variable.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
