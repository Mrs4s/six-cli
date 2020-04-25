package six_cloud

import (
	"errors"
	"github.com/Mrs4s/six-cli/models"
	"github.com/tidwall/gjson"
	"time"
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
		Shared         bool   `json:"share"`

		owner *SixUser
	}

	SixOfflineTask struct {
		Identity       string               `json:"identity"`
		UserIdentity   int64                `json:"userIdentity"`
		CreateTime     int64                `json:"createTime"`
		Name           string               `json:"name"`
		Type           int32                `json:"type"`
		Status         SixOfflineTaskStatus `json:"status"`
		TotalSize      int64                `json:"size"`
		DownloadedSize int64                `json:"downloadSize"`
		Progress       int32                `json:"progress"`
		AccessPath     string               `json:"accessPath"`

		ErrorCode    int32  `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}

	SixOfflineTaskStatus int
)

const (
	Failed           SixOfflineTaskStatus = -1
	Downloaded                            = 1000
	Downloading                           = 100
	AlmostDownloaded                      = 900
)

var (
	ErrWaitingLogin     = errors.New("waiting for login")
	ErrStateWrong       = errors.New("state wrong")
	ErrCreateDestFailed = errors.New("create destination failed")
	ErrDestExpired      = errors.New("destination expired")
	ErrLoginFailed      = errors.New("login failed")

	ErrInvalidToken = errors.New("invalid token")
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

func CreateDestination() (string, int64, error) {
	cli := models.NewSixHttpClient("")
	res := gjson.Parse(cli.PostJsonObject("https://api.6pan.cn/v3/user/createDestination", models.B{"ts": time.Now().Unix()}))
	dest := res.Get("destination")
	if !dest.Exists() {
		return "", 0, ErrCreateDestFailed
	}
	return dest.Str, res.Get("expireTime").Int(), nil
}

func LoginWithWebToken(dest, state string) (*SixUser, error) {
	cli := models.NewSixHttpClient("")
	res := gjson.Parse(cli.PostJsonObject("https://api.6pan.cn/v3/user/checkDestination", models.B{"destination": dest}))
	switch res.Get("status").Int() {
	case 10:
		return nil, ErrWaitingLogin
	case 100:
		if res.Get("state").Str != state {
			return nil, ErrStateWrong
		}
		user := &SixUser{
			Client: models.NewSixHttpClient(res.Get("token").Str),
		}
		user.RefreshUserInfo()
		return user, nil
	case -10:
		return nil, ErrDestExpired
	}
	return nil, ErrLoginFailed
}

func LoginWithAccessToken(token string) (*SixUser, error) {
	cli := models.NewSixHttpClient(token)
	user := &SixUser{
		Client: cli,
	}
	user.RefreshUserInfo()
	if user.Username == "" {
		return nil, ErrInvalidToken
	}
	return user, nil
}

func (task SixOfflineTask) StatusStr() string {
	switch task.Status {
	case Failed:
		return "下载失败"
	case Downloaded:
		return "下载完成"
	case Downloading:
		return "下载中"
	case AlmostDownloaded:
		return "部分下载完成"
	default:
		return "未知状态"
	}
}
