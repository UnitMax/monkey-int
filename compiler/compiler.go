package compiler

import (
	"fmt"
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
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpPop)
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		case "-":
			c.emit(bytecode.OpSub)
		case "/":
			c.emit(bytecode.OpDiv)
		case "*":
			c.emit(bytecode.OpMul)
		case ">":
			c.emit(bytecode.OpGreaterThan)
		case "<":
			c.emit(bytecode.OpLessThan)
		case "==":
			c.emit(bytecode.OpEqual)
		case "!=":
			c.emit(bytecode.OpNotEqual)
		default:
			return fmt.Errorf("Unknown operator: %s", node.Operator)
		}
	case *ast.Boolean:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(integer))
	}
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

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	ins := bytecode.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
