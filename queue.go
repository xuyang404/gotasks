package gotasks

import (
	"gotasks/brokers"
	"gotasks/tasks"
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

func (q *Queue) UseRedisBroker(redisURL string) {
	rb := brokers.NewRedisBroker(redisURL)
	q.SetBroker(rb)
}

func (q *Queue) Enqueue(taskName string, args tasks.ArgsMap) (string, error) {
	task := tasks.NewTask(q.Name, taskName, args)
	return q.broker.Enqueue(task)
}
