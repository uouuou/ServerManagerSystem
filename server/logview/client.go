package logview

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	controlC    = "^C"
	messageType = websocket.TextMessage
	tailN       = "n"
)

// Request WS数据接受结构体
type Request struct {
	MsgType int    `json:"msg_type"` //如果为1进行ssh连接验证为2进行代码执行
	Token   string `json:"token"`    //用户token
	LogFile string `json:"log_file"` //日志位置
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	// The websocket connection.
	conn *websocket.Conn
	// Context撤销函数，用于关闭日志文件的读取
	ctxCancelFunc context.CancelFunc
	// 需要读取的日志文件
	logFile string
	// tail 参数
	tailOptions []string
	//连接状态
	IsClient bool `json:"is_client"`
}

// read 读取客户端websocket消息
func (c *Client) read() {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			c.ctxCancelFunc()
			break
		}
		req := Request{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			mid.Log().Error(err.Error())
			return
		}
		message = bytes.TrimSpace(message)
		switch req.MsgType {
		case 1:
			if !checkUserToken(req.Token) {
				c.IsClient = false
				c.ctxCancelFunc()
				break
			} else {
				c.IsClient = true
			}
		case 2:
			c.ctxCancelFunc()
			break
		}
		if string(message) == controlC {
			c.ctxCancelFunc()
			break
		}
		c.logFile = req.LogFile
	}
}

// write 向客户端写websocket消息
func (c *Client) write(ctx context.Context) {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}()
	if !c.IsClient {
		message, err := WriteMessage("Token验证失败", 1)
		if err != nil {
			return
		}
		err = c.conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
		return
	}
	if c.logFile == "" {
		message, err := WriteMessage("未指定日志文件或日志文件不存在", 1)
		if err != nil {
			return
		}
		err = c.conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
		return
	}
	if !FileExists(c.logFile) {
		message, err := WriteMessage("日志文件["+c.logFile+"]不存在", 1)
		if err != nil {
			return
		}
		err = c.conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
		return
	}
	c.tailOptions = append(c.tailOptions, "-f", c.logFile)
	cmd := exec.CommandContext(ctx, "tail", c.tailOptions...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		message, err := WriteMessage("调用tail命令出错:"+err.Error(), 1)
		if err != nil {
			return
		}
		err = c.conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		msg := map[string]interface{}{
			"msg":      line,
			"msg_type": 1,
		}
		marshal, err := json.Marshal(msg)
		if err != nil {
			mid.Log().Error(err.Error())
			return
		}
		err = c.conn.WriteMessage(messageType, marshal)
		if err != nil {
			return
		}
	}
	err = cmd.Wait()
	if err != nil {
		return
	}
}

// serveWs 处理相应客户端的websocket请求
func serveWs(w http.ResponseWriter, r *http.Request, logFile string, tailOptions []string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		mid.Log().Error(err.Error())
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	client := &Client{conn: conn, ctxCancelFunc: cancelFunc, logFile: logFile, tailOptions: tailOptions}

	go client.read()
	for {
		if client.logFile != "" {
			go client.write(ctx)
			break
		}
		continue
	}
}

// LogWs 日志文件输出
func LogWs(c *gin.Context) {
	//paths := c.Query("log")
	var path string
	tailOptions := make([]string, 0)
	_, err := strconv.Atoi(tailN)
	if err == nil {
		tailOptions = append(tailOptions, "-n", tailN)
	}
	serveWs(c.Writer, c.Request, path, tailOptions)

}

//检查token的合法性
func checkUserToken(token string) bool {
	if _, err := mid.ParseToken(token); err != nil {
		return false
	} else {
		return true
	}
}

// FileExists 判断文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

// WriteMessage 数据反馈
func WriteMessage(line string, msgType int) (send []byte, err error) {
	msg := map[string]interface{}{
		"msg":      line,
		"msg_type": msgType,
	}
	send, err = json.Marshal(msg)
	if err != nil {
		mid.Log().Error(err.Error())
		return send, err
	}
	return send, nil
}
