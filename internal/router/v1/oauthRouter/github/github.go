package github

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/handle/oauth/github"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
)

func InitGithubRouter(e *gin.RouterGroup) {
	githubRouter := e.Group("/github")
	{
		githubRouter.GET("/login", github.HandleLogin)
		githubRouter.GET("/bind", middleware.JwtVerify, github.HandleBindAccount)
		githubRouter.GET("/callback", github.HandleCallBack)
	}
}
