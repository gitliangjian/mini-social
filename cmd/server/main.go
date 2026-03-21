package main

import (
	"log"
	"net/http"

	"mini-social/internal/bootstrap"
	"mini-social/internal/config"

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

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//加载服务端口号
	addr := ":" + cfg.App.Port
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server run failed:%v", err)
	}

	//r.Run(":8080")
}
