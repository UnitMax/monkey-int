package vm

import (
	"fmt"
	"monkey-int/bytecode"
	"monkey-int/compiler"
	"monkey-int/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

var VmTrue = &object.Boolean{Value: true}
var VmFalse = &object.Boolean{Value: false}
var VmNull = &object.Null{}

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int // stackpointer, always pointing to the next FREE slot in the stack

	globals []object.Object

	frames      []*Frame
	framesIndex int
}

func New(myBytecode *compiler.MyBytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: myBytecode.Instructions}
	mainFrame := NewFrame(mainFn)
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   myBytecode.Constants,
		stack:       make([]object.Object, StackSize),
		sp:          0,
		globals:     make([]object.Object, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(myBytecode *compiler.MyBytecode, s []object.Object) *VM {
	vm := New(myBytecode)
	vm.globals = s
	return vm
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	var ip int
	var ins bytecode.Instructions
	var op bytecode.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = bytecode.Opcode(ins[ip])

		switch op {
		case bytecode.OpPop:
			vm.pop()
		case bytecode.OpConstant:
			constIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
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
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case bytecode.OpMinus:
			val := vm.pop()
			if ival, ok := val.(*object.Integer); ok {
				vm.push(&object.Integer{Value: -ival.Value})
			} else {
				return fmt.Errorf("Unsupported type for negation: %s", val.Type())
			}
		case bytecode.OpJump:
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case bytecode.OpJumpNotTruthy:
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case bytecode.OpNull:
			err := vm.push(VmNull)
			if err != nil {
				return err
			}
		case bytecode.OpSetGlobal:
			globalIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case bytecode.OpGetGlobal:
			globalIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case bytecode.OpArray:
			numElements := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			err := vm.push(array)
			if err != nil {
				return err
			}
		case bytecode.OpHash:
			numElements := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)
			if err != nil {
				return err
			}
		case bytecode.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
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

func (vm *VM) executeBinaryOperation(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("Unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryStringOperation(op bytecode.Opcode, left object.Object, right object.Object) error {
	if op != bytecode.OpAdd {
		return fmt.Errorf("Unknown string operator: %d", op)
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return vm.push(&object.String{Value: leftVal + rightVal})
}

func (vm *VM) executeBinaryIntegerOperation(op bytecode.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64
	switch op {
	case bytecode.OpAdd:
		result = leftValue + rightValue
	case bytecode.OpSub:
		result = leftValue - rightValue
	case bytecode.OpMul:
		result = leftValue * rightValue
	case bytecode.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("Unknown integer operator: %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return VmTrue
	}
	return VmFalse
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()
	switch operand {
	case VmTrue:
		return vm.push(VmFalse)
	case VmFalse:
		return vm.push(VmTrue)
	case VmNull:
		return vm.push(VmTrue)
	default:
		return vm.push(VmFalse)
	}

}
func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]
		pair := object.HashPair{Key: key, Value: val}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("Invalid hash key: %s. Hashable not implemented.", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("Index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObj := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)
	if i < 0 || i > max {
		return vm.push(VmNull)
	}
	return vm.push(arrayObj.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("Unusuable as hash key: %s", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(VmNull)
	}
	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}
