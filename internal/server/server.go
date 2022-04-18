package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/config"
	_ "github.com/wujunyi792/crispy-waffle-be/internal/corn"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/redis"
	v1 "github.com/wujunyi792/crispy-waffle-be/internal/router/v1"
)

var E *gin.Engine

func init() {
	logger.Info.Println("start init gin")
	gin.SetMode(config.GetConfig().MODE)
	E = gin.New()
	E.Use(middleware.GinRequestLog, gin.Recovery(), middleware.Cors(E))
}

func Run() {
	redis.GetRedis()
	v1.MainRouter(E)
	if err := E.Run("0.0.0.0:8080"); err != nil {
		logger.Error.Fatalln(err)
	}
}
