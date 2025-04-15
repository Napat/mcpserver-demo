package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Napat/mcpserver-demo/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTConfig is the configuration for JWT middleware
type JWTConfig struct {
	Secret string
}

// getJWTSecret retrieves the secret key from the environment
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your_jwt_secret_key_here" // Default value from .env
	}
	return secret
}

// JWTMiddleware checks JWT token
func JWTMiddleware() echo.MiddlewareFunc {
	secret := getJWTSecret()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Split token from "Bearer "
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header format must be Bearer {token}",
				})
			}

			tokenString := parts[1]

			// Check token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// Retrieve claims from token
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("user", claims)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims",
			})
		}
	}
}

// RoleMiddleware verifies user permissions
func RoleMiddleware(requiredRole models.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User not authenticated",
				})
			}

			// Extract role from claims
			userRole, ok := userClaims["role"].(float64)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Invalid role information",
				})
			}

			// Check if user has the necessary permissions
			if models.UserRole(uint8(userRole))&requiredRole == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Insufficient permissions",
				})
			}

			return next(c)
		}
	}
}
