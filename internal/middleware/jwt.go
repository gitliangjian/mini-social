package middleware

import (
	"mini-social/internal/config"
	jwtutil "mini-social/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "authorization header is required",
				"data": nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "invalid authorization header",
				"data": nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtutil.ParseToken(cfg.JWT.Secret, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "invalid or expired token",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Next()
	}
}
