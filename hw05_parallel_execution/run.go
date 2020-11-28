package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Функция выполнения задачи из очереди.
func worker(tasksChn <-chan Task, tasksStatusChn chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// Запрашиваем задачу
	tasksStatusChn <- nil

	// Обрабатываем канал входящих задач
	for task := range tasksChn {
		if task == nil { // команда прекратить обработку
			return
		}

		res := task()
		tasksStatusChn <- res
	}
}

// Функция для параллельного выполнения заданий переданным числом горутин.
func Run(tasks []Task, workersCount int, maxErrors int) (result error) {
	var tasksChn = make(chan Task)
	var tasksStatusChn = make(chan error)

	tasksFailed := 0
	tasksToRun := len(tasks)
	if workersCount > tasksToRun {
		workersCount = tasksToRun
	}

	nextTaskNo := 0 // Номер следующей задачи на обработку
	workInProgress := workersCount

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go worker(tasksChn, tasksStatusChn, &wg)
	}

	// Обрабатываем канал состояния задач
	for err := range tasksStatusChn {
		if err != nil {
			if tasksFailed == maxErrors {
				// Достигнут максимум ошибок, а он может быть и равен 0
				result = ErrErrorsLimitExceeded
				// Больше новых задач не раздаём
				nextTaskNo = tasksToRun
			}

			tasksFailed++
		}

		// nil используется в качестве команды остановки обработчика - пропускаем такие значения
		for nextTaskNo < tasksToRun && tasks[nextTaskNo] == nil {
			nextTaskNo++
		}

		if nextTaskNo < tasksToRun {
			tasksChn <- tasks[nextTaskNo]
			nextTaskNo++
		} else {
			// Команда на прекращение работы
			tasksChn <- nil
			workInProgress--

			// Всем отправлены команда на прекращение работы
			if workInProgress == 0 {
				break
			}
		}
	}

	// Ждём выхода из всех функция обработки
	wg.Wait()

	return
}
