package middleware

import (
	"mini-social/internal/config"
	jwtutil "mini-social/pkg/jwt"
	"mini-social/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			response.Unauthorized(c, "invalid authorization header")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtutil.ParseToken(cfg.JWT.Secret, tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Next()
	}
}
