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
		gl.run(scanner.Bytes())
	}
}

func (gl *Glox) runFile(path string) {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	gl.run(file)
}

func (gl *Glox) run(source []byte) {
	tokens, _ := newScanner(source).scan()

	stmts, _ := newParser(tokens).parse()

	// fmt.Println(AstStringer{stmts: stmts})

	gl.Interpreter.interpret(stmts)
}
