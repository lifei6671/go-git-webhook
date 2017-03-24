package goclient

import (
	"github.com/gorilla/websocket"
	"net/http"
	"errors"
	"strings"
	"io/ioutil"
	"github.com/widuu/gojson"
	"strconv"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type WebHookClient struct {
	token string
	conn *websocket.Conn

}
func Connection(remoteUrl  string,token string) (*WebHookClient ,error){
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
func GetToken(remoteUrl,account string,password string) (string,error) {


	t := strconv.Itoa(time.Now().Nanosecond())
	hash := sha256.New()
	hash.Write([]byte(account + password + t))
	md := hash.Sum(nil)
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