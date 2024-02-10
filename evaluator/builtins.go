package evaluator

import (
	"fmt"
	"monkey-int/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", len(args), 1)}
			}
			if val, ok := args[0].(*object.String); ok {
				return &object.Integer{Value: int64(len(val.Value))}
			}
			return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
		},
	},
}
