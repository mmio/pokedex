package main

import (
	"strings"
)

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	tokens := strings.Fields(lowerText)
	
	return tokens
}

