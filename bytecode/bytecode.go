package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Opcode byte

type Instructions []byte

const (
	OpConstant      Opcode = 0x0A
	OpPop           Opcode = 0x0B
	OpAdd           Opcode = 0x0C
	OpSub           Opcode = 0x0D
	OpMul           Opcode = 0x0E
	OpDiv           Opcode = 0x0F
	OpFalse         Opcode = 0xA0
	OpTrue          Opcode = 0xA1
	OpNull          Opcode = 0xA2
	OpEqual         Opcode = 0xB0
	OpNotEqual      Opcode = 0xB1
	OpGreaterThan   Opcode = 0xB2
	OpLessThan      Opcode = 0xB3
	OpMinus         Opcode = 0xC0
	OpBang          Opcode = 0xC1
	OpJumpNotTruthy Opcode = 0xD0
	OpJump          Opcode = 0xD1
	OpCall          Opcode = 0xD2
	OpGetGlobal     Opcode = 0xE0
	OpSetGlobal     Opcode = 0xE1
	OpArray         Opcode = 0xF0
	OpHash          Opcode = 0xF1
	OpIndex         Opcode = 0xF2
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpPop:           {"OpPop", []int{}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpFalse", []int{}},
	OpFalse:         {"OpTrue", []int{}},
	OpNull:          {"OpNull", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpLessThan:      {"OpLessThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpCall:          {"OpCall", []int{}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	definition, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("Opcode %x undefined", op)
	}
	return definition, nil
}

func Make(op Opcode, operands ...int) []byte {
	definition, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instrunctionLen := 1
	for _, w := range definition.OperandWidths {
		instrunctionLen += w
	}

	instruction := make([]byte, instrunctionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := definition.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += (1 + read)
	}
	return out.String()
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
