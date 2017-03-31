package tasks

import (
	"strconv"
	"strings"
	"time"
	"net/url"
	"bytes"

	"github.com/astaxie/beego/logs"
	"github.com/lifei6671/go-git-webhook/modules/queue"
	"github.com/lifei6671/go-git-webhook/conf"
	"github.com/lifei6671/go-git-webhook/models"
	"github.com/lifei6671/go-git-webhook/modules/goclient"
	"errors"
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
			logs.Error("",err.Error())
			return
		}
		if scheduler.Status != "wait" {
			logs.Info("%s","Scheduler status does not wait.")
			return
		}
		server := models.NewServer()
		server.ServerId = task.ServerId

		if err := server.Find();err != nil {
			logs.Error("",err.Error())
			return
		}

		hook := models.NewWebHook()
		hook.WebHookId = task.WebHookId

		if err := hook.Find();err != nil {
			logs.Error("",err.Error())
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

		client,err := CreateClient(server.Type)

		if err != nil {
			logs.Error("",err.Error())
			return
		}


		host := server.IpAddress + ":" + strconv.Itoa(server.Port)

		u,err := url.Parse(host)

		scheme := "http"

		if strings.HasPrefix(server.IpAddress,"http://") {
			scheme = "http"
		}else if strings.HasPrefix(server.IpAddress,"https://") {
			scheme = "https"
		}else{
			scheme = "ssh"
		}

		if err != nil {
			u = &url.URL{ Scheme : scheme, Host: host}
		}

		logs.Info("connecting ", u)

		go client.Command(*u, server.Account,server.PrivateKey, hook.Shell,channel)


		buf := bytes.NewBufferString("")

		isChannelClosed := false
		var lastInfo string

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
						lastInfo = string(out)
					}

				}
			}
			if buf.Len() > 0 {
				logs.Info("%s", "The command was executed successfully")
				scheduler.LogContent = buf.String();
				scheduler.Status = "executing"
				scheduler.EndExecTime = time.Now()
				scheduler.Save()
			}
		}

		if strings.HasPrefix(lastInfo,"Error") {
			scheduler.Status = "failure"
		}else{
			scheduler.Status = "success"
		}
		scheduler.Save()

	}else{
		logs.Error("Can not be converted to Task:",value)
	}
}

func CreateClient(t string) (goclient.ClientInterface,error) {
	if t == "ssh" {
		return &goclient.SSHClient{},nil
	}else if t == "client" {
		return &goclient.WebHookClient{},nil
	}else{
		return nil,errors.New("未知的客户端类型")
	}
}

func init()  {
	queues.Handle = Handle
}