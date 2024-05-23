package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	filename := "messages.txt"

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
	// Читаем существующие данные из файла
	existingData, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	// Обновляем содержимое файла
	newData := append(existingData, []byte(message)...)
	// Записываем обновленные данные в файл
	err = ioutil.WriteFile(filename, newData, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
	}
}