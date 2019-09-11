package models

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
)

// 6盘
type SixHttpClient struct {
	QingzhenToken string

	client *http.Client
}

func NewSixHttpClient(token string) *SixHttpClient {
	cli := &SixHttpClient{
		QingzhenToken: token,
		client:        &http.Client{},
	}
	return cli
}

func (cli *SixHttpClient) PostJson(url, body string) string {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return ""
	}
	defer req.Body.Close()
	if cli.QingzhenToken != "" {
		req.Header["Qingzhen-Token"] = []string{cli.QingzhenToken}
	}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:67.0) Gecko/20100101 Firefox/67.0"}
	req.Header["Content-Type"] = []string{"application/json"}
	resp, err := cli.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	token, ok := resp.Header["Qingzhen-Token"]
	if ok && len(token) > 0 {
		cli.QingzhenToken = token[0]
	}
	return string(b)
}

func (cli *SixHttpClient) GetBytes(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		buffer := bytes.NewBuffer(body)
		r, _ := gzip.NewReader(buffer)
		unCom, err := ioutil.ReadAll(r)
		return unCom, err
	}
	return body, nil
}
func (cli *SixHttpClient) GetString(url string) string {
	bytes, err := cli.GetBytes(url)
	if err != nil {
		return ""
	}
	return string(bytes)
}
