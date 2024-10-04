package main

import (
	"bufio"
	"fmt"
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
		hadRuntimeError = false
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)

	try(err)

	run(file)

	switch {
	case hadError:
		os.Exit(65)
	case hadRuntimeError:
		os.Exit(70)
	default:
		os.Exit(0)
	}
}

func run(source []byte) {
	tokens := newScanner(source).scan()

	if hadError {
		return
	}

	ast := newParser(tokens).parse()

	if hadError {
		return
	}

	println(AstStringer{}.String(ast))

	res := newInterpreter().evaluate(ast)

	if hadRuntimeError {
		return
	}

	fmt.Printf("%v\n", res)
}

func try(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
