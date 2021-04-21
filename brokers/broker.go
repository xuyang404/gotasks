package brokers

import (
	"errors"
	"gotasks/tasks"
)

var NoDatas = errors.New("queue no data")

type Broker interface {
	Acquire(queueName string) (*tasks.Task, error)
	Ack() bool
	Update(task *tasks.Task) error
	Enqueue(task *tasks.Task) (string, error)
	QueueLen(queueName string) int64
}
