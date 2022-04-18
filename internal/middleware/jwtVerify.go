package middleware

import (
	"github.com/gin-gonic/gin"
	err2 "github.com/wujunyi792/crispy-waffle-be/internal/dto/err"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/jwtTokenGen"
)

func JwtVerify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		entry, err := jwtTokenGen.ParseToken(token)
		if err == nil {
			//c.Set("token", token)
			c.Set("uid", entry.Info.ID)
			c.Next()
			return
		} else {
			Fail(c, err2.JWTErr)
			return
		}
	}
	Fail(c, err2.VerifyErr)
	return
}
