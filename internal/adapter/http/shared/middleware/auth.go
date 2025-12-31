package middleware

import (
	"net/http"
	"nexus/internal/adapter/http/shared/response"
	jwtpkg "nexus/pkg/jwt"
	"nexus/pkg/uuidv7"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	authorizationPrefix = "Bearer "
	userIDKey           = "user_id"
	userEmailKey        = "user_email"
)

type AuthMiddleware struct {
	jwtManager *jwtpkg.JWTManager
}

func NewAuthMiddleware(jwtManager *jwtpkg.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Requires valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header required", nil)
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid or expired token", err)
			c.Abort()
			return
		}

		// Save user data to context
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.Email)

		c.Next()
	}
}

// Tries to extract token but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err == nil {
			c.Set(userIDKey, claims.UserID)
			c.Set(userEmailKey, claims.Email)
		}

		c.Next()
	}
}

func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader(authorizationHeader)
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, authorizationPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, authorizationPrefix)
}

// Helper functions to get data from the context
func GetUserID(c *gin.Context) (uuidv7.UUID, bool) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return uuidv7.Nil, false
	}

	userID, ok := value.(uuidv7.UUID)
	return userID, ok
}

func GetUserEmail(c *gin.Context) (string, bool) {
	value, exists := c.Get(userEmailKey)
	if !exists {
		return "", false
	}

	email, ok := value.(string)
	return email, ok
}

// gets UserID or panics (for protected routes)
func MustGetUserID(c *gin.Context) uuidv7.UUID {
	userID, ok := GetUserID(c)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}
