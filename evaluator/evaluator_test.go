package evaluator

import (
	"monkey-int/lexer"
	"monkey-int/object"
	"monkey-int/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	return Eval(p.ParseProgram())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an Integer. Got=%T (%+v) instead.", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. Got=%d, wanted=%d", result.Value, expected)
		return false
	}
	return true
}
