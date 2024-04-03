package vm

import (
	"monkey-int/bytecode"
	"monkey-int/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.fn.Instructions
}
