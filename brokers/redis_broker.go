package brokers

import (
	"github.com/go-redis/redis"
	"github.com/xuyang404/gotasks/json"
	"github.com/xuyang404/gotasks/tasks"
	"log"
)

var _ Broker = &RedisBroker{}

type RedisBroker struct {
	client *redis.Client
}

func (r *RedisBroker) Acquire(queueName string) (*tasks.Task, error) {
	//TODO implement me
	task := &tasks.Task{}
	str, err := r.client.RPop(queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, NoDatas
		}
		return nil, err
	}

	if err := json.Json.Unmarshal([]byte(str), &task); err != nil {
		return nil, err
	}

	return task, nil
}

func (r *RedisBroker) Ack() bool {
	//TODO implement me
	return true
}

func (r *RedisBroker) Update(task *tasks.Task) error {
	//TODO implement me
	b, err := json.Json.Marshal(task)
	if err != nil {
		return err
	}
	r.client.HSet(genSaveTaskName(task), task.ID, b)
	return nil
}

func (r *RedisBroker) Enqueue(task *tasks.Task) (string, error) {
	//TODO implement me
	b, err := json.Json.Marshal(task)
	if err != nil {
		return "", err
	}
	_, err = r.client.LPush(task.QueueName, string(b)).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	log.Printf("task %s enqueued to %s, taskID is %s", task.TaskName, task.QueueName, task.ID)
	return task.ID, nil
}

func (r *RedisBroker) QueueLen(queueName string) int64 {
	//TODO implement me
	l, _ := r.client.LLen(queueName).Result()
	return l
}

func NewRedisBroker(options *redis.Options) *RedisBroker {
	rb := &RedisBroker{}
	rb.client = redis.NewClient(options)
	return rb
}

func genSaveTaskName(task *tasks.Task) string {
	return "gt:" + task.QueueName + ":" + task.TaskName
}
