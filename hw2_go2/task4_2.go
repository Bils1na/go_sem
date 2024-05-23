package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := "messages.txt"

	// Читаем данные из файла
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Файл не существует.")
		} else {
			fmt.Println("Ошибка при чтении файла:", err)
		}
		return
	}

	// Проверяем, пустой ли файл
	if len(data) == 0 {
		fmt.Println("Файл пуст.")
		return
	}

	// Выводим содержимое файла
	fmt.Println("Содержимое файла:")
	fmt.Println(string(data))
}