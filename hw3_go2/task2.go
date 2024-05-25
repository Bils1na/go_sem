package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Канал для получения сигналов завершения работы
	sigChan := make(chan os.Signal, 1)
	// Канал для передачи состояния завершения горутинам
	done := make(chan bool)

	// Указываем, какие сигналы будем обрабатывать
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для обработки сигналов
	go func() {
		<-sigChan
		fmt.Println("Выхожу из программы")
		close(done)
	}()

	// Горутина для вывода квадратов натуральных чисел
	go func() {
		num := 1
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("Квадрат %d: %d\n", num, num*num)
				num++
				time.Sleep(1 * time.Second) // Задержка для имитации работы
			}
		}
	}()

	// Ожидание завершения программы
	<-done
	fmt.Println("Программа завершена.")
}
