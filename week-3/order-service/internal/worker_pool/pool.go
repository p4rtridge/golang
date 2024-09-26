package workerpool

import (
	"context"
	"errors"
	"fmt"
	"order_service/services/auth/entity"
	"order_service/services/auth/usecase"
	"sync"
)

var (
	once     sync.Once
	wg       sync.WaitGroup
	jobQueue chan Job
	result   chan Result
)

type Job struct {
	Ctx     context.Context
	Data    interface{}
	Request string
}

type Result struct {
	Data interface{}
	Err  error
}

type WorkerPool struct {
	authUc usecase.AuthUseCase
}

func NewWorkerPool(numWorkers, queueSize int, authUc usecase.AuthUseCase) *WorkerPool {
	var wp WorkerPool

	once.Do(func() {
		jobQueue = make(chan Job, queueSize)
		result = make(chan Result, queueSize)

		wp = WorkerPool{
			authUc: authUc,
		}

		for i := 1; i <= numWorkers; i++ {
			wg.Add(1)

			go wp.worker(i, jobQueue, result, &wg)
		}
	})

	return &wp
}

func (pool *WorkerPool) AddJob(job Job) {
	jobQueue <- job
}

func (pool *WorkerPool) GetResult() <-chan Result {
	return result
}

func (pool *WorkerPool) Close() {
	close(jobQueue)
	close(result)
	wg.Wait()
}

func (pool *WorkerPool) worker(id int, jobs <-chan Job, result chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("[Worker %d]: processing job: %s\n", id, job.Request)

		switch job.Request {
		case "login":
			token, err := pool.processLogin(job)
			result <- Result{Data: token, Err: err}

		default:
			fmt.Println("cannot be found any process")
		}
	}
}

func (pool *WorkerPool) processLogin(job Job) (*entity.TokenResponse, error) {
	data, ok := job.Data.(entity.AuthLogin)
	if !ok {
		return nil, errors.New("cannot be parse")
	}

	return pool.authUc.Login(job.Ctx, data)
}
