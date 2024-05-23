package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
 // Открываем файл для записи, создаем его если он не существует, добавляем данные в конец файла
 	file, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
 	}
 	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 1

	for {
		fmt.Print("Введите сообщение: ")
		if !scanner.Scan() {
			break
		}
  		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Получаем текущую дату и время
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		// Формируем строку для записи в файл
		message := fmt.Sprintf("%d %s %s\n", lineNumber, currentTime, input)
		// Записываем строку в файл
		if _, err := file.WriteString(message); err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}

		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
	}
}