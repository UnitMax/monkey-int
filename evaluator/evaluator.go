package evaluator

import (
	"fmt"
	"monkey-int/ast"
	"monkey-int/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, ctx *object.Context) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, ctx)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, ctx)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, ctx)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BlockStatement:
		return evalBlockStatement(node, ctx)
	case *ast.LetStatement:
		val := Eval(node.Value, ctx)
		if isError(val) {
			return val
		}
		ctx.Set(node.Name.Value, val)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, ctx)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, ctx)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, ctx)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, ctx)
	case *ast.Identifier:
		return evalIdentifier(node, ctx)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Ctx: ctx, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, ctx)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, ctx)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.ArrayLiteral:
		var els []object.Object
		for _, el := range node.Elements {
			evalEl := Eval(el, ctx)
			if isError(evalEl) {
				return evalEl
			}
			els = append(els, evalEl)
		}
		return &object.Array{Elements: els}
	case *ast.IndexExpression:
		leftEval := Eval(node.Left, ctx)
		arr, ok := leftEval.(*object.Array)
		if isError(leftEval) {
			return leftEval
		}
		if !ok {
			return &object.Error{Message: fmt.Sprintf("Index operator not supported for %q", leftEval.Type())}
		}
		indexEval := Eval(node.Index, ctx)
		if isError(indexEval) {
			return indexEval
		}
		indexNr, ok := indexEval.(*object.Integer)
		if !ok || int(indexNr.Value) >= len(arr.Elements) || int(indexNr.Value) < 0 {
			return NULL
			// return &object.Error{Message: "out of bounds array access"}
		}
		return arr.Elements[indexNr.Value]
	}
	return nil
}

func evalProgram(statements []ast.Statement, ctx *object.Context) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement, ctx)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}

		if errorValue, ok := result.(*object.Error); ok {
			return errorValue
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// check for integers before checking for booleans because booleans are only
	// evaluated with pointer comparison, not object comparison
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, ctx *object.Context) object.Object {
	condition := Eval(ie.Condition, ctx)
	if isError(condition) {
		return condition
	} else if isTruthy(condition) {
		return Eval(ie.Consequence, ctx)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, ctx)
	}
	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, ctx *object.Context) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, ctx)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, ctx *object.Context) object.Object {
	if value, ok := ctx.Get(node.Value); ok {
		return value
	}
	// look up built-in functions
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, ctx *object.Context) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, ctx)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if ok {
		extendedCtx := extendedFunctionCtx(function, args)
		evaluated := Eval(function.Body, extendedCtx)
		return unwrapReturnValue(evaluated)
	}
	builtin, ok := fn.(*object.Builtin)
	if ok {
		// no need to unwrap since builtins don't return the custom *object.ReturnValue type
		return builtin.Fn(args...)
	}
	return newError("Not a function: %s", fn.Type())
}

func unwrapReturnValue(obj object.Object) object.Object {
	// We don't want the wrapper here, we want the actual value
	// This is only necessary for return statements, expression statements already yield the value directly
	// e.g.:
	//   return x;  => Object.ReturnValue
	//   x;         => Object.Integer
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func extendedFunctionCtx(fn *object.Function, args []object.Object) *object.Context {
	ctx := object.NewEnclosedContext(fn.Ctx)
	for paramIdx, param := range fn.Parameters {
		ctx.Set(param.Value, args[paramIdx])
	}
	return ctx
}
