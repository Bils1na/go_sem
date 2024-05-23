package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Имя файла
	filename := "messages.txt"

	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return
	}

	// Проверяем, пустой ли файл
	if fileInfo.Size() == 0 {
		fmt.Println("Файл пуст.")
		return
	}

	// Читаем и выводим строки из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Проверяем ошибки при чтении
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}
}