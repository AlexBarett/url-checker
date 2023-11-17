package workerpool

import "fmt"

type WorkerPool chan struct{}

func New(limit int) WorkerPool {
	return make(chan struct{}, limit)
}

func (wp WorkerPool) Exec(work func()) error {
	errChan := make(chan error, 1)
	wp <- struct{}{}

	go func() {
		defer func() {
			r := recover()
			if r != nil {
				errChan <- fmt.Errorf("worker pool recover: %v", r)
			}
		}()

		work()
		close(errChan)
		<-wp
	}()

	return <-errChan
}
