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

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction

	symbolTable *SymbolTable
}

type EmittedInstruction struct {
	Opcode   bytecode.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions:        bytecode.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
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
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(bytecode.OpBang)
		case "-":
			c.emit(bytecode.OpMinus)
		default:
			return fmt.Errorf("Unknown operator: %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		jnTruthyPos := c.emit(bytecode.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		jumpPos := c.emit(bytecode.OpJump, 9999)
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jnTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(bytecode.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(bytecode.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("Unknown symbol: %s", node.Value)
		}
		c.emit(bytecode.OpGetGlobal, symbol.Index)
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
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == bytecode.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := bytecode.Opcode(c.instructions[opPos])
	newInstruction := bytecode.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}
