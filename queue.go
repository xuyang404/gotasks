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

func NewQueue(Name string) *Queue {
	return &Queue{
		Name: Name,
	}
}

func (q *Queue) SetBroker(broker brokers.Broker) {
	q.broker = broker
}

func (q *Queue) UseRedisBroker(options *redis.Options) {
	rb := brokers.NewRedisBroker(options)
	q.SetBroker(rb)
}

func (q *Queue) Enqueue(taskName string, args tasks.ArgsMap) (string, error) {
	task := tasks.NewTask(q.Name, taskName, args)
	return q.broker.Enqueue(task)
}
