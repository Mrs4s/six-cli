package six_cloud

import (
	"encoding/json"
	"errors"
	"github.com/Mrs4s/six-cli/models"
	"github.com/tidwall/gjson"
	"time"
)

func (user *SixUser) RefreshUserInfo() {
	res := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/user/info", models.B{"ts": time.Now().Unix()}))
	if res.Get("success").Exists() {
		return
	}
	user.Identity = res.Get("identity").Int()
	user.Username = res.Get("name").Str
	user.TotalSpace = res.Get("spaceCapacity").Int()
	user.UsedSpace = res.Get("spaceUsed").Int()
}

func (user *SixUser) GetFilesByPath(path string) ([]*SixFile, error) {
	res := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/files/list",
		models.B{"parentPath": path, "skip": 0, "limit": 500}))
	arr := parseFiles(res.Get("dataList").Array())
	for _, file := range arr {
		file.owner = user
	}
	return arr, nil
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

func (user *SixUser) GetOfflineTasks() ([]*SixOfflineTask, error) {
	var (
		body = `{"page": 1,"pageSize": 200}`
		info = gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/offline/page", body))
		res  []*SixOfflineTask
	)
	if !info.Get("success").Bool() {
		return nil, errors.New(info.Get("message").Str)
	}
	for _, token := range info.Get("result.list").Array() {
		var task *SixOfflineTask
		err := json.Unmarshal([]byte(token.Raw), &task)
		if err == nil {
			res = append(res, task)
		}
	}
	return res, nil
}

func (user *SixUser) GetDownloadAddressByPath(path string) (string, error) {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/file/download", models.B{"path": path}))
	if info.Get("success").Exists() {
		return "", errors.New(info.Get("message").Str)
	}
	return info.Get("downloadAddress").Str, nil
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

func (user *SixUser) CopyFile(source, target string) error {
	body := `{"source": [{"path": "` + source + `"}],"destination": {"path": "` + target + `"}}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/copy", body))
	if !info.Get("success").Bool() {
		return errors.New(info.Get("message").Str)
	}
	return nil
}

func (user *SixUser) SearchFilesByName(name string) ([]*SixFile, error) {
	body := `{"pageSize":200,"name":"` + name + `"}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/pageAll", body))
	if !info.Get("success").Bool() {
		return nil, errors.New(info.Get("message").Str)
	}
	files := parseFiles(info.Get("result.list").Array())
	return files, nil
}

func (user *SixUser) PreparseOffline(url, pass string) (string, string, int64, error) {
	body := `{"url": "` + url + `","password": "` + pass + `"}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/offline/parseUrl", body))
	if !info.Get("success").Bool() {
		return "", "", 0, errors.New(info.Get("messages").Str)
	}
	if len(info.Get("result").Array()) == 0 {
		return "", "", 0, errors.New("not any results")
	}
	return info.Get("result.0.identity").Str, info.Get("result.0.name").Str, info.Get("result.0.size").Int(), nil
}

func (user *SixUser) AddOfflineTask(identity, path string) error {
	body := `{"path": "` + path + `","task":[{"identity" : "` + identity + `"}]}`
	info := gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/offline/add", body))
	if !info.Get("success").Bool() {
		return errors.New(info.Get("message").Str)
	}
	if !info.Get("result.success").Bool() {
		return errors.New("unknown error")
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
