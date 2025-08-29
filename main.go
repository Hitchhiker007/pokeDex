package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	var myslice []string
	text = strings.ToLower(text)
	words := strings.Fields(text)

	for _, word := range words {
		myslice = append(myslice, word)
	}
	return myslice
}
