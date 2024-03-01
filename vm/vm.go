package vm

import (
	"monkey-int/bytecode"
	"monkey-int/compiler"
	"monkey-int/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions bytecode.Instructions

	stack []object.Object
	sp    int // stackpointer
}

func New(myBytecode *compiler.MyBytecode) *VM {
	return &VM{
		instructions: myBytecode.Instructions,
		constants:    myBytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}
