package queue

import (
	"sync"
	"container/list"
	"runtime"
)

type Queue struct {
	maxLength int
	mutex *sync.RWMutex
	cache *list.List
	temp list.List
}

func NewQueue(maxLength int) *Queue  {
	return &Queue{
		maxLength	: maxLength,
		mutex 		: &sync.RWMutex{},
		cache 		: list.New(),
	}
}

func (q *Queue) Enqueue(value interface{}) {
	for q.cache.Len() >= q.maxLength {
		runtime.Gosched()
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.cache.PushBack(value)
}

func (q *Queue) Dequue() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if element := q.cache.Front();element != nil {
		q.cache.Remove(element)
		return element.Value
	}
	return nil
}
func (q *Queue) Length() int  {

	return q.cache.Len()
}
