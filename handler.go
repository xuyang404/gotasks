package gotasks

import (
	"errors"
	"gotasks/tasks"
	"reflect"
	"runtime"
)

type TaskHandler func(tasks.ArgsMap) (tasks.ArgsMap, error)

type ReentrantOptions struct {
	maxLimit int //最大重试次数
	sleepMs  int //重试间隔数，毫秒
}

type ReentrantOption func(options *ReentrantOptions)

func getHandleName(handle interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
}

func CallHandle(handle interface{}, args ...interface{}) (interface{}, error) {
	value := reflect.ValueOf(handle)
	if value.Kind() != reflect.Func {
		return nil, errors.New("function is not function")
	}

	if value.Type().NumIn() != len(args) {
		return nil, errors.New("The parameter length does not match")
	}

	argValues := make([]reflect.Value, 0, value.Type().NumIn())
	for i := 0; i < value.Type().NumIn(); i++ {
		//param := value.Type().In(i)
		argValue := reflect.ValueOf(args[i])
		//argType := argValue.Type().String()
		//if param.String() != argType {
		//	return nil, errors.New(fmt.Sprintf("function %s argument %d is not valid", getHandleName(handle), i+1))
		//}
		argValues = append(argValues, argValue)
	}

	resultValues := value.Call(argValues)

	var results []interface{}
	for _, result := range resultValues {

		isErr1 := result.Type().AssignableTo(reflect.TypeOf((*error)(nil)).Elem())

		isErr2 := result.Type().Implements(reflect.TypeOf((*error)(nil)).Elem())

		if isErr1 || isErr2 {
			if result.Interface() != nil {
				return nil, result.Interface().(error)
			}
		}

		results = append(results, result.Interface())
	}

	return results, nil
}

func WithMaxLimit(maxLimit int) ReentrantOption {
	return func(options *ReentrantOptions) {
		options.maxLimit = maxLimit
	}
}

func WithSleepMs(ms int) ReentrantOption {
	return func(options *ReentrantOptions) {
		options.sleepMs = ms
	}
}

func NewReentrantOptions(opts ...ReentrantOption) *ReentrantOptions {
	o := &ReentrantOptions{
		maxLimit: 3,
		sleepMs:  20,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
