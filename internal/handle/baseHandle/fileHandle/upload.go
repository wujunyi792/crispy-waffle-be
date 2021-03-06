package fileHandle

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/oss"
)

func HandleGetAliUploadToken(c *gin.Context) {
	data := oss.GetPolicyToken()
	middleware.Success(c, data)
}

// HandleAliUpLoad 通过业务服务器中转文件至OSS 表单提交 字段名upload
func HandleAliUpLoad(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		middleware.FailWithCode(c, 20008, err.Error())
	} else {
		url := oss.UploadFileToOss(header.Filename, file)
		if url == "" {
			middleware.FailWithCode(c, 50006, err.Error())
		}
		middleware.Success(c, url)
	}
}
