package gotasks

import (
	"github.com/go-redis/redis"
	"github.com/xuyang404/gotasks/brokers"
	"github.com/xuyang404/gotasks/tasks"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type Worker struct {
	broker        brokers.Broker
	errorHandler  func(interface{})
	taskMap       map[string][]TaskHandler
	reentrantMap  map[string]*ReentrantOptions
	taskMapLock   sync.RWMutex
	reentrantLock sync.RWMutex
}

func (w *Worker) SetErrorHandler(errorHandler func(interface{})) {
	w.errorHandler = errorHandler
}

func NewWorker() *Worker {
	return &Worker{
		taskMap:      make(map[string][]TaskHandler),
		reentrantMap: make(map[string]*ReentrantOptions),
	}
}

func (w *Worker) SetBroker(broker brokers.Broker) {
	w.broker = broker
}

func (w *Worker) UseRedisBroker(options *redis.Options) {
	rb := brokers.NewRedisBroker(options)
	w.SetBroker(rb)
}

func (w *Worker) Reentrant(handler TaskHandler, opts ...ReentrantOption) TaskHandler {
	handlerName := getHandleName(handler)
	reentrant := NewReentrantOptions(opts...)
	w.reentrantLock.Lock()
	defer w.reentrantLock.Unlock()
	if _, ok := w.reentrantMap[handlerName]; !ok {
		w.reentrantMap[handlerName] = reentrant
	}

	return handler
}

func (w *Worker) RegisterHandler(taskName string, handlers ...TaskHandler) {
	w.taskMapLock.Lock()
	defer w.taskMapLock.Unlock()
	w.taskMap[taskName] = handlers
}

func (w *Worker) handlerTask(task *tasks.Task) {

	defer func() {
		if r := recover(); r != nil {
			task.PanicLog = string(debug.Stack())
			w.broker.Enqueue(task) //再塞回队列重试
		}
	}()

	w.taskMapLock.Lock()
	handlers, ok := w.taskMap[task.TaskName]
	w.taskMapLock.Unlock()

	if !ok {
		log.Printf("task %s handler not exits", task.TaskName)
		return
	}

	var err error
	var argMap tasks.ArgsMap

	for i, handler := range handlers {
		handlerName := getHandleName(handler)
		//panic被捕获之后塞回队列重试要跳过之前执行过的handler
		if task.CurrentHandlerIndex > i {
			log.Printf("The task %s has already executed the handle %s, skip", task.ID, handlerName)
			continue
		}

		task.CurrentHandlerIndex = i

		w.reentrantLock.Lock()
		reentrant, ok := w.reentrantMap[handlerName]
		w.reentrantLock.Unlock()

		if ok {
			for i := 0; i < reentrant.maxLimit; i++ {
				argMap, err = handler(task)
				if err == nil {
					break
				}

				time.Sleep(time.Duration(reentrant.sleepMs) * time.Millisecond)
			}
		} else {
			argMap, err = handler(task)
		}

		if err != nil {
			w.errorHandler(err.Error())
			task.ResultLog = err.Error()
		}

		if argMap != nil {
			w.taskMapLock.Lock()
			argMap["funcName"] = handlerName
			w.taskMapLock.Unlock()
		}

		task.Result = append(task.Result, argMap)
		task.UpdateAt = time.Now().Format("2006-01-02 15:04:05")

		if err := w.broker.Update(task); err != nil {
			w.errorHandler("task update err: " + err.Error())
		}
	}
}

func (w *Worker) ListenMany(queueNames []string, concurrency int) {
	for _, queueName := range queueNames {
		w.listen(queueName, concurrency)
	}
}

func (w *Worker) Listen(queueName string, concurrency int) {
	w.listen(queueName, concurrency)
}

func (w *Worker) listen(queueName string, concurrency int) {
	log.Println("Listening......")
	pool := NewWaitGo(concurrency)
	pool.PanicHandler = w.errorHandler
	defer pool.Wait()

	for {
		fn := func() {
			task, err := w.broker.Acquire(queueName)

			if err == brokers.NoDatas {
				return
			}

			if err != nil {
				w.errorHandler("task Acquire err: " + err.Error())
				return
			}

			w.handlerTask(task)
		}

		if concurrency > 0 {
			pool.Add(fn)
		} else {
			fn()
		}
	}

}
