package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
)

type Result struct {
	Id         string `json:"id"`
	Base64Blog string `json:"raw"`
}

// 默认存储10240个验证码，每个验证码10分钟过期
var store = base64Captcha.DefaultMemStore

func GenerateCaptcha() (*Result, error) {
	// 生成默认数字
	//driver := base64Captcha.DefaultDriverDigit
	// 此尺寸的调整需要根据网站进行调试，链接：
	// https://captcha.mojotv.cn/
	driver := base64Captcha.NewDriverDigit(70, 130, 4, 0.8, 100)
	// 生成base64图片
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 获取
	id, b64s, err := captcha.Generate()
	if err != nil {
		logger.Error.Println("Register GetCaptchaPhoto get base64Captcha has err:", err)
		return nil, err
	}

	return &Result{
		Id:         id,
		Base64Blog: b64s,
	}, nil
}

func VerifyCaptcha(id string, value string) bool {
	verifyResult := store.Verify(id, value, true)
	return verifyResult
}
