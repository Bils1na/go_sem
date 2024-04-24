package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	fmt.Println("Enter sequence")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	var wordLine, symbolNumber, bytesNumber int
	var prev rune

	for _, v := range text {
		if unicode.IsLetter(v) {
			symbolNumber++
		} else if unicode.IsSpace(v) && !unicode.IsSpace(prev) {
			wordLine++
		} else {
			symbolNumber++
		}
		prev = v
	}
	bytesNumber = len(text)

	fmt.Printf("Number of words: %d ", wordLine)
	fmt.Printf("Number of letters: %d ", symbolNumber)
	fmt.Printf("Number of bytes: %d ", bytesNumber)
}
