package middleware

import (
	"os"
	"time"

	"github.com/Napat/mcpserver-demo/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTClaims structure for storing data in JWT
type JWTClaims struct {
	UserID uint            `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for the user
func GenerateToken(userID uint, role models.UserRole) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your_jwt_secret_key_here"
	}

	expirationStr := os.Getenv("JWT_EXPIRATION")
	if expirationStr == "" {
		expirationStr = "24h"
	}

	expiration, err := time.ParseDuration(expirationStr)
	if err != nil {
		expiration = 24 * time.Hour
	}

	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserIDFromToken extracts UserID from token
func GetUserIDFromToken(c echo.Context) uint {
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return 0
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0
	}

	return uint(userID)
}

// GetUserRoleFromToken extracts UserRole from token
func GetUserRoleFromToken(c echo.Context) models.UserRole {
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return models.RoleUser
	}

	role, ok := claims["role"].(float64)
	if !ok {
		return models.RoleUser
	}

	return models.UserRole(uint8(role))
}

// AdminMiddleware middleware for checking if user is an Admin
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, ok := c.Get("user").(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(403, "Forbidden: Authentication required")
		}

		roleValue, ok := claims["role"].(float64)
		if !ok {
			return echo.NewHTTPError(403, "Forbidden: Invalid role information")
		}

		role := models.UserRole(uint8(roleValue))
		if !role.HasRole(models.RoleAdmin) && !role.HasRole(models.RoleSuperAdmin) {
			return echo.NewHTTPError(403, "Forbidden: Admin role required")
		}

		return next(c)
	}
}

// SuperAdminMiddleware middleware for checking if user is a SuperAdmin
func SuperAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, ok := c.Get("user").(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(403, "Forbidden: Authentication required")
		}

		roleValue, ok := claims["role"].(float64)
		if !ok {
			return echo.NewHTTPError(403, "Forbidden: Invalid role information")
		}

		role := models.UserRole(uint8(roleValue))
		if !role.HasRole(models.RoleSuperAdmin) {
			return echo.NewHTTPError(403, "Forbidden: SuperAdmin role required")
		}

		return next(c)
	}
}
