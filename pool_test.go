package gotasks

import (
	"fmt"
	"testing"
	"time"
)

func TestWaitGo(t *testing.T) {
	wg := NewWaitGo(3)

	for i := 0; i < 10; i++ {
		a := i
		wg.Add(func() {
			fmt.Println("a", a)
			time.Sleep(1 * time.Second)
		})
	}

	wg.Wait()
}
