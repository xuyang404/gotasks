package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func call(function interface{}, args ...string) {
	value := reflect.ValueOf(function)
	if value.Kind() != reflect.Func {
		fmt.Println("function is not function")
		return
	}

	parameters := make([]reflect.Type, 0, value.Type().NumIn())
	for i := 0; i < value.Type().NumIn(); i++ {
		arg := value.Type().In(i)
		fmt.Printf("argument %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
		parameters = append(parameters, arg)
	}

	if value.Type().NumIn() != len(args) {
		fmt.Printf("argument %d length doesn't equal to provide length %d \n", value.Type().NumIn(), len(args))
		return
	}

	outs := make([]reflect.Type, 0, value.Type().NumOut())
	for i := 0; i < value.Type().NumOut(); i++ {
		arg := value.Type().Out(i)
		fmt.Printf("out %d is %s[%s] type \n", i, arg.Kind(), arg.Name())
		outs = append(outs, arg)
	}

	if value.Type().NumOut() < 1 {
		fmt.Println("outs length must greater than 0")
		return
	}

	if !outs[len(outs)-1].AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		fmt.Println("last output must be error")
		return
	}
	if !outs[len(outs)-1].Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		fmt.Println("last output must be error")
		return
	}

	argValues := make([]reflect.Value, 0, len(parameters))
	for i := 0; i < len(args); i++ {
		switch parameters[i] {
		case reflect.TypeOf(int(0)):
			v, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				fmt.Printf("argument %d %s convert int failed, %v \n", i, args[i], err)
				return
			}
			argValues = append(argValues, reflect.ValueOf(int(v)))
		case reflect.TypeOf(string("")):
			argValues = append(argValues, reflect.ValueOf(args[i]))
		default:
			fmt.Printf("unsupport type %s[%s] \n", parameters[i].Kind(), parameters[i].Name())
			return
		}
	}

	resultValues := value.Call(argValues)
	for i := 0; i < len(resultValues); i++ {
		switch resultValues[i].Type() {
		case reflect.TypeOf(int(0)):
			fmt.Println("result: ", i, ", ", resultValues[i].Interface().(int))
		case reflect.TypeOf(string("")):
			fmt.Println("result: ", i, ", ", resultValues[i].Interface().(string))
		default:
			fmt.Printf("type: %s[%s], value: %v \n", resultValues[i].Type().Kind(), resultValues[i].Type().Name(), resultValues[i].Interface())
		}
	}
}

func f1(arg1 string, arg2 int) (string, error) {
	fmt.Println("arg1: ", arg1, "  arg2: ", arg2)
	return "f1", errors.New("出错了")
}

func main() {
	call(f1, "", "100")
}