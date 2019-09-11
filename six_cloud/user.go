package six_cloud

import (
	"encoding/json"
	"errors"
	"github.com/Mrs4s/six-cli/models"
	"github.com/tidwall/gjson"
	"strconv"
)

func (user *SixUser) GetFilesByPath(path string) ([]*SixFile, error) {
	var (
		page = 2
		body = `{"path":"` + path + `","pageSize":50,"page": 1}`
		info = gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/page", body))
	)
	if !info.Get("success").Bool() {
		return nil, errors.New(info.Get("message").Str)
	}
	if info.Get("result.parent").Type == gjson.Null {
		return nil, errors.New("path not exists")
	}
	res := parseFiles(info.Get("result.list").Array())
	for int64(page) <= info.Get("result.totalPage").Int() {
		body = `{"path":"` + path + `","pageSize":50,"page": ` + strconv.FormatInt(int64(page), 10) + `}`
		info = gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/page", body))
		if !info.Get("success").Bool() {
			return res, nil
		}
		res = append(res, parseFiles(info.Get("result.list").Array())...)
		page++
	}
	for _, file := range res {
		file.owner = user
	}
	return res, nil
}

func (user *SixUser) GetFileByPath(path string) (*SixFile, error) {
	files, err := user.GetFilesByPath(models.GetParentPath(path))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name == models.GetFileName(path) {
			return file, nil
		}
	}
	return nil, errors.New("not found")
}

func (user *SixUser) GetDownloadAddressByPath(path string) (string, error) {
	body := `{"path":"` + path + `"}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/get", body))
	if !info.Get("success").Bool() {
		return "", errors.New(info.Get("message").Str)
	}
	return info.Get("result.downloadAddress").Str, nil
}

func (user *SixUser) CreateDirectory(path string) error {
	body := `{"path":"` + path + `"}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/createDirectory", body))
	if !info.Get("success").Bool() {
		return errors.New(info.Get("message").Str)
	}
	return nil
}

func (user *SixUser) DeleteFile(path string) error {
	body := `{"source":[{"path":"` + path + `"}]}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/delete", body))
	if !info.Get("success").Bool() {
		return errors.New(info.Get("message").Str)
	}
	return nil
}

func parseFiles(list []gjson.Result) []*SixFile {
	var res []*SixFile
	for _, r := range list {
		var file *SixFile
		err := json.Unmarshal([]byte(r.Raw), &file)
		if err == nil {
			res = append(res, file)
		}
	}
	return res
}
