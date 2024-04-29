package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	fmt.Println("Enter text: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	letterCount := make(map[rune]int)

	for _, char := range text {
		char = unicode.ToLower(char)
		if unicode.IsLetter(char) {
			letterCount[char]++
		}
	}

	for char, count := range letterCount {
		percentage := float64(count) / float64(len(text)-1) * 100
		fmt.Printf("%c - %d %.1f%%\n", char, count, percentage)
	}
}
