package expression

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// map the function name to the Go implementation to call
var fns map[string]interface{}

func init() {
	fns = map[string]interface{}{
		"ABS":    abs,
		"CHR":    chr,
		"CONCAT": concat,
		"LTRIM":  ltrim,
		"RTRIM":  rtrim,
	}
}

// Evaluate will lex, parse, and finally evaluate the input and return the result
func Evaluate(input string, vars []Variable) (result string, err error) {
	node, err := parse([]byte(input), vars)
	if err != nil {
		return
	}

	node, err = evaluateNode(node)
	if err != nil {
		return
	}

	result = node.Args[0].(string)

	return
}

func evaluateNode(node Node) (result Node, err error) {
	// Get the function from the node
	if _, ok := fns[node.Exp]; !ok {
		err = fmt.Errorf("the function %s either is invalid or hasn't been implemented by this library", node.Exp)
		return
	}
	function := reflect.ValueOf(fns[node.Exp])

	// If any arg is a node, evaluate it
	for i, arg := range node.Args {
		if n, ok := arg.(Node); !ok {
			err = fmt.Errorf("expected Node but got '%s'\n%+v", reflect.TypeOf(n), n)
			return
		} else if ok {
			if n.Exp != "NUMBER" && n.Exp != "STRING" && n.Exp != "NULL" {
				node.Args[i], err = evaluateNode(arg.(Node))
				if err != nil {
					return
				}
			}
		}
	}

	// Convert the args to Value for usage with the reflected function
	in := make([]reflect.Value, len(node.Args))
	for k, arg := range node.Args {
		in[k] = reflect.ValueOf(arg)
	}

	// Call the function
	value := function.Call(in)
	if !value[1].IsNil() {
		err = value[1].Interface().(error)
		return
	}
	result = value[0].Interface().(Node)

	return
}

func abs(args ...Node) (result Node, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("incorrect number of arguments, %d, to ABS", len(args))
		return
	}

	if args[0].Exp == "NULL" {
		result.Exp = "NULL"
		result.Args = append(result.Args, "NULL")
		return
	}

	num, err := strconv.ParseFloat(args[0].Args[0].(string), 64)
	if err != nil {
		return
	}

	num = math.Abs(num)

	result.Exp = "NUMBER"
	result.Args = append(result.Args, fmt.Sprintf("%f", num))

	return
}

func chr(args ...Node) (result Node, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("incorrect number of arguments, %d, to CHR", len(args))
		return
	}

	if args[0].Exp == "NULL" {
		result.Exp = "NULL"
		result.Args = append(result.Args, "NULL")
		return
	}

	result.Exp = "STRING"

	f, err := strconv.ParseFloat(args[0].Args[0].(string), 64)
	if err != nil {
		return
	}
	ascii := int(f)

	result.Args = append(result.Args, string(rune(ascii)))
	return
}

func concat(args ...Node) (result Node, err error) {
	if len(args) != 2 {
		err = fmt.Errorf("incorrect number of arguments, %d, to CONCAT", len(args))
		return
	}

	result.Exp = "STRING"

	if args[0].Exp == "NULL" {
		result.Args = append(result.Args, args[1].Args[0].(string))
		return
	}
	if args[1].Exp == "NULL" {
		result.Args = append(result.Args, args[0].Args[0].(string))
		return
	}

	result.Args = append(result.Args, args[0].Args[0].(string)+args[1].Args[0].(string))
	return
}

func ltrim(args ...Node) (result Node, err error) {
	result.Exp = "STRING"

	switch len(args) {
	case 1:
		v := strings.TrimLeft(args[0].Args[0].(string), ` `)
		result.Args = append(result.Args, v)
	case 2:
		if args[0].Exp == "NULL" || args[1].Exp == "NULL" {
			result.Exp = "NULL"
			result.Args = append(result.Args, "NULL")
		} else {
			v := strings.TrimLeft(args[0].Args[0].(string), args[1].Args[0].(string))
			result.Args = append(result.Args, v)
		}
	default:
		err = fmt.Errorf("incorrect number of arguments, %d, to LTRIM", len(args))
	}

	return
}

func rtrim(args ...Node) (result Node, err error) {
	result.Exp = "STRING"

	switch len(args) {
	case 1:
		v := strings.TrimRight(args[0].Args[0].(string), ` `)
		result.Args = append(result.Args, v)
	case 2:
		if args[0].Exp == "NULL" || args[1].Exp == "NULL" {
			result.Exp = "NULL"
			result.Args = append(result.Args, "NULL")
		} else {
			v := strings.TrimRight(args[0].Args[0].(string), args[1].Args[0].(string))
			result.Args = append(result.Args, v)
		}
	default:
		err = fmt.Errorf("incorrect number of arguments, %d, to LTRIM", len(args))
	}
	return
}
