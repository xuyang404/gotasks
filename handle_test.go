package gotasks

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

func add(arg1 int, arg2 int) (int, error) {
	sum := arg1 + arg2
	return sum, errors.New("出错了哦")
}

func add2(m interface{}) (map[string]interface{},error) {
	a := m.(map[string]interface{})
	return a, nil
}

func add3(arg1 , arg2 interface{}) (map[string]interface{}, error) {
	sum := arg1.(int) + arg2.(int)
	return map[string]interface{}{"sum": sum}, nil
}

func TestCallHandle(t *testing.T) {
	res, err := CallHandle(add2, map[string]interface{}{"a": "abcd"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
