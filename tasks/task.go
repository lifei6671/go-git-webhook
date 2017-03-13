package tasks

import (
	"sync"
	"bytes"
	"github.com/astaxie/beego/logs"
)

var (
	SyncRWLock = &sync.RWMutex{}
	ConsoleMap 	= make(map[int]bytes.Buffer,20)
	SchedulerQueue = make(chan Task,20)
)

var(
	taskState = 0
)

// 设置队列的容量
func SetSchedulerQueueMaxCount(count int)  {
	SyncRWLock.Lock()
	queue := SchedulerQueue

	SchedulerQueue = make(chan Task,count)

	if len(queue) > 0 {
		for _,item := range queue {
			SchedulerQueue <- item
		}
	}

	SyncRWLock.Unlock()
}


type Task struct {
	SchedulerId int
}


func Start() {
	taskState = 1
	for {
		if taskState == 0{
			logs.Info("%s","The queue has stopped")
			break
		}
		if taskState == 2 {
			continue
		}
		select {
		case task,isClose := <-SchedulerQueue:
			{
				if !isClose {
					logs.Info("channel closed!")
					break
				}
				logs.Info("%s", "Start the task : ", task.SchedulerId)
			}
		}
	}
}
//暂停
func Paused()  {
	taskState = 2
}
//恢复
func Resume(){
	taskState = 1
}

func Stop(){
	taskState = 0
}

//获取队列状态
func State() int {
	return taskState
}