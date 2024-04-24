package main

import (
	"fmt"
	"log"
	"os"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func main() {
	var value int
	fmt.Println("Enter number:")
	if _, err := fmt.Fscan(os.Stdin, &value); err != nil {
		log.Fatal(err)
	}

	if isPrime(value) {
		fmt.Println("The number", value, "is simple")
	} else {
		fmt.Println("The number", value, "is not simple")
	}
}
