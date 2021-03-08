package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Функция выполнения задачи из очереди.
func startWorker(tasksChn <-chan Task) <-chan error {
	chnOut := make(chan error)

	go func() {
		// Запрашиваем задачение
		chnOut <- nil
		defer close(chnOut)
		// Обрабатываем канал входящих задач
		for task := range tasksChn {
			if task == nil {
				chnOut <- nil
			} else {
				chnOut <- task()
			}
		}
	}()

	return chnOut
}

func fanIn(channels []<-chan error) <-chan error {
	var wg sync.WaitGroup
	multiplexedStream := make(chan error)

	multiplex := func(c <-chan error) {
		defer wg.Done()

		for i := range c {
			multiplexedStream <- i
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

// Функция для параллельного выполнения заданий переданным числом горутин.
func Run(tasks []Task, workersCount int, maxErrors int) (result error) {
	tasksToRun := len(tasks)

	if tasksToRun == 0 {
		return
	}

	if workersCount > tasksToRun {
		workersCount = tasksToRun
	}

	var tasksChn = make(chan Task)
	var errorsChn = make([]<-chan error, workersCount)

	for i := 0; i < workersCount; i++ {
		errorsChn[i] = startWorker(tasksChn)
	}

	tasksFailed := 0
	nextTaskNo := 0 // Номер следующей задачи на обработку

	// Обрабатываем канал состояния задач
	for err := range fanIn(errorsChn) {
		if err != nil {
			if tasksFailed == maxErrors {
				// Достигнут максимум ошибок, а он может быть равен и 0
				result = ErrErrorsLimitExceeded
				// Больше новых задач не раздаём
				nextTaskNo = tasksToRun
			}

			tasksFailed++
		}

		if nextTaskNo < tasksToRun {
			tasksChn <- tasks[nextTaskNo]
			nextTaskNo++
		}

		if nextTaskNo == tasksToRun && tasksChn != nil {
			close(tasksChn)
			tasksChn = nil
		}
	}

	return
}
