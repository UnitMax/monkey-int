package compiler

import (
	"monkey-int/ast"
	"monkey-int/bytecode"
	"monkey-int/object"
)

type Compiler struct {
	instructions bytecode.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) Bytecode() *MyBytecode {
	return &MyBytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type MyBytecode struct {
	Instructions bytecode.Instructions
	Constants    []object.Object
}
