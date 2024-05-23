package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Функция для генерации комбинаций скобок
func generateParenthesis(n int) []string {
 result := []string{}
 generate(&result, "", 0, 0, n)
 return result
}

// Вспомогательная функция для генерации комбинаций скобок
func generate(result *[]string, current string, open, close, max int) {
	if len(current) == max*2 {
		*result = append(*result, current)
		return
 	}

	if open < max {
		generate(result, current+"(", open+1, close, max)
	}
	if close < open {
		generate(result, current+")", open, close+1, max)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите количество пар скобок: ")

	if scanner.Scan() {
		input := scanner.Text()
		n, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка: введите целое число")
			return
		}

		combinations := generateParenthesis(n)
		fmt.Println(combinations)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
	}
}