package fileRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/baseHandle/fileHandle"
)

func InitFileRouter(e *gin.Engine) {
	file := e.Group("/fileHandle")
	{
		file.GET("/ali/token", fileHandle.HandleGetAliUploadToken)
		file.POST("/ali/upload", fileHandle.HandleAliUpLoad)
	}
}
