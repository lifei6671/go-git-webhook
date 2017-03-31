package goclient

import (

	"net/http"
	"errors"
	"strings"
	"io/ioutil"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"net/url"

	"github.com/lifei6671/go-git-webhook/modules/hash"

	"github.com/astaxie/beego/logs"
	"github.com/widuu/gojson"
	"github.com/gorilla/websocket"
)

type WebHookClient struct {
	token string
	conn *websocket.Conn

}

func (p *WebHookClient)Connection(remoteUrl  string,token string) (*WebHookClient ,error){
	client :=  &WebHookClient{

	}

	header := http.Header{}
	header.Add("x-smarthook-token",token)

	c, _, err := websocket.DefaultDialer.Dial(remoteUrl,header)

	if err != nil {
		return client,err
	}

	client.conn = c
	return client,nil
}

func (c *WebHookClient) SetCloseHandler(h func(code int, text string) error) {
	c.conn.SetCloseHandler(h)
}

func (c *WebHookClient)Send(msg []byte) error {

	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *WebHookClient) SendJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *WebHookClient) Read() ([]byte,error) {
	_, message, err := c.conn.ReadMessage()

	return message,err
}

func (c *WebHookClient) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}
func (c *WebHookClient) Close() {
	c.conn.Close()
}

// GetToken 获取与服务器端连接的认证密钥.
func GetToken(remoteUrl,account string,password string) (string,error) {

	t := strconv.Itoa(time.Now().Nanosecond())
	h := sha256.New()
	h.Write([]byte(account + password + t))
	md := h.Sum(nil)
	mdStr := hex.EncodeToString(md)

	response, err := http.Post(remoteUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader("account=" + account +"&password=" + mdStr + "&time=" + t))
	if err != nil {
		return "",err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "",err
	}

	if js := gojson.Json(string(body)).Get("error_code");js.IsValid() {
		error_code,err := strconv.Atoi(js.Tostring())

		if err != nil {
			return "",err
		}
		if error_code != 0 {
			message := gojson.Json(string(body)).Get("message");

			return "",errors.New(message.Tostring())
		}
		token := gojson.Json(string(body)).Get("data")

		return token.Tostring(),nil

	}
	return "",errors.New("Data error")

}

// Command 执行命令.
func (p *WebHookClient) Command (host url.URL,account,password ,shell string,channel chan <-[]byte) {

	defer close(channel)


	token,err := GetToken(host.String() +"/token",account,password)

	if err != nil {
		logs.Error("Connection remote server error:", err.Error())

		channel <- []byte("Error: Connection remote server error => " + err.Error())
		return
	}


	u := &url.URL{Scheme: "ws", Host: host.Host , Path: "/socket"}

	client,err := (&WebHookClient{}).Connection(u.String(),token)

	if err != nil {
		logs.Error("Remote server error:", err.Error())

		channel <- []byte("Error:Remote server error => " + err.Error())
		return
	}

	defer client.Close()

	client.SetCloseHandler(func(code int, text string) error {

		return nil
	})

	msg_id :=  hash.Md5(shell + time.Now().String())

	command := JsonResult{
		ErrorCode	: 0,
		Message		: "ok",
		Command		: "shell",
		MsgId		: msg_id,
		Data		: shell,
	}


	err = client.SendJSON(command)

	if err != nil {
		logs.Error("Remote server error:", err.Error())

		channel <- []byte("Error:Remote server error => " + err.Error())
		return
	}

	for {
		var response JsonResult

		err := client.ReadJSON(&response)

		if err != nil {
			logs.Error("Remote server error:", err.Error())

			channel <- []byte("Error:Remote server error => " + err.Error())
			return
		}
		if response.ErrorCode == 0 {
			if response.Command == "end" {
				return
			}
			body := response.Data.(string)

			channel <- []byte(body)
		}
	}
}