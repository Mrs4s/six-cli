package six_cloud

import (
	"errors"
	"github.com/Mrs4s/six-cli/models"
	"github.com/tidwall/gjson"
)

type (
	SixUser struct {
		Username   string
		Identity   int64
		UsedSpace  int64
		TotalSpace int64

		Client *models.SixHttpClient
	}

	SixFile struct {
		Identity       string `json:"identity"`
		ETag           string `json:"hash"`
		UserIdentity   int64  `json:"userIdentity"`
		Path           string `json:"path"`
		Name           string `json:"name"`
		Size           int64  `json:"size"`
		CreateTime     int64  `json:"ctime"`
		Mime           string `json:"mime"`
		ParentIdentity string `json:"parent"`
		IsDir          bool   `json:"directory"`
	}
)

func LoginWithUsernameOrPhone(value, password string) (*SixUser, error) {
	var (
		body = `{"value":"` + value + `","password":"` + models.ToMd5(password) + `","code":""}`
		cli  = models.NewSixHttpClient("")
		res  = cli.PostJson("https://api.6pan.cn/v2/user/login", body)
	)
	if res == "" {
		return nil, errors.New("login failed")
	}
	info := gjson.Parse(res)
	if !info.Get("success").Bool() {
		return nil, errors.New(info.Get("message").Str)
	}
	return LoginWithAccessToken(cli.QingzhenToken)
}

func LoginWithAccessToken(token string) (*SixUser, error) {
	cli := models.NewSixHttpClient(token)
	info := gjson.Parse(cli.PostJson("https://api.6pan.cn/v2/user/info", "{}"))
	if !info.Get("success").Bool() {
		return nil, errors.New("login failed: token error")
	}
	user := &SixUser{
		Username:   info.Get("result.name").Str,
		Identity:   info.Get("result.identity").Int(),
		UsedSpace:  info.Get("result.spaceUsed").Int(),
		TotalSpace: info.Get("result.spaceCapacity").Int(),
		Client:     cli,
	}
	return user, nil
}
