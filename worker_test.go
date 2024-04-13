package gotasks

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/xuyang404/gotasks/tasks"
	"testing"
)

type Mul struct {
	Args1 int `json:"args1"`
	Args2 int `json:"args2"`
	Sum   int `json:"sum"`
}

func TestWorker(t *testing.T) {
	worker := NewWorker()
	worker.UseRedisBroker(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	worker.SetErrorHandler(func(i interface{}) {
		fmt.Println(i)
	})
	//
	add := func(task *tasks.Task) (tasks.ArgsMap, error) {
		m := map[string]int{}

		fmt.Println("taskId", task.ID)

		err := tasks.ArgsMapTo(task.ArgsMap, &m)
		if err != nil {
			return nil, err
		}

		m["sum"] = m["args1"] + m["args2"]

		am, err := tasks.ToArgsMap(m)
		if err != nil {
			return nil, err
		}

		return am, nil
	}

	multiplication := func(task *tasks.Task) (tasks.ArgsMap, error) {

		mul := &Mul{}
		err := tasks.ArgsMapTo(task.ArgsMap, mul)
		if err != nil {
			return nil, err
		}

		mul.Sum = mul.Args1 * mul.Args2

		am, err := tasks.ToArgsMap(mul)
		if err != nil {
			return nil, err
		}

		return am, nil
	}
	worker.RegisterHandler("add", add, worker.Reentrant(multiplication, WithMaxLimit(3), WithSleepMs(3)))

	worker.Listen("test_queue", 0)
}
