package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	inputChan := make(chan int)
	squareChan := make(chan int)
	productChan := make(chan int)

	// Горутинa для ввода данных
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Введите число (или 'стоп' для выхода): ")
			scanner.Scan()
			text := scanner.Text()
			if strings.ToLower(text) == "стоп" {
				close(inputChan)
				return
			}
			num, err := strconv.Atoi(text)
			if err != nil {
				fmt.Println("Ошибка ввода. Пожалуйста, введите целое число.")
				continue
			}
			inputChan <- num
		}
	}()

	// Горутинa для вычисления квадрата числа
	go func() {
		for num := range inputChan {
			square := num * num
			fmt.Printf("Квадрат: %d\n", square)
			squareChan <- square
		}
		close(squareChan)
	}()

	// Горутинa для умножения квадрата на 2
	go func() {
		for square := range squareChan {
			product := square * 2
			fmt.Printf("Произведение: %d\n", product)
			productChan <- product
		}
		close(productChan)
	}()

	// Основной поток ожидает завершения всех горутин
	wg.Add(1)
	go func() {
		for range productChan {
			// просто считываем из канала, чтобы главный поток не завершился раньше времени
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Программа завершена.")
}
