package oauthRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/router/v1/oauthRouter/github"
)

func InitOauthRouter(e *gin.Engine) {
	oauthGroup := e.Group("/oauth")
	github.InitGithubRouter(oauthGroup)
}
