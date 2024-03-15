package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey-int/compiler"
	"monkey-int/evaluator"
	"monkey-int/lexer"
	"monkey-int/object"
	"monkey-int/parser"
	"monkey-int/vm"
	"os"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	useInterpreter := false
	if len(os.Args) >= 2 && os.Args[1] == "-int" {
		io.WriteString(out, "\nRunning in interpreter mode\n")
		useInterpreter = true
	} else {
		io.WriteString(out, "\nRunning in compiler mode\n")
	}

	scanner := bufio.NewScanner(in)

	// Evaluator context
	ctx := object.NewContext()

	// Compiler context
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, message := range p.Errors() {
				io.WriteString(out, "\t"+message+"\n")
			}
			continue
		}

		if useInterpreter {
			evaluated := evaluator.Eval(program, ctx)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}

		} else {
			comp := compiler.NewWithState(symbolTable, constants)
			err := comp.Compile(program)
			if err != nil {
				fmt.Fprintf(out, "Compilation error:\n %s\n", err)
				continue
			}

			code := comp.Bytecode()
			constants = code.Constants

			machine := vm.NewWithGlobalsStore(code, globals)
			err = machine.Run()
			if err != nil {
				fmt.Fprintf(out, "Executing bytecode failed:\n %s\n", err)
				continue
			}

			lastPopped := machine.LastPoppedStackElem()
			io.WriteString(out, lastPopped.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
