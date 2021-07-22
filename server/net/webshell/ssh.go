package webshell

import (
	"bufio"
	"encoding/json"
	"errors"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/go-homedir"
	goSsh "golang.org/x/crypto/ssh"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	ModeList string
}

type SshLoginModel struct {
	Addr     string
	UserName string
	SshType  int
	Pwd      string
	PemKey   string
	PtyCols  uint32
	PtyRows  uint32
}

//todo 新增安全审计，对登录的信息进行记录并记录用户操作SSH返回的数据以及增加该条操作的正常和非正常标记 正常 警告  危险 致命 VIM
//todo 记录SFTP日志，对SFTP操作进行记录
//todo 增加跳板机功能，连接时可以选择跳板机进行连接（新增设备时予以选择）

// 创建一个ssh连接
func sshConnect(login SshLoginModel) (client *goSsh.Client, ch goSsh.Channel, session *goSsh.Session, err error) {
	//创建ssh登陆配置
	config := &goSsh.ClientConfig{}
	config.SetDefaults()
	config.Timeout = time.Second * 2
	config.User = login.UserName
	if login.Pwd == "" {
		return
	} else {
		if login.SshType == 1 {
			config.Auth = []goSsh.AuthMethod{goSsh.Password(login.Pwd)}
		} else if login.SshType == 2 {
			config.Auth = []goSsh.AuthMethod{publicKeyAuthFunc(login.Pwd)}
		} else {
			return
		}
	}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key goSsh.PublicKey) error { return nil }
	//dial 获取ssh client
	client, err = goSsh.Dial("tcp", login.Addr, config)
	if err != nil {
		return
	}
	channel, incomingRequests, err := client.Conn.OpenChannel("session", nil)
	if err != nil {
		return
	}
	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				err = req.Reply(false, nil)
				if err != nil {
					return
				}
			}
		}
	}()
	modes := goSsh.TerminalModes{
		goSsh.ECHO:          1,
		goSsh.TTY_OP_ISPEED: 14400,
		goSsh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, goSsh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  login.PtyCols,
		Rows:     login.PtyRows,
		Width:    login.PtyCols * 8,
		Height:   login.PtyRows * 17,
		ModeList: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, goSsh.Marshal(&req))
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("e001")
		return
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("e002")
		return
	}
	ch = channel
	//创建ssh-session
	session, _ = client.NewSession()

	return
}

// Request WS数据接受结构体
type Request struct {
	MsgType  int    `json:"msg_type"`  //如果为1进行ssh连接验证为2进行代码执行
	Token    string `json:"token"`     //用户token
	ServerID string `json:"server_id"` //wehShell服务器id
	Command  string `json:"command"`   //用户命令
	Cols     int    `json:"cols"`      //终端窗口的列数, 可以在创建Terminal指定cols（大概和分辨率的比拟是1：7.5）
	Rows     int    `json:"rows"`      //终端窗口的行数, 可以在创建Terminal指定rows（大概和分辨率的比拟是1：18）
}

// WebSocketHandler 启动一个WS并进行SSH数据交互
func WebSocketHandler(w http.ResponseWriter, r *http.Request, checkUserToken func(string) bool, getServerInfo func(string, int, int) SshLoginModel) {
	ws, err := upGrader.Upgrade(w, r, nil)
	if nil != err {
		return
	}
	isConnect := false
	var channel goSsh.Channel
	var client *goSsh.Client
	var sshSession *goSsh.Session
	defer func() {
		if isConnect {
			_ = channel.Close()
			_ = client.Close()
		}
		_ = ws.Close()
	}()
	done := make(chan bool, 2)
	go func() {
		defer func() {
			done <- true
		}()
		for {
			_, msgByte, err := ws.ReadMessage()
			if err != nil {
				return
			}
			req := Request{}
			err = json.Unmarshal(msgByte, &req)
			if err != nil {
				return
			}
			switch req.MsgType {
			case 1:
				if !isConnect {
					if !checkUserToken(req.Token) {
						return
					}
					loginInfo := getServerInfo(req.ServerID, req.Cols, req.Rows)
					if loginInfo.Addr == "" {
						return
					}
					client, channel, sshSession, err = sshConnect(loginInfo)
					if err != nil {
						err := ws.WriteMessage(1, []byte("\n\n\n\033[31m 登录失败!\033[0m"))
						if err != nil {
							return
						}
						err = ws.Close()
						return
					} else {
						isConnect = true
					}
				}
			case 2:
				if isConnect {
					if !checkUserToken(req.Token) {
						return
					}
					if _, err := channel.Write([]byte(req.Command)); nil != err {
						return
					}
				}
			case 3:
				if err := sshSession.WindowChange(req.Cols, req.Rows); err != nil {
					err := ws.WriteMessage(1, []byte("\n\n\n\033[31m 分辨率重置成功!\033[0m"))
					if err != nil {
						return
					}
					continue
				}
			default:
				return
			}
		}
	}()
	go func() {
		for {
			if !isConnect {
				time.Sleep(time.Millisecond * 200)
				continue
			}
			defer func() {
				done <- true
			}()
			br := bufio.NewReader(channel)
			var buf []byte
			t := time.NewTimer(time.Millisecond * 50)
			defer t.Stop()
			r := make(chan rune)
			go func() {
				for {
					x, size, err := br.ReadRune()
					if err != nil {
						err := ws.WriteMessage(1, []byte("\n\n\n\033[31m 已经关闭连接!\033[0m"))
						if err != nil {
							return
						}
						err = ws.Close()
						if err != nil {
							return
						}
						return
					}
					if size > 0 {
						r <- x
					}
				}
			}()
			for {
				select {
				case <-t.C:
					if len(buf) != 0 {
						err = ws.WriteMessage(websocket.TextMessage, buf)
						buf = []byte{}
						if err != nil {
							return
						}
					}
					t.Reset(time.Millisecond * 50)
				case d := <-r:
					if d != utf8.RuneError {
						p := make([]byte, utf8.RuneLen(d))
						utf8.EncodeRune(p, d)
						buf = append(buf, p...)
					} else {
						buf = append(buf, []byte("@")...)
					}
				}
			}
		}
	}()
	<-done
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func publicKeyAuthFunc(kPath string) goSsh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		mid.Log().Errorf("find key's home dir failed %v", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		mid.Log().Errorf("ssh key file read failed %v", err)
	}
	// Create the Signer for this private key.
	signer, err := goSsh.ParsePrivateKey(key)
	if err != nil {
		mid.Log().Errorf("ssh key signer failed %v", err)
	}
	return goSsh.PublicKeys(signer)
}
