package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	servicedto "vesuvio/internal/dto/service"
	"vesuvio/internal/service"
)

const ContextUserKey = "currentUser"

// AuthMiddleware loads the user from header X-User-ID.
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDHeader := c.GetHeader("X-User-ID")
		if userIDHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-ID"})
			return
		}

		uid, err := strconv.ParseUint(userIDHeader, 10, 64)
		if err != nil || uid == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid X-User-ID"})
			return
		}

		user, err := authService.GetUserByID(c.Request.Context(), uint(uid))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set(ContextUserKey, *user)
		c.Next()
	}
}

// AdminOnly ensures current user is admin.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(ContextUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user := val.(servicedto.User)
		if !user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
