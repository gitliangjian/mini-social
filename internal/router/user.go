package router

import (
	"mini-social/internal/config"
	"mini-social/internal/handler"
	"mini-social/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, cfg *config.Config, userHandler *handler.UserHandler) {
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		//访问前必须先通过JWT中间件
		users.Use(middleware.JWTAuth(cfg))
		{
			users.GET("/me", userHandler.Me)
		}
	}
}
