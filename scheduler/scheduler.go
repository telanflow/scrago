package scheduler

import (
	"context"
	"sync"
	"time"
)

type TaskChan chan QueueElement
type SchedulerOptions func(*Scheduler)

type Scheduler struct {
	list		Queue
	handler		func(QueueElement)

	mu 			*sync.Mutex
	num			uint
	threads		uint
	taskChan	TaskChan
	stopSign	context.CancelFunc
}

func New() *Scheduler {
	return &Scheduler{
		list:		NewListQueue(),
		taskChan:	make(TaskChan, 20),
		mu:			new(sync.Mutex),
		threads:	10,
		num:		0,
	}
}

func (self *Scheduler) Use(options ...SchedulerOptions) {
	for _, v := range options {
		v(self)
	}
}

func (self *Scheduler) Push(v QueueElement) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if ok := self.list.Push(v); ok {
		self.num++
	}
}

func (self *Scheduler) Pop() QueueElement {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.num--
	return self.list.Pop()
}

func (self *Scheduler) Count() uint {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.num
}

func (self *Scheduler) task() <-chan QueueElement {
	return self.taskChan
}

func (self *Scheduler) taskComplete() {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.num--
}

// 调度器：停止分发任务
func (self *Scheduler) StopDispatch() {
	if self.stopSign != nil {
		self.stopSign()
	}
}

// 调度器：分发任务
func (self *Scheduler) Dispatch() {
	var (
		wg	sync.WaitGroup
		ctx	context.Context
	)

	ctx, self.stopSign = context.WithCancel(context.Background())
	go func() {
		defer self.stopSign()
		defer close(self.taskChan)

		for {
			if self.Count() <= 0 {
				break
			}

			req := self.list.Pop()
			if req != nil {
				self.taskChan <- req
				continue
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()

	for i := 0; uint(i) < self.threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case task, ok := <-self.task():
					if !ok {
						// TaskChan Closed
						return
					}

					self.handler(task)

					self.taskComplete()

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
}

func WithHandler(handler func(QueueElement)) SchedulerOptions {
	return func(s *Scheduler) {
		s.handler = handler
	}
}