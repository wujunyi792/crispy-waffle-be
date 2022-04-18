package captchaHandle

import (
	"github.com/gin-gonic/gin"
	err2 "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/captcha"
)

func HandleGetCaptcha(c *gin.Context) {
	cp, err := captcha.GenerateCaptcha()
	if err != nil {
		middleware.Fail(c, err2.InternalErr)
		return
	}
	middleware.Success(c, *cp)
}
