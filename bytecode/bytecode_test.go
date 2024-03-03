package bytecode

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		if len(instruction) != len(tt.expected) {
			t.Errorf("Instruction has wrong length. Wanted=%d, got=%d instead.", len(tt.expected), len(instruction))
		}

		for i := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("Wrong byte at pos %d. Wanted=0x%X, got=0x%X instead.", i, tt.expected[i], instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := "0000 OpAdd\n0001 OpConstant 2\n0004 OpConstant 65535\n"

	concatted := Instructions{}
	for _, instruction := range instructions {
		concatted = append(concatted, instruction...)
	}

	if concatted.String() != expected {
		t.Errorf("Instructions wrongly formatted.\nWanted=%q\ngot=%q instead.", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("Definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. Wanted=%d, got=%d instead.", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("Operand wrong. Wanted=%d, got=%q instead.", want, operandsRead[i])
			}
		}
	}
}
