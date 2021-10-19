package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrWrongCounter = errors.New("ErrCounter less or equal zero")

type Task func() error

type Pool struct {
	poolSize int
	taskChan chan Task
	resChan  chan error
	done     chan struct{}

	wg sync.WaitGroup
}

func InitPool(poolsize int, tasksize int) *Pool {
	return &Pool{
		poolSize: poolsize,
		taskChan: make(chan Task, tasksize),
		resChan:  make(chan error),
		done:     make(chan struct{}),
	}
}

func (p *Pool) work() {
	defer p.wg.Done()
	for task := range p.taskChan {
		err := task()

		select {
		case <-p.done:
			return
		case p.resChan <- err:

		}
	}
}

func (p *Pool) RunWorkers() {
	p.wg.Add(p.poolSize)

	for i := 0; i < p.poolSize; i++ {
		go p.work()
	}

	p.wg.Wait()
	close(p.resChan)

}

func (p *Pool) CheckResult(result error, counter, size *int) error {
	if result != nil {
		*counter++
	}
	if *counter >= *size {
		close(p.done)
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrWrongCounter
	}
	if len(tasks) < n {
		n = len(tasks)
	}

	workPool := InitPool(n, 1)

	go workPool.RunWorkers()

	var err error
	errCount := 0

	for i := 0; i < len(tasks); {
		select {
		case workPool.taskChan <- tasks[i]:
			i++
		case res := <-workPool.resChan:
			res = workPool.CheckResult(res, &errCount, &m)
			if res != nil {
				err = res
			}
		}
	}

	close(workPool.taskChan)
	for res := range workPool.resChan {
		res = workPool.CheckResult(res, &errCount, &m)
		if res != nil {
			err = res
		}
	}

	close(workPool.done)

	workPool.wg.Wait()

	if err != nil {
		return err
	}
	return nil

}
