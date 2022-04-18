package user

import (
	"github.com/gin-gonic/gin"
	"github.com/wujunyi792/crispy-waffle-be/internal/controller/users"
	serviceErr "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/user"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/middleware"
	"github.com/wujunyi792/crispy-waffle-be/internal/model/Mysql"
	"github.com/wujunyi792/crispy-waffle-be/internal/redis"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/jwtTokenGen"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/tecentCMS"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/captcha"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/check"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/crypto"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/gen/cmscode"
	"github.com/wujunyi792/crispy-waffle-be/pkg/utils/gen/xrandom"
	"time"
)

func HandleSendRegisterCode(c *gin.Context) {
	var req user.SendCode
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40204, "验证码错误")
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40205, "手机号已经注册，可以直接登录")
		return
	}

	_, err := redis.GetRedis().Get(req.Phone + "_register")
	if err == nil {
		middleware.FailWithCode(c, 40203, "发送过于频繁，稍后再试")
		return
	}
	data := cmscode.GenValidateCode(6)
	err = redis.GetRedis().Set(req.Phone+"_register", data, 5*time.Minute)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	tecentCMS.SendCMS(req.Phone, []string{data})
	middleware.Success(c, nil)
}

func HandleCheckPhoneExist(c *gin.Context) {
	var req user.CheckPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	exist := users.CheckPhoneExist(req.Phone)
	middleware.Success(c, user.CheckPhoneResponse{
		Phone: req.Phone,
		Exist: exist,
	})
}

func HandleRegister(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40201, "手机号已存在")
		return
	}

	err := check.PasswordStrengthCheck(6, 20, 3, req.Password)
	if err != nil {
		middleware.FailWithCode(c, 40202, err.Error())
		return
	}

	code, err := redis.GetRedis().Get(req.Phone + "_register")
	if err != nil || code != req.Code {
		middleware.Fail(c, serviceErr.CodeErr)
		return
	}

	redis.GetRedis().RemoveKey(req.Phone+"_passwordReset", false)

	salt := xrandom.GetRandom(5, xrandom.RAND_ALL)
	entity := Mysql.User{
		NickName:  "hgame_" + xrandom.GetRandom(10, xrandom.RAND_NUM),
		Sex:       -1,
		Phone:     req.Phone,
		Signature: "这位用户没有任何想法~",
		Password:  crypto.PasswordGen(req.Password, salt),
		Salt:      salt,
	}
	err = users.RegisterUser(&entity)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, entity.Phone)
}

func HandleGeneralLogin(c *gin.Context) {
	var req user.LoginGeneralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	entity := Mysql.User{Phone: req.Info}
	users.GetEntity(&entity)

	if entity.ID == "" {
		middleware.Fail(c, serviceErr.LoginErr)
		return
	}

	if !crypto.PasswordCompare(req.Password, entity.Password, entity.Salt) {
		middleware.Fail(c, serviceErr.LoginErr)
		return
	}
	token, err := jwtTokenGen.GenToken(jwtTokenGen.Info{UID: entity.ID})
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, user.LoginResponse{Token: token})
	users.SetLoginLog(entity.ID, token)
}

func HandleSendPasswordResetCode(c *gin.Context) {
	var req user.SendCode
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40204, "验证码错误")
		return
	}

	_, err := redis.GetRedis().Get(req.Phone + "_passwordReset")
	if err == nil {
		middleware.FailWithCode(c, 40203, "发送过于频繁，稍后再试")
		return
	}
	data := cmscode.GenValidateCode(6)
	err = redis.GetRedis().Set(req.Phone+"_register", data, 5*time.Minute)
	if err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	tecentCMS.SendCMS(req.Phone, []string{data})
	middleware.Success(c, nil)
}

func HandleResetPassword(c *gin.Context) {
	var req user.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if req.Password != req.PasswordAgain {
		middleware.FailWithCode(c, 40206, "密码不一致")
		return
	}
	if !users.CheckPhoneExist(req.Phone) {
		middleware.FailWithCode(c, 40207, "手机号不存在")
		return
	}
	err := check.PasswordStrengthCheck(6, 20, 3, req.Password)
	if err != nil {
		middleware.FailWithCode(c, 40202, err.Error())
		return
	}

	code, err := redis.GetRedis().Get(req.Phone + "_passwordReset")
	if err != nil || code != req.Code {
		middleware.Fail(c, serviceErr.CodeErr)
		return
	}

	_ = redis.GetRedis().RemoveKey(req.Phone+"_passwordReset", false)

	salt := xrandom.GetRandom(5, xrandom.RAND_ALL)
	passwordHashed := crypto.PasswordGen(req.Password, salt)
	err = users.UpdatePasswordAndSalt(req.Phone, passwordHashed, salt)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, req.Phone)
}

func HandleUpdateNickName(c *gin.Context) {
	var req user.UpdateNickNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if users.CheckUserNameExist(req.NickName) {
		middleware.FailWithCode(c, 40208, "用户名已存在")
		return
	}
	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateUserName(uid, req.NickName)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateAvatar(c *gin.Context) {
	var req user.UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateAvatar(uid, req.AvatarUrl)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateSex(c *gin.Context) {
	var req user.UpdateSexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateSex(uid, req.Sex)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateStatus(c *gin.Context) {
	var req user.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateStatus(uid, req.Status)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateSignature(c *gin.Context) {
	var req user.UpdateSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateSignature(uid, req.Signature)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleUpdateEmail(c *gin.Context) {
	var req user.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}

	if !check.VerifyEmailFormat(req.Email) {
		middleware.FailWithCode(c, 40209, "邮箱格式错误")
		return
	}

	cuid, _ := c.Get("uid")
	uid := cuid.(string)

	err := users.UpdateEmail(uid, req.Email)
	if err != nil {
		logger.Error.Println(err)
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	middleware.Success(c, nil)
}

func HandleGetUserInfo(c *gin.Context) {
	cuid, _ := c.Get("uid")
	entity := Mysql.User{}
	entity.ID = cuid.(string)
	users.GetEntity(&entity)
	middleware.Success(c, entity)
}

func HandleDelAccount(c *gin.Context) {
	middleware.Success(c, "功能暂不支持，如有需求，请联系工作人员")
}
