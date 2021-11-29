package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrWrongCounter        = fmt.Errorf("wrong errors size")
)

type Task func() error

type Pool struct {
	errorCount *int32
	taskCh     chan Task
	wg         sync.WaitGroup
	maxError   int32
}

func (p *Pool) InitPool(errCount *int32, chSize int, maxErr int32) *Pool {
	return &Pool{
		errorCount: errCount,
		taskCh:     make(chan Task, chSize),
		wg:         sync.WaitGroup{},
		maxError:   maxErr,
	}
}

func (p *Pool) work() {
	defer p.wg.Done()

	for task := range p.taskCh {
		localCount := atomic.LoadInt32(p.errorCount)
		if localCount >= p.maxError {
			return
		}
		err := task()
		if err != nil {
			atomic.AddInt32(p.errorCount, 1)
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrWrongCounter
	}

	var errCounter int32
	var taskPool Pool
	pool := taskPool.InitPool(&errCounter, len(tasks), int32(m))

	for _, task := range tasks {
		pool.taskCh <- task
	}
	close(pool.taskCh)

	for i := 0; i < n; i++ {
		pool.wg.Add(1)
		go pool.work()
	}

	pool.wg.Wait()

	localErrorCount := int(*pool.errorCount)
	if localErrorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
