package main

import (
	"fmt"
	"sync"
	"time"
)

type TType struct {
	ID          int
	CreatedAt   time.Time // Поменял строку на временную структуру
	CompletedAt time.Time
	TaskResult  string // Байты заменил на строку для удобства в определении выполнения задачи и дальнейшего вывода
}

const numWorkers = 10

func main() {
	tasks := make(chan TType, 10)
	doneTasks := make(chan TType, 10)
	errors := make(chan error, 10)

	var workersWg sync.WaitGroup // Waitgroup для отслеживания выполнения задач
	var createWg sync.WaitGroup  // Waitgroup для создания тасков

	// Генерация задач
	createWg.Add(1)
	go func() {
		defer createWg.Done()
		for i := 0; i < 10; i++ {
			createdAt := time.Now()
			task := TType{
				ID:        i,
				CreatedAt: createdAt,
			}
			tasks <- task
			time.Sleep(time.Millisecond * 10)
		}
	}()
	createWg.Wait()

	// Выполнение задач
	workersWg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer workersWg.Done()
			for task := range tasks {
				if task.CreatedAt.UnixNano()%2 > 0 {
					task.TaskResult = "unsuccessful"
					errors <- fmt.Errorf("Task ID : %d, Error : %s", task.ID, task.TaskResult)
				} else {
					task.TaskResult = "successful"
					task.CompletedAt = time.Now()
					doneTasks <- task
				}

				time.Sleep(time.Millisecond * 10)
			}
		}()
	}

	// Ожидание завершения выполнения всех задач, а дальше закрываем каналы для ошибок и успешных задач
	go func() {
		workersWg.Wait()
		close(errors)
		close(doneTasks)
	}()

	go func() {
		time.Sleep(5 * time.Millisecond)
		close(tasks)
	}()

	// Выводим все успешные задачи
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range doneTasks {
			fmt.Printf("Done task : %d at %s with result : %s\n", task.ID, task.CreatedAt.Format(time.RFC3339), task.TaskResult)
		}
	}()
	wg.Wait()

	// Таймаут для наглядности выполнения
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("timeout reached")
	}
}
