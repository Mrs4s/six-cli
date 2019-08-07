package main

import (
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/six_cloud"
	"io/ioutil"
)

var (
	currentUser *six_cloud.SixUser
	currentPath = "/"
)

func main() {
	if models.PathExists("token.info") {
		bytes, err := ioutil.ReadFile("token.info")
		if err == nil {
			currentUser, _ = six_cloud.LoginWithAccessToken(string(bytes))
		}
	}
	runAsShell()
	/*
		if len(os.Args) == 1 {
			fmt.Println("usage: six-cli <command> or six-cli shell")
			return
		}
		if os.Args[1] == "shell" {
			runAsShell()
			return
		}
	*/
}
