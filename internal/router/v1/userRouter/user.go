package userRouter

import (
	"github.com/gin-gonic/gin"
	user "github.com/wujunyi792/crispy-waffle-be/internal/handle/userHandle"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitUserRouter(e *gin.Engine) {
	userGroup := e.Group("/user")
	{
		userGroup.POST("/register/code", user.HandleSendRegisterCode)
		userGroup.POST("/reset/code", user.HandleSendPasswordResetCode)
		userGroup.POST("/register", user.HandleRegister)
		userGroup.POST("/exist/phone", user.HandleCheckPhoneExist)
		userGroup.POST("/login/general", user.HandleGeneralLogin)
		userGroup.POST("/pwd/reset", user.HandleResetPassword)

		userGroup.Use(middleware.JwtVerify)
		{
			userGroup.GET("/info", user.HandleGetUserInfo)

			userGroup.POST("/update/avatar", user.HandleUpdateAvatar)
			userGroup.POST("/update/nickname", user.HandleUpdateNickName)
			userGroup.POST("/update/sex", user.HandleUpdateSex)
			userGroup.POST("/update/signature", user.HandleUpdateSignature)
			userGroup.POST("/update/status", user.HandleUpdateStatus)
			userGroup.POST("/update/email", user.HandleUpdateEmail)

			userGroup.GET("/delete", user.HandleDelAccount)
		}
	}
}
