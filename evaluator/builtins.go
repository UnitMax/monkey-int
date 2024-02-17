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
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			}
			return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", len(args), 1)}
			}
			val, ok := args[0].(*object.Array)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to `first` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())}
			}
			if len(val.Elements) <= 0 {
				return NULL
			}
			return val.Elements[0]
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", len(args), 1)}
			}
			val, ok := args[0].(*object.Array)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to `last` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())}
			}
			if len(val.Elements) <= 0 {
				return NULL
			}
			return val.Elements[len(val.Elements)-1]
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", len(args), 2)}
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to `push` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())}
			}

			newArray := make([]object.Object, len(arr.Elements)+1)
			copy(newArray, arr.Elements)
			newArray[len(arr.Elements)] = args[1]
			return &object.Array{Elements: newArray}
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=%d", len(args), 1)}
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to `tail` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())}
			}
			if len(arr.Elements) <= 0 {
				return NULL
			}
			newArray := make([]object.Object, len(arr.Elements)-1)
			copy(newArray, arr.Elements[1:len(arr.Elements)])
			return &object.Array{Elements: newArray}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
