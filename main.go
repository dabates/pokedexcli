package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	stringVals := strings.Fields(text)
	return stringVals
}

func main() {
	fmt.Println(cleanInput("Hello, World!"))
}
