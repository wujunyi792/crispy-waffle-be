package github

import (
	"github.com/parnurzeal/gorequest"
	"github.com/wujunyi792/crispy-waffle-be/config"
	"github.com/wujunyi792/crispy-waffle-be/internal/logger"
	"net/http"
	"net/url"
	"time"
)

type Code2TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserInfo struct {
	Login                   string      `json:"login"`
	Id                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                interface{} `json:"location"`
	Email                   interface{} `json:"email"`
	Hireable                interface{} `json:"hireable"`
	Bio                     interface{} `json:"bio"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		Collaborators int    `json:"collaborators"`
		PrivateRepos  int    `json:"private_repos"`
	} `json:"plan"`
}

func Code2Token(code, state string) (token string) {
	query := make(url.Values)
	query.Add("client_id", config.GetConfig().OAUTH.GITHUB.ClientId)
	query.Add("client_secret", config.GetConfig().OAUTH.GITHUB.ClientSecret)
	query.Add("code", code)
	reqUrl := url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/access_token",
		RawQuery: query.Encode(),
	}
	res := Code2TokenResponse{}
	_, _, err := gorequest.New().Get(reqUrl.String()).AppendHeader("Accept", "application/json").
		Retry(3, time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		EndStruct(&res)
	if err != nil {
		logger.Error.Println(err)
		return ""
	}
	return res.AccessToken
}
