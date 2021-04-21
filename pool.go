package gotasks

import "log"

type WaitGo struct {
	token        chan struct{}
	concurrency  int
	PanicHandler func(interface{})
}

func NewWaitGo(concurrency int) *WaitGo {

	gpool := &WaitGo{token: make(chan struct{}, concurrency), concurrency: concurrency}

	for i := 0; i < concurrency; i++ {
		gpool.token <- struct{}{}
	}

	return gpool
}

func (pool *WaitGo) Add(fn func()) {

	<-pool.token

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if pool.PanicHandler != nil {
					pool.PanicHandler(r)
				} else {
					log.Printf("task paniced: %s", r)
				}
			}
			pool.token <- struct{}{}
		}()
		fn()
	}()
}

func (pool *WaitGo) Wait() {
	for i := 0; i < pool.concurrency; i++ {
		<-pool.token
	}

	close(pool.token)
}
