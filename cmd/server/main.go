package main

import (
	"log"
	"net/http"

	"mini-social/internal/bootstrap"
	"mini-social/internal/config"
	"mini-social/internal/handler"
	"mini-social/internal/repository"
	"mini-social/internal/router"
	"mini-social/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	//读取配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed:%v", err)
	}

	//数据库初始化和迁移
	db, err := bootstrap.InitDB(cfg)
	if err != nil {
		log.Fatalf("init db failed:%v", err)
	}
	if err := bootstrap.AutoMigrate(db); err != nil {
		log.Fatalf("migrate db failed:%v", err)
	}

	//用户
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	log.Println("jwt secret:", cfg.JWT.Secret)
	userHandler := handler.NewUserHandler(userService)

	//动态
	momentRepo := repository.NewMomentRepository(db)
	momentService := service.NewMomentService(momentRepo)
	momentHandler := handler.NewMomentHandler(momentService)

	//评论
	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	likeRepo := repository.NewLikeRepository(db)
	likeService := service.NewLikeService(likeRepo)
	likeHandler := handler.NewLikeHandler(likeService)

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.RegisterUserRoutes(r, cfg, userHandler)
	router.RegisterMomentRoutes(r, cfg, momentHandler, commentHandler, likeHandler)

	//加载服务端口号
	addr := ":" + cfg.App.Port
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server run failed:%v", err)
	}

	//r.Run(":8080")
}
