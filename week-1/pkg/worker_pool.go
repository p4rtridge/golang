package pkg

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

const (
	MAX_WORKER  uint8 = 10
	MAX_CHANNEL uint8 = 10
	JOB_NUM     int   = 100
)

type Result struct {
	workerId  uint8
	ranNum    int
	factorial int
}

type Worker struct {
	taskQueue <-chan int
	resultCh  chan<- Result
	id        uint8
}

func (w *Worker) start(wg *sync.WaitGroup) {
	for num := range w.taskQueue {
		result, error := calcFactorial(num)
		if error == nil {
			w.resultCh <- Result{workerId: w.id, ranNum: num, factorial: result}
		}
	}

	wg.Done()
}

type WorkerPool struct {
	taskQueue chan int
	resultCh  chan Result
	maxWorker uint8
}

func (wp *WorkerPool) start() {
	var wg sync.WaitGroup

	for i := uint8(0); i < wp.maxWorker; i++ {
		wg.Add(1)
		worker := Worker{id: i, taskQueue: wp.taskQueue, resultCh: wp.resultCh}
		go worker.start(&wg)
	}

	wg.Wait()
	close(wp.resultCh)
}

func (wp *WorkerPool) sendJobs(jobs []int) {
	for _, job := range jobs {
		wp.add(job)
	}

	close(wp.taskQueue)
}

func (wp *WorkerPool) add(num int) {
	wp.taskQueue <- num
}

func (wp *WorkerPool) printResult(done chan<- bool) {
	for result := range wp.resultCh {
		fmt.Printf("[Worker ID: %d] Num: %d, Factorial: %d\n", []any{int(result.workerId), result.ranNum, result.factorial}...)
	}
	done <- true
}

func Workerpool() {
	jobs := createJobs(JOB_NUM)

	workerPool := WorkerPool{
		maxWorker: MAX_WORKER,
		taskQueue: make(chan int, MAX_CHANNEL),
		resultCh:  make(chan Result, MAX_CHANNEL),
	}

	done := make(chan bool)

	go workerPool.sendJobs(jobs)
	go workerPool.printResult(done)
	workerPool.start()

	<-done
}

func createJobs(jobNum int) []int {
	jobs := []int{}

	for i := 0; i < jobNum; i++ {
		jobs = append(jobs, rand.Intn(10))
	}

	return jobs
}

func calcFactorial(num int) (int, error) {
	if num < 0 {
		return 0, errors.New("undefined")
	}

	if num == 0 {
		return 1, nil
	}

	result := 1
	for i := 2; i <= num; i++ {
		result *= i
	}

	return result, nil
}
