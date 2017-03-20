package tasks

import (
	"go-git-webhook/modules/queue"
	"go-git-webhook/conf"
	"strconv"
	"go-git-webhook/models"
	"github.com/astaxie/beego/logs"
	"go-git-webhook/modules/ssh"
	"strings"
	"time"
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

		host := server.IpAddress + ":" + strconv.Itoa(server.Port)

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

		//str := strings.TrimSpace(strings.Replace(hook.Shell,"\r","",-1))
		//
		//shells := strings.Split(str,"\n")
		//
		//logs.Info("%s",shells)
		//
		//stdout := bytes.NewBufferString("")
		//
		//for _,shell := range shells {
		//	out,err := session.CombinedOutput(shell);
		//
		//	if err != nil{
		//		logs.Error("%+v",err)
		//		break
		//	}
		//	logs.Info("Command : ",shell)
		//	logs.Info("Execute Result:",string(out))
		//	stdout.Write(out)
		//
		//	session.Stderr = nil
		//	session.Stdout = nil
		//}

		out,err := session.CombinedOutput(hook.Shell);

		if err != nil{
			logs.Error("CombinedOutput:",err.Error())
			scheduler.Status = "failure"
			scheduler.LogContent = err.Error()
			scheduler.EndExecTime = time.Now()
			scheduler.Save()
			return
		}
		logs.Info("%s","The command was executed successfully")
		scheduler.LogContent = string(out);
		scheduler.Status = "success"
		scheduler.EndExecTime = time.Now()
		scheduler.Save()
	}

}

func init()  {
	queues.Handle = Handle
}