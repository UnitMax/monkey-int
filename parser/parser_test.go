package parser

import (
	"monkey-int/ast"
	"monkey-int/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 12345;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got=%v instead", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. Got=%q instead", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. Got=%T instead", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s'. Got=%s instead", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name is not '%s'. Got=%s instead", name, letStatement.Name)
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 123456;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. Got=%d instead.", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not *ast.returnStatement. Got=%T", statement)
			continue
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement.TokenLiteratl not 'return', got=%q instead.", returnStatement.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors!", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q.", msg)
	}
	t.FailNow()
}
