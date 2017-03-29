package tasks

import (
	"github.com/lifei6671/go-git-webhook/modules/queue"
	"github.com/lifei6671/go-git-webhook/conf"
	"strconv"
	"github.com/lifei6671/go-git-webhook/models"
	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/go-git-webhook/modules/ssh"
	"strings"
	"time"
	"github.com/lifei6671/go-git-webhook/modules/goclient"
	"net/url"
	"github.com/lifei6671/go-git-webhook/modules/hash"
	"fmt"
	"bufio"
	"io"
	"io/ioutil"
	"bytes"
)

var (

	queues = queue.NewQueue(conf.QueueSize())

)

type Task struct {
	SchedulerId int
	WebHookId int
	ServerId int
}

func Add(task Task)  {
	name := strconv.Itoa(task.WebHookId) + "-" + strconv.Itoa(task.ServerId)

	queues.Enqueue(name,task)
}

func Handle(value interface{})  {

	if task,ok := value.(Task);ok {
		scheduler := models.NewScheduler()
		scheduler.SchedulerId = task.SchedulerId

		if err := scheduler.Find(); err != nil {
			logs.Error("%s",err.Error())
			return
		}
		if scheduler.Status != "wait" {
			logs.Info("%s","Scheduler status does not wait.")
			return
		}
		server := models.NewServer()
		server.ServerId = task.ServerId

		if err := server.Find();err != nil {
			logs.Error("%s",err.Error())
			return
		}

		hook := models.NewWebHook()
		hook.WebHookId = task.WebHookId

		if err := hook.Find();err != nil {
			logs.Error("%s",err.Error())
			return
		}
		if strings.TrimSpace(hook.Shell) == "" {
			logs.Warn("","Shell command does not exist.")
			return
		}
		scheduler.StartExecTime = time.Now()
		scheduler.Status = "executing"
		scheduler.Save()





		channel := make(chan []byte,10)

		if server.Type == "ssh" {
			host := server.IpAddress + ":" + strconv.Itoa(server.Port)
			logs.Info("connecting ", host)
			go sshClient(host, scheduler, server, hook,channel)
		}else{
			host := server.IpAddress  + ":" + strconv.Itoa(server.Port)
			u ,err := url.Parse(host)
			if err != nil {
				u = &url.URL{ Host: host }
			}

			go clientClient(u.String(),scheduler,server,hook,channel)
		}
		buf := bytes.NewBufferString("")

		isChannelClosed := false
		for {
			if isChannelClosed {
				break
			}
			select {
			case out, ok := <-channel:
				{
					if !ok {
						isChannelClosed = true
						break
					}
					if len(out) > 0 {
						buf.Write(out)
					}

				}
			}
			if buf.Len() > 0 {
				logs.Info("%s", "The command was executed successfully")
				scheduler.LogContent = buf.String();
				scheduler.Status = "success"
				scheduler.EndExecTime = time.Now()
				scheduler.Save()
			}
		}

	}else{
		logs.Error("Can not be converted to Task:",value)
	}
}

func sshClient(host string,scheduler *models.Scheduler,server *models.Server,hook *models.WebHook,channel chan <-[]byte) {

	defer close(channel)

	logs.Info("connecting ", host)
	_,session,err := ssh.Connection(server.Account,host,server.PrivateKey)

	if err != nil {
		logs.Error("Connection remote server error:", err.Error())
		scheduler.Status = "failure"
		scheduler.LogContent = err.Error()
		scheduler.EndExecTime = time.Now()
		scheduler.Save()
		return
	}

	defer func() {
		if session != nil {
			session.Close()
		}
	}()

	logs.Info("SSH Server connectioned: " , host)

	stdout, err := session.StdoutPipe()

	if err != nil {
		fmt.Println("StdoutPipe: " + err.Error())
		channel <- []byte("StdoutPipe: " + err.Error())
		return
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		channel <- []byte("StderrPipe: " + err.Error())
		return
	}

	if err := session.Start(hook.Shell); err != nil {
		fmt.Println("Start: ", err.Error())
		channel <- []byte("Start: " + err.Error())
		return
	}

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line ,err2 := reader.ReadBytes('\n')

		if err2 != nil || io.EOF == err2 {
			break
		}
		channel <- line
	}

	bytesErr, err := ioutil.ReadAll(stderr)

	if err == nil {
		channel <- bytesErr

	}else{
		scheduler.Status = "failure"
		scheduler.LogContent = err.Error()
		scheduler.EndExecTime = time.Now()
		scheduler.Save()

		channel <- []byte("Stderr: " + err.Error())
	}

	if err := session.Wait(); err != nil {

		fmt.Println("Wait: ", err.Error())
		channel <- []byte("Wait: " +err.Error())
		return
	}

}

func clientClient(host string,scheduler *models.Scheduler,server *models.Server,hook *models.WebHook,channel chan<- []byte)  {

	defer close(channel)

	token,err := goclient.GetToken(host +"/token",server.Account,server.PrivateKey)

	if err != nil {
		logs.Error("Connection remote server error:", err.Error())
		scheduler.Status = "failure"
		scheduler.LogContent = err.Error()
		scheduler.EndExecTime = time.Now()
		scheduler.Save()
		return
	}

	u ,err := url.Parse(host)
	if err != nil {
		u = &url.URL{Scheme: "ws", Host: host , Path: "/socket"}
	}else{
		u = &url.URL{Scheme: "ws", Host: u.Host , Path: "/socket"}
	}


	client,err := goclient.Connection(u.String(),token)

	if err != nil {
		logs.Error("Remote server error:", err.Error())
		scheduler.Status = "failure"
		scheduler.LogContent = err.Error()
		scheduler.EndExecTime = time.Now()
		scheduler.Save()
		return
	}
	defer client.Close()

	client.SetCloseHandler(func(code int, text string) error {

		return nil
	})

	msg_id :=  hash.Md5(hook.Shell + time.Now().String())

	command := JsonResult{
		ErrorCode:0,
		Message:"ok",
		Command: "shell",
		MsgId: msg_id,
		Data:hook.Shell,
	}


	err = client.SendJSON(command)

	if err != nil {
		logs.Error("Remote server error:", err.Error())
		scheduler.Status = "failure"
		scheduler.LogContent = err.Error()
		scheduler.EndExecTime = time.Now()
		scheduler.Save()
		return
	}

	for {
		var response JsonResult

		 err := client.ReadJSON(&response)

		if err != nil {
			logs.Error("Remote server error:", err.Error())
			scheduler.Status = "failure"
			scheduler.LogContent = err.Error()
			scheduler.EndExecTime = time.Now()
			scheduler.Save()
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

type JsonResult struct {
	ErrorCode int                 	`json:"error_code"`
	Message string                	`json:"message"`
	Command string			`json:"command,omitempty"`
	MsgId string			`json:"msg_id,omitempty"`
	Data interface{}	      	`json:"data,omitempty"`
}

func init()  {
	queues.Handle = Handle
}