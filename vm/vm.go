package vm

import (
	"fmt"
	"monkey-int/bytecode"
	"monkey-int/compiler"
	"monkey-int/object"
)

const StackSize = 2048

var VmTrue = &object.Boolean{Value: true}
var VmFalse = &object.Boolean{Value: false}

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
		case bytecode.OpPop:
			vm.pop()
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
		case bytecode.OpSub:
			val1 := vm.pop()
			val2 := vm.pop()
			val1int1, _ := val1.(*object.Integer)
			val1int2, _ := val2.(*object.Integer)
			subVal := val1int2.Value - val1int1.Value
			vm.push(&object.Integer{Value: subVal})
		case bytecode.OpMul:
			val1 := vm.pop()
			val2 := vm.pop()
			val1int1, _ := val1.(*object.Integer)
			val1int2, _ := val2.(*object.Integer)
			multVal := val1int1.Value * val1int2.Value
			vm.push(&object.Integer{Value: multVal})
		case bytecode.OpDiv:
			val1 := vm.pop()
			val2 := vm.pop()
			val1int1, _ := val1.(*object.Integer)
			val1int2, _ := val2.(*object.Integer)
			divVal := val1int2.Value / val1int1.Value
			vm.push(&object.Integer{Value: divVal})
		case bytecode.OpFalse:
			err := vm.push(VmFalse)
			if err != nil {
				return err
			}
		case bytecode.OpTrue:
			err := vm.push(VmTrue)
			if err != nil {
				return err
			}
		case bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGreaterThan, bytecode.OpLessThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case bytecode.OpBang:
			val := vm.pop()
			if val == VmFalse {
				vm.push(VmTrue)
			} else { // everything that's literally true or just "truthy"
				vm.push(VmFalse)
			}
		case bytecode.OpMinus:
			val := vm.pop()
			if ival, ok := val.(*object.Integer); ok {
				vm.push(&object.Integer{Value: -ival.Value})
			} else {
				return fmt.Errorf("Unsupported type for negation: %s", val.Type())
			}
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

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp] // sp points to the next free element, so this is technically "free"
}

func (vm *VM) executeComparison(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftValue := left.(*object.Integer).Value
		rightValue := right.(*object.Integer).Value

		switch op {
		case bytecode.OpEqual:
			return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
		case bytecode.OpNotEqual:
			return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
		case bytecode.OpGreaterThan:
			return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
		case bytecode.OpLessThan:
			return vm.push(nativeBoolToBooleanObject(leftValue < rightValue))
		default:
			return fmt.Errorf("Unknown operator: %d", op)
		}
	}

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case bytecode.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("Unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return VmTrue
	}
	return VmFalse
}
