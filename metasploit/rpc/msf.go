package rpc

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/vmihailenco/msgpack.v2"
)

// sessionlist 请求的结构体
type SessionListReq struct {
	_msgpack struct{} `msgpack:",asArray"` // 当做索引数组解析
	Method   string
	Token    string
}

// sessionList 的响应
type SessionListRes struct {
	ID          uint32 `msgpack:",omitempty"` // 可选参数
	Type        string `msgpack:"type"`
	TunnelLocal string `msgpack:"tunnel_local"`
	TunnelPeer  string `msgpack:"tunnel_peer"`
	ViaExploit  string `msgpack:"via_exploit"`
	ViaPayload  string `msgpack:"via_payload"`
	Desc        string `msgpack:"desc"`
	Info        string `msgpack:"info"`
	Workspace   string `msgpack:"workspack"`
	SessionHost string `msgpack:"session_host"`
	SessionPort int    `msgpack:"session_port"`
	Username    string `msgpack:"username"`
	UUID        string `msgpack:"uuid"`
	ExploitUUID string `msgpack:"exploit_uuid"`
}

// 登录请求
type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Pass     string
}

// 登录返回
type loginRes struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

//登出请求
type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

// 登出响应
type logoutRes struct {
	Result string `msgpack:"result"`
}

// 通用信息
type Msf struct {
	host  string
	user  string
	pass  string
	token string
}

// 初始化
func New(host, user, pass string) (*Msf, error) {
	rtn := &Msf{
		host: host,
		user: user,
		pass: pass,
	}

	if err := rtn.Login(); err != nil {
		return nil, err
	}

	return rtn, nil
}

func (msf *Msf) send(req interface{}, res interface{}) error {
	buf := new(bytes.Buffer) //https://blog.csdn.net/flyfreelyit/article/details/80291945  bytes.Buffer 使用
	// encodereq放到buf中

	msgpack.NewEncoder(buf).Encode(req)
	dst := fmt.Sprintf("http://%s/api", msf.host)
	resp, err := http.Post(dst, "binary/message-pack", buf)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	defer resp.Body.Close()

	if err = msgpack.NewDecoder(resp.Body).Decode(res); err != nil {
		log.Printf("%s", err)
		return err
	}
	fmt.Println(res)

	return nil

}

func (msf *Msf) Login() error {
	ctx := &loginReq{
		Method:   "auth.login",
		Username: msf.user,
		Pass:     msf.pass,
	}
	var res loginRes
	// send 的第二个参数为interface 可以接收任何类型
	if err := msf.send(ctx, &res); err != nil {
		log.Printf("%s", err)
		return err
	}
	msf.token = res.Token
	return nil
}

func (msf *Msf) Logout() error {
	ctx := &logoutReq{
		Method:      "auth.logout",
		Token:       msf.token,
		LogoutToken: msf.token,
	}
	var res logoutRes

	if err := msf.send(ctx, &res); err != nil {
		log.Println(err)
		return err
	}

	msf.token = ""
	return nil
}

func (msf *Msf) SessionList() (map[uint32]SessionListRes, error) {
	req := &SessionListReq{
		Method: "session.list",
		Token:  msf.token,
	}

	res := make(map[uint32]SessionListRes)

	if err := msf.send(req, &res); err != nil {
		log.Fatal(err)
		return nil, err
	}
	

	for id, session := range res {
		session.ID = id
		res[id] = session
	}

	return res, nil
}
