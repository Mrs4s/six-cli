package six_cloud

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	user, err := LoginWithAccessToken("--")
	files, err := user.GetFilesByPath("/宿星のガールフレンド/123")
	fmt.Println(user, files, err)
}
