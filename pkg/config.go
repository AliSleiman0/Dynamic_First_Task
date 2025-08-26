package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	JWTTTL    time.Duration
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	DBCharset string
	AppPort   string
	Env       string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env") //without this env fields get empty
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")

	ttl := 72 * time.Hour // ttl yaane time to live
	if v := os.Getenv("JWT_TTL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ttl = time.Duration(n) * time.Hour
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbCharset := os.Getenv("DB_CHARSET")
	if dbCharset == "" {
		dbCharset = "utf8mb4"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000"
	}
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		JWTSecret: secret,
		JWTTTL:    ttl,
		DBUser:    dbUser,
		DBPass:    dbPass,
		DBHost:    dbHost,
		DBPort:    dbPort,
		DBName:    dbName,
		DBCharset: dbCharset,
		AppPort:   appPort,
		Env:       env,
	}
}

// GenerateTokens generates a new access and refresh JWT token for a user.
func GenerateTokens(userID int, email string) (string, string, error) {
	accessSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" {
		accessSecret = "default_access_secret"
	}
	if refreshSecret == "" {
		refreshSecret = "default_refresh_secret"
	}
	ttl := 72 * time.Hour // ttl yaane time to live
	if v := os.Getenv("JWT_TTL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ttl = time.Duration(n) * time.Hour
		}
	}
	// Access token
	//NewWithClaims does the following:
	// 	Creates JWT object token li ha hteja bl auth.
	//  Attaches ll claims (payload data) ana bde hutta bl token.
	//  Prepares it to be signed later with a secret key.
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   strconv.Itoa(int(userID)), //"sub" usaully mnhutta user ID
		"email": email,
		"exp":   time.Now().Add(ttl).Unix(), //when does it expire
		"iat":   time.Now().Unix(),          // when it was created
	})
	//singedstring does the following:
	//Signs the JWT with a secret key or private key.

	//Generates the final JWT string that you can send to clients.

	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(userID)),
		"exp": time.Now().Add(ttl * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
