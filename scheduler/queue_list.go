package scheduler

import (
	"container/list"
)

type ListQueue struct {
	list	*list.List
	only	map[string] bool
}

// 创建一个新的双向链表
func NewListQueue() *ListQueue {
	return &ListQueue {
		list: list.New(),
		only: make(map[string]bool),
	}
}

// 推入元素
func (self *ListQueue) Push(v QueueElement) bool {
	var s = v.Id()
	if _, ok := self.only[s]; !ok {
		self.only[s] = true
		self.list.PushBack(v)
		return true
	} else {
		return false
	}
}

// 弹出元素
func (self *ListQueue) Pop() QueueElement {

	if self.list.Len() == 0 {
		return nil
	}

	return self.list.Remove(self.list.Front()).(QueueElement)
}

// 统计数量
func (self *ListQueue) Count() uint {
	return uint(self.list.Len())
}