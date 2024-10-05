package main

import (
	"bufio"
	"log"
	"os"
)

type Glox struct {
	Interpreter *Interpreter
}

func main() {
	glox := &Glox{newInterpreter()}
	switch len(os.Args) {
	case 1:
		glox.runPrompt()
	case 2:
		glox.runFile(os.Args[1])
	default:
		println("Usage: glox [file]")
		os.Exit(1)
	}
}

func (gl *Glox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	for {
		print("> ")
		if !scanner.Scan() {
			break
		}
		gl.run(scanner.Bytes(), true)
	}
}

func (gl *Glox) runFile(path string) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	runErrors := gl.run(file, false)

	if len(runErrors) != 0 {
		for _, e := range runErrors {
			log.Print(e)
		}

		os.Exit(1)
	}

}

func (gl *Glox) run(source []byte, interactive bool) []error {
	tokens, err := newScanner(source).scan()

	if err != nil {
		return []error{err}
	}

	stmts, parseErrs := newParser(tokens).parse()

	if len(parseErrs) != 0 && !interactive {
		return parseErrs
	}

	println(AstStringer{}.String(stmts))

	return gl.Interpreter.interpret(stmts)
}
