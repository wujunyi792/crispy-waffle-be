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

func makeNewBindSession(redirect string, uid string) string {
	sid := uuid.NewV4().String()
	err := redis.GetRedis().Set(sid, uid, time.Minute*10)
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
	if err != nil || (sid != "oauth_state" && len(sid) != 36) {
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

	// ????????????
	if sid != "oauth_state" {
		userEntity := Mysql.User{}
		userEntity.ID = sid
		users.GetEntity(&userEntity)
		found, entity := users.GetEntityByGithubId(githubInfo.Id)
		if found {
			if userEntity.ID == entity.ID {
				middleware.FailWithCode(c, 40213, "????????????????????????????????????")
				return
			}
			middleware.FailWithCode(c, 40212, "??? Github ??????????????????????????????"+entity.NickName)
			return
		}
		err = users.AddGithubOauth(&Mysql.Oauth{
			UserID:     sid,
			OauthType:  "github",
			OauthId:    strconv.Itoa(int(githubInfo.Id)),
			UnionId:    githubInfo.NodeId,
			Credential: token,
		})
		if err != nil {
			logger.Error.Println(err)
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		redirect, _ := redis.GetRedis().Get(state + "_redirect")
		if redirect == "" {
			redirect = config.GetConfig().FrontendLogin
		}
		c.Redirect(302, redirect)

		return
	}

	// ????????????
	found, entity := users.GetEntityByGithubId(githubInfo.Id)
	// ?????????????????? ????????????
	if found {
		serviceToken, err := jwtTokenGen.GenToken(jwtTokenGen.Info{UID: entity.ID, InfoComplete: true})
		if err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		middleware.Success(c, user.LoginResponse{Token: serviceToken})
		users.SetLoginLog(entity.ID, token)
		return
	}

	// ??????????????????????????????????????????
	salt := xrandom.GetRandom(6, xrandom.RAND_LOWER)
	newUser := Mysql.User{
		Base:      Mysql.Base{},
		NickName:  githubInfo.Login,
		Sex:       -1,
		Phone:     xrandom.GetRandom(10, xrandom.RAND_NUM),
		Signature: "??????????????????????????????~",
		Password:  crypto.PasswordGen(xrandom.GetRandom(6, xrandom.RAND_LOWER), salt),
		Salt:      salt,
		Avatar:    githubInfo.AvatarUrl,
		Oauth: []Mysql.Oauth{
			{
				OauthType:  "github",
				OauthId:    strconv.Itoa(int(githubInfo.Id)),
				UnionId:    githubInfo.NodeId,
				Credential: token,
			},
		},
	}
	logger.Debug.Println("????????????", entity)

	err = users.RegisterUser(&newUser)

	if err != nil {
		logger.Error.Println(err)
		middleware.FailWithCode(c, 40210, "????????????????????????")
		return
	}

	serviceToken, err := jwtTokenGen.GenToken(jwtTokenGen.Info{UID: entity.ID, InfoComplete: false})
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	c.Redirect(302, config.GetConfig().FrontendLogin+"/#/auth?token="+serviceToken)
	users.SetLoginLog(entity.ID, token)
}

func HandleBindAccount(c *gin.Context) {
	cuid, _ := c.Get("uid")
	uid := cuid.(string)
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = config.GetConfig().FrontendLogin
	}
	state := makeNewBindSession(redirect, uid)
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
