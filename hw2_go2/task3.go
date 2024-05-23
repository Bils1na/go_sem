package main

import (
	"fmt"
	"os"
)

func main() {
	// Имя файла
	filename := "readonly.txt"

	// Создаем файл с правами только для чтения
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0444)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	// Записываем начальные данные в файл
	_, err = file.WriteString("Это файл только для чтения.")
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	// Закрываем файл, чтобы изменить его права доступа
	file.Close()

	// Открываем файл снова в режиме только для чтения
	file, err = os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Проверяем, можно ли записать в файл
	_, err = file.WriteString("Попытка записи.")
	if err != nil {
		fmt.Println("Успешно проверено: запись в файл невозможна. Ошибка:", err)
	} else {
		fmt.Println("Не удалось проверить: запись в файл выполнена успешно.")
	}

	// Читаем данные из файла и выводим их
	content := make([]byte, 1024)
	n, err := file.Read(content)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}
	fmt.Println("Содержимое файла:", string(content[:n]))
}