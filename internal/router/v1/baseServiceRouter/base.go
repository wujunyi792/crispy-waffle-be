package baseServiceRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/baseHandle/captchaHandle"
)

func InitBaseServiceRouter(e *gin.Engine) {
	baseGroup := e.Group("/base")
	{
		baseGroup.GET("/captcha", captchaHandle.HandleGetCaptcha)
	}
}
