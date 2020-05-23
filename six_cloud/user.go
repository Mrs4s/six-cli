package six_cloud

import (
	"encoding/json"
	"errors"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/models/fs"
	"github.com/tidwall/gjson"
	"time"
)

func (user *SixUser) GetRootFile() *SixFile {
	return &SixFile{
		Identity:     "root",
		UserIdentity: user.Identity,
		Path:         "/",
		Name:         "root",
		IsDir:        true,
		owner:        user,
	}
}

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
	id := models.ToIdentity(path)
	b, err := user.Client.GetBytes("https://api.6pan.cn/v3/file/" + id)
	if err != nil {
		return nil, err
	}
	var file *SixFile
	err = json.Unmarshal(b, &file)
	if err != nil {
		return nil, err
	}
	file.owner = user
	return file, nil
}

func (user *SixUser) GetOfflineTasks() ([]*SixOfflineTask, error) {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/offline/list", models.B{"skip": 0, "limit": 500}))
	var res []*SixOfflineTask
	if info.Get("success").Exists() {
		return nil, errors.New(info.Get("message").Str)
	}
	for _, token := range info.Get("dataList").Array() {
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
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/file", models.B{"path": path}))
	if info.Get("success").Exists() {
		return errors.New(info.Get("message").Str)
	}
	return nil
}

func (user *SixUser) DeleteFile(path string) error {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/file/trash", models.B{"sourcePath": []string{path}}))
	if info.Get("successCount").Int() != 1 {
		return errors.New("delete failed")
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

func (user *SixUser) SearchFilesByName(parent, name string) ([]*SixFile, error) {
	if parent == "" {
		parent = "::all"
	}
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/files/list/", models.B{"parentIdentity": parent, "name": name}))
	arr := parseFiles(info.Get("dataList").Array())
	for _, file := range arr {
		file.owner = user
	}
	return arr, nil
}

func (user *SixUser) PreparseOffline(url, pass string) (string, string, int64, error) {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/offline/parse", models.B{"textLink": url, "password": pass}))
	if info.Get("success").Exists() {
		return "", "", 0, errors.New(info.Get("message").Str)
	}
	return info.Get("hash").Str, info.Get("info.name").Str, info.Get("info.size").Int(), nil
}

func (user *SixUser) AddOfflineTask(hash, path string) error {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/offline/add", models.B{"task": []models.B{{"hash": hash}}, "savePath": path}))
	if info.Get("successCount").Int() != 1 {
		return errors.New("unsuccessful")
	}
	return nil
}

func (user *SixUser) CreateUploadTree(remote string, files []string) map[string]string {
	res := make(map[string]string)
	for _, f := range files {
		p := remote + "/" + fs.GetFileName(f)
		if fs.IsDir(f) {
			_ = user.CreateDirectory(p)
			for k, v := range user.CreateUploadTree(p, fs.GetDirEntities(f)) {
				res[k] = v
			}
			continue
		}
		res[p] = f
	}
	return res
}

func (user *SixUser) CreateUploadToken(path, name, hash string) SixUploadToken {
	info := gjson.Parse(user.Client.PostJsonObject("https://api.6pan.cn/v3/file/uploadToken", models.B{"path": path, "name": name, "hash": hash}))
	if info.Get("created").Bool() {
		return SixUploadToken{Cached: true}
	}
	return SixUploadToken{
		UploadToken: info.Get("uploadToken").Str,
		UploadUrl:   info.Get("partUploadUrl").Str,
	}
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
