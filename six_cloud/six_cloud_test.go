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

func TestOffline(t *testing.T) {
	user, _ := LoginWithAccessToken("--")
	tasks, err := user.GetOfflineTasks()
	fmt.Println(tasks, err)
	identity, name, size, err := user.PreparseOffline("magnet:?xt=urn:btih:1536cc0e486e5d649dacacdb13947cb72c64e8d5&dn=ZOMBIE%20LAND%20SAGA%20Special%20Disc%20Collection%20%5BFLAC%5D%20%5B44.1kHz%2F16bit%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce", "")
	fmt.Println(identity, name, size, err)
}
