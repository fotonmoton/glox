package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	expr := &Binary{
		&Unary{Token{MINUS, "-", nil, 1}, &Literal{123}},
		Token{STAR, "*", nil, 1},
		&Grouping{&Grouping{&Binary{
			&Unary{Token{MINUS, "-", nil, 1}, &Literal{123}},
			Token{STAR, "*", nil, 1},
			&Grouping{&Grouping{&Literal{45.67}}}}}},
	}

	println(AstStringer{}.String(expr))

	switch len(os.Args) {
	case 1:
		runPrompt()
	case 2:
		runFile(os.Args[1])
	default:
		println("Usage: glox [file]")
		os.Exit(1)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for {
		print("> ")
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		run([]byte(scanner.Text()))
		hadError = false
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)

	panic(err)

	run(file)

	if hadError {
		os.Exit(1)
	}
}

func run(source []byte) {
	tokens := newScanner(source).scan()

	for _, token := range tokens {
		println(token.String())
	}
}

func panic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
