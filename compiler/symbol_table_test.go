package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("Expected a=%+v, got=%+v instead", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("Expected b=%+v, got=%+v instead", expected["b"], b)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("Name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("Expected %s to resolve to %+v, got=%+v instead.", sym.Name, sym, result)
		}
	}
}
