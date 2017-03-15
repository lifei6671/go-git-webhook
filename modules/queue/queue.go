package queue

import (
	"sync"
	"time"
)

type Queue struct {
	mutex *sync.RWMutex
	contents map[string]*Element
	Handle func(interface{})
}

func NewQueue(queueLength int) *Queue {
	return &Queue{
		mutex	: &sync.RWMutex{},
		contents: make(map[string]*Element,queueLength),
	}
}

type Element struct {
	mutex *sync.RWMutex
	contents chan interface{}
}

func (self *Element) Push(value interface{}) {
	self.contents <- value
}


func (self *Element) worker(worker func(interface{})) {

	for {
		select {
		case element := <- self.contents : {
			worker(element)
		}
		case <- time.After(time.Second * 1):
			return
		}
	}
}


func (q *Queue) Enqueue(name string,value interface{})  {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if element,ok := q.contents[name]; ok {
		element.Push(value)
	}else{
		element := &Element{
			mutex		: &sync.RWMutex{},
			contents	: make(chan interface{},20),
		}
		element.Push(value)
		q.contents[name] = element
		if q.Handle != nil {
			go element.worker(q.Handle)
		}

	}
}
