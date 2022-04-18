package github

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/wujunyi792/crispy-waffle-be/config"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/users"
	serviceErr "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/user"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"github.com/wujunyi792/crispy-waffle-be/internal/redis"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/github"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/jwtTokenGen"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/crypto"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/gen/xrandom"
	"net/url"
	"strconv"
	"time"
)

func makeNewSession(redirect string) string {
	sid := uuid.NewV4().String()
	err := redis.GetRedis().Set(sid, "oauth_state", time.Minute*10)
	if err != nil {
		logger.Error.Println(err)
		return ""
	}
	err = redis.GetRedis().Set(sid+"_redirect", redirect, time.Minute*10)
	if err != nil {
		logger.Error.Println(err)
		return ""
	}
	return sid
}

func HandleLogin(c *gin.Context) {
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = config.GetConfig().FrontendLogin
	}
	state := makeNewSession(redirect)
	query := url.Values{}
	query.Add("response_type", "code")
	query.Add("client_id", config.GetConfig().OAUTH.GITHUB.ClientId)
	query.Add("redirect_uri", config.GetConfig().OAUTH.GITHUB.RedirectUri)
	query.Add("state", state)
	query.Add("scope", config.GetConfig().OAUTH.GITHUB.Scope)
	redirectUrl := url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: query.Encode(),
	}
	c.Redirect(302, redirectUrl.String())
}

func HandleCallBack(c *gin.Context) {
	code, state := c.Query("code"), c.Query("state")
	token := ""
	if code == "" || state == "" {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	sid, err := redis.GetRedis().Get(state)
	if err != nil || sid != "oauth_state" {
		middleware.Fail(c, serviceErr.OauthErr)
		return
	}
	token = github.Code2Token(code)
	if token == "" {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	//logger.Debug.Println(token)

	githubInfo, err := github.GetGithubUserInfo(token)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	found, entity := users.GetEntityByGithubId(githubInfo.Id)
	// 找到用户记录 直接登录
	if found {
		serviceToken, err := jwtTokenGen.GenToken(jwtTokenGen.Info{UID: entity.ID})
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, user.LoginResponse{Token: serviceToken})
		users.SetLoginLog(entity.ID, token)
		return
	}

	// TODO 找不到记录，进入注册绑定流程（需要前端支持，暂不做）

	// 自动注册没有绑定手机号的账号
	salt := xrandom.GetRandom(6, xrandom.RAND_LOWER)
	newUser := Mysql.User{
		Base:      Mysql.Base{},
		NickName:  githubInfo.Login,
		Sex:       -1,
		Phone:     xrandom.GetRandom(10, xrandom.RAND_NUM),
		Signature: "这位用户没有任何想法~",
		Password:  crypto.PasswordGen(xrandom.GetRandom(6, xrandom.RAND_LOWER), salt),
		Salt:      salt,
		Avatar:    githubInfo.AvatarUrl,
		Oauth: []Mysql.Oauth{
			{
				OauthType:  "github",
				OauthId:    strconv.Itoa(int(githubInfo.Id)),
				UnionId:    "",
				Credential: token,
			},
		},
	}

	middleware.Success(c, "用户不存在")
}
