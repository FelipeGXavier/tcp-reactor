package data

import (
	"container/list"
)

type Queue interface {
	Pop() interface{}
	Len() int 
	Add(interface{})
	Seek() interface{}
}

type LinkedQueue struct {
	list *list.List
}

func MakeLinkedQueue() Queue {
	queue := LinkedQueue{list: list.New()}
	return &queue
}

func (self *LinkedQueue) Add(element interface{}) {
	self.list.PushBack(element)
}

func (self *LinkedQueue) Len() int {
	return self.list.Len()
}

func (self *LinkedQueue) Seek() interface{} {
	return self.list.Front().Value
}

func (self *LinkedQueue) Pop() interface{} {
	element := self.list.Front()
	if element == nil {
		return nil
	}
	self.list.Remove(element)
	return element.Value
}
