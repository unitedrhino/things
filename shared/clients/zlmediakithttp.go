package clients

//部分门接口直接http访问
import (
	"fmt"
	"github.com/i-Things/things/src/vidsvr/pb/vid"
	"io"
	"net/http"
	"net/url"
)

type ApiResponse struct {
	Data string `json:"data"`
}

var PreUrl string

// ZLMediakit初始化数据结构
type MediaConfig struct {
	ID     string
	Ipv4   string
	Port   int64
	Secret string
	PreUrl string
}

const (
	VIDMGRTIMEOUT = 60
)

func NewMeidaServer(vmgrInfo *vid.VidmgrInfo) *MediaConfig {
	return &MediaConfig{
		Ipv4:   vmgrInfo.VidmgrIpV4,
		Port:   vmgrInfo.VidmgrPort,
		Secret: vmgrInfo.VidmgrSecret,
		PreUrl: fmt.Sprintf("http://%s:%d/index/api/", vmgrInfo.VidmgrIpV4, vmgrInfo.VidmgrPort),
	}
}

func (f *MediaConfig) PostMediaServer(strurl string, values url.Values) (data []byte, err error) {
	values.Add("secret", f.Secret)
	//fmt.Println("[---------------]PostMediaServer -", f.PreUrl+strurl, "param:", values)
	resp, err := http.PostForm(f.PreUrl+strurl, values)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}
