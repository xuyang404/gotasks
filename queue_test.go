package gotasks_test

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/xuyang404/gotasks"
	"github.com/xuyang404/gotasks/tasks"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := gotasks.NewQueue("test_queue")
	queue.UseRedisBroker(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	taskID, err := queue.Enqueue("add", tasks.ArgsMap{
		"args1": 1,
		"args2": 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(taskID)
}
