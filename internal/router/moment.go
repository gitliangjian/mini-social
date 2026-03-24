package router

import (
	"mini-social/internal/config"
	"mini-social/internal/handler"
	"mini-social/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterMomentRoutes(r *gin.Engine, cfg *config.Config, momentHandler *handler.MomentHandler) {
	api := r.Group("/api/v1")
	{
		moments := api.Group("/moments")
		{
			moments.GET("", momentHandler.List)
			moments.GET("/:id", momentHandler.Detail)
		}

		authMoments := api.Group("/moments")
		authMoments.Use(middleware.JWTAuth(cfg))
		{
			authMoments.POST("", momentHandler.Create)
			authMoments.DELETE("/:id", momentHandler.Delete)
		}
	}
}
