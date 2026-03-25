package router

import (
	"mini-social/internal/config"
	"mini-social/internal/handler"
	"mini-social/internal/middleware"

	"github.com/gin-gonic/gin"
)

// comment/like是moment的子资源，所以将相关路由也放在这里
func RegisterMomentRoutes(r *gin.Engine, cfg *config.Config, momentHandler *handler.MomentHandler, commentHandler *handler.CommentHandler, likeHandler *handler.LikeHandler) {
	api := r.Group("/api/v1")
	{
		moments := api.Group("/moments")
		{
			moments.GET("", momentHandler.List)
			moments.GET("/:moment_id", momentHandler.Detail)
			moments.GET("/:moment_id/comments", commentHandler.List)
		}

		authMoments := api.Group("/moments")
		authMoments.Use(middleware.JWTAuth(cfg))
		{
			//动态相关路由
			authMoments.POST("", momentHandler.Create)
			authMoments.DELETE("/:moment_id", momentHandler.Delete)

			//评论相关路由
			comments := authMoments.Group("/:moment_id/comments")
			{
				comments.POST("", commentHandler.Create)
				comments.DELETE("/:comment_id", commentHandler.Delete)
			}

			// 点赞动态： POST /api/v1/moments/:id/like
			authMoments.POST("/:moment_id/like", likeHandler.Like) // target_type = "moment", target_id = :id

			// 点赞评论： POST /api/v1/moments/:id/comments/:comment_id/like
			authMoments.POST("/:moment_id/comments/:comment_id/like", likeHandler.Like)

			//取消点赞
			// 取消点赞同理
			authMoments.DELETE("/:moment_id/like", likeHandler.UnLike)
			authMoments.DELETE("/:moment_id/comments/:comment_id/like", likeHandler.UnLike)
		}
	}
}
