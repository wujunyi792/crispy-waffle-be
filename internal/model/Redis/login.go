package Redis

import (
	"encoding/json"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"github.com/wujunyi792/crispy-waffle-be/internal/service/github"
)

type GithubLogin struct {
	Info  github.GithubUserInfo
	State string
}

func NewGithubLoginString(m *github.GithubUserInfo, state string) []byte {
	var res GithubLogin
	res.State = state
	res.Info = *m
	marshal, _ := json.Marshal(res)
	return marshal
}

func RevertGithubLoginStruct(raw []byte) (res GithubLogin) {
	err := json.Unmarshal(raw, &res)
	if err != nil {
		logger.Error.Println(err)
	}
	return
}
