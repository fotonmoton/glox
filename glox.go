package main

import (
	"bufio"
	"fmt"
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

	doRun := func(line []byte) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		gl.run(line)
	}

	for {
		print("> ")
		if !scanner.Scan() {
			break
		}
		doRun(scanner.Bytes())
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
	tokens, err := newScanner(source).scan()

	if err != nil {
		panic(err)
	}

	stmts, parseErrs := newParser(tokens).parse()

	if parseErrs != nil {
		panic(parseErrs)
	}

	fmt.Println(AstStringer{stmts: stmts})

	resolveErrs := newResolver(gl.Interpreter).resolveStmts(stmts...)

	if resolveErrs != nil {
		panic(resolveErrs)
	}

	interpreterErrs := gl.Interpreter.interpret(stmts)

	if interpreterErrs != nil {
		panic(interpreterErrs)
	}
}
