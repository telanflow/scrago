package scheduler

import "net/http"

type QueueElement interface {
	Id() string
	GetRequest() *http.Request
}

type Queue interface {
	Push(QueueElement) bool
	Pop() QueueElement
	Count() uint
}
