package bytecode

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
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
