package tasks

import (
)
import (
	"go-git-webhook/modules/queue"
	"go-git-webhook/conf"

)

var (

	queues = queue.NewQueue(conf.QueueSize())

	queueState uint64 = 0
)

type Task struct {
	SchedulerId int
}

func Add(task Task)  {
	queues.Enqueue(task)
}

func Start() {
	//如果等于0标识执行方法未开启
	if queueState == 0 {
		go func() {
			queueState = 1
			for {
				//停止了队列读取
				if queueState == 0 {
					break
				}
				value := queues.Dequue()

				if _,ok := value.(Task);ok {

					//fmt.Println(task.SchedulerId)
				}

			}
		}()
	}

}

func Stop(){
	queueState = 0
}