package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
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

	try(err)

	run(file)
}

func run(source []byte) {
	tokens := newScanner(source).scan()

	ast := newParser(tokens).parse()

	println(AstStringer{}.String(ast))
	println(AstToRPN{}.String(ast))
}

func try(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
