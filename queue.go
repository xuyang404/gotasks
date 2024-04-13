package gotasks

import (
	"github.com/go-redis/redis"
	"github.com/xuyang404/gotasks/brokers"
	"github.com/xuyang404/gotasks/tasks"
)

type Queue struct {
	broker brokers.Broker
	Name   string
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) SetBroker(broker brokers.Broker) {
	q.broker = broker
}

func (q *Queue) UseRedisBroker(options *redis.Options) {
	rb := brokers.NewRedisBroker(options)
	q.SetBroker(rb)
}

func (q *Queue) Enqueue(task *tasks.Task) (string, error) {
	return q.broker.Enqueue(task)
}
