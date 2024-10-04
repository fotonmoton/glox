package main

import (
	"fmt"
	"log"
)

var hadError = false
var hadRuntimeError = false

func printError(token Token, message string) {
	if token.typ == EOF {
		report(token.line, " at and", message)
	} else {
		report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), message)
	}
}

func report(line int, where string, message string) {
	log.Printf("[%d] Error %s: %s", line, where, message)
	hadError = true
}

func reportRuntimeError(token Token, message string) {
	log.Printf("[%d] Error: %s", token.line, message)
	hadRuntimeError = true
}
