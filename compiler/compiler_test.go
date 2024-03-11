package compiler

import (
	"fmt"
	"monkey-int/ast"
	"monkey-int/bytecode"
	"monkey-int/lexer"
	"monkey-int/object"
	"monkey-int/parser"
	"testing"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []bytecode.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1; 2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 + 2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 - 2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 * 2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 / 2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compiler error: %s", err)
		}

		myBytecode := compiler.Bytecode()
		err = testInstructions(tt.expectedInstructions, myBytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, myBytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []bytecode.Instructions, actual bytecode.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("Wrong instruction length.\nWanted=%q\ngot=%q instead.", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("Wrong instruction at %d.\nWanted=%q\ngot=%q instead.", i, concatted, actual)
		}
	}
	return nil
}

func concatInstructions(s []bytecode.Instructions) bytecode.Instructions {
	out := bytecode.Instructions{}
	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("Wrong number of constants. Wanted=%d, got=%d instead.", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("Constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not an Integer. Got=%T (%+v) instead.", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Object has wrong value. Wanted=%d, got=%d instead.", expected, result.Value)
	}

	return nil
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpLessThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (true) {10}; 3333;
			`,
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 7),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpPop),
				// 0008
				bytecode.Make(bytecode.OpConstant, 1),
				// 0011
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			if (true) {10;} else {20;}; 3333;
			`,
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []bytecode.Instructions{
				// 0000
				bytecode.Make(bytecode.OpTrue),
				// 0001
				bytecode.Make(bytecode.OpJumpNotTruthy, 10),
				// 0004
				bytecode.Make(bytecode.OpConstant, 0),
				// 0007
				bytecode.Make(bytecode.OpJump, 13),
				// 0010
				bytecode.Make(bytecode.OpConstant, 1),
				// 0013
				bytecode.Make(bytecode.OpPop),
				// 0014
				bytecode.Make(bytecode.OpConstant, 2),
				// 0017
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
