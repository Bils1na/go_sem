package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func fillFileName(strs []string) string {
	var result string
	if len(strs) <= 1 {
		result += strs[0]
		return result
	}
	for i := 0; i <= len(strs)-2; i++ {
		if i == len(strs)-2 {
			result += strs[i]
		} else {
			result += strs[i] + "."
		}
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Укажите польный путь до файла вторым аргументом")
	}

	filePath := os.Args[1]

	var fileName, fileExt string

	pathEls := strings.Split(filePath, "/")
	lastEl := strings.Split(pathEls[len(pathEls)-1], ".")

	fileName = fillFileName(lastEl)
	if len(lastEl) > 1 {
		fileExt = lastEl[len(lastEl)-1]
	}

	fmt.Printf("filename: %s\n", fileName)
	fmt.Printf("extension: %s", fileExt)
}
