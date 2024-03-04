package vm

import (
	"fmt"
	"monkey-int/bytecode"
	"monkey-int/compiler"
	"monkey-int/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions bytecode.Instructions

	stack []object.Object
	sp    int // stackpointer, always pointing to the next FREE slot in the stack
}

func New(myBytecode *compiler.MyBytecode) *VM {
	return &VM{
		instructions: myBytecode.Instructions,
		constants:    myBytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := bytecode.Opcode(vm.instructions[ip])

		switch op {
		case bytecode.OpConstant:
			constIndex := bytecode.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case bytecode.OpAdd:
			val1 := vm.pop()
			val2 := vm.pop()
			val1int1, _ := val1.(*object.Integer)
			val1int2, _ := val2.(*object.Integer)
			addVal := val1int1.Value + val1int2.Value
			vm.push(&object.Integer{Value: addVal})
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("Stack overflow!")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	returnVal := vm.stack[vm.sp-1]
	vm.sp--
	return returnVal
}
