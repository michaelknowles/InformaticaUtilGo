// the parser takes tokens from the lexer and creates an AST

package expression

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/timtadh/lexmachine"
)

// all Informatica function names
// https://docs.informatica.com/data-integration/powercenter/10-4-0/transformation-language-reference/functions.html
var functions map[string]interface{}

func init() {
	// set the allowed function names
	functions = map[string]interface{}{
		"ABORT":            nil,
		"ABS":              nil,
		"ADD_TO_DATE":      nil,
		"AES_DECRYPT":      nil,
		"AES_ENCRYPT":      nil,
		"ASCII":            nil,
		"AVG":              nil,
		"BINARY_COMPARE":   nil,
		"BINARY_CONCAT":    nil,
		"BINARY_LENGTH":    nil,
		"BINARY_SECTION":   nil,
		"CEIL":             nil,
		"CHOOSE":           nil,
		"CHR":              nil,
		"CHRCODE":          nil,
		"COMPRESS":         nil,
		"CONCAT":           nil,
		"CONVERT_BASE":     nil,
		"COS":              nil,
		"COSH":             nil,
		"COUNT":            nil,
		"CRC32":            nil,
		"CUME":             nil,
		"DATE_COMPARE":     nil,
		"DATE_DIFF":        nil,
		"DEC_BASE64":       nil,
		"DEC_HEX":          nil,
		"DECODE":           nil,
		"DECOMPRESS":       nil,
		"EBCDIC_ISO88591":  nil,
		"ENC_BASE64":       nil,
		"ENC_HEX":          nil,
		"ERROR":            nil,
		"EXP":              nil,
		"FIRST":            nil,
		"FLOOR":            nil,
		"FV":               nil,
		"GET_DATE_PART":    nil,
		"GREATEST":         nil,
		"IIF":              nil,
		"IN":               nil,
		"INDEXOF":          nil,
		"INITCAP":          nil,
		"INSTR":            nil,
		"ISNULL":           nil,
		"IS_DATE":          nil,
		"IS_NUMBER":        nil,
		"IS_SPACES":        nil,
		"LAG":              nil,
		"LAST":             nil,
		"LAST_DAY":         nil,
		"LEAD":             nil,
		"LEAST":            nil,
		"LENGTH":           nil,
		"LN":               nil,
		"LOG":              nil,
		"LOOKUP":           nil,
		"LOWER":            nil,
		"LPAD":             nil,
		"LTRIM":            nil,
		"MAKE_DATE_TIME":   nil,
		"MAX":              nil, // Dates, Numbers, String
		"MD5":              nil,
		"MEDIAN":           nil,
		"METAPHONE":        nil,
		"MIN":              nil, // Dates, Numbers, String
		"MOD":              nil,
		"MOVINGAVG":        nil,
		"MOVINGSUM":        nil,
		"NPER":             nil,
		"PERCENTILE":       nil,
		"PMT":              nil,
		"POWER":            nil,
		"PV":               nil,
		"RAND":             nil,
		"RATE":             nil,
		"REG_EXTRACT":      nil,
		"REG_MATCH":        nil,
		"REG_REPLACE":      nil,
		"REPLACECHR":       nil,
		"REPLACESTR":       nil,
		"REVERSE":          nil,
		"ROUND":            nil, // Dates, Numbers
		"RPAD":             nil,
		"RTRIM":            nil,
		"SETCOUNTVARIABLE": nil,
		"SET_DATE_PART":    nil,
		"SETMAXVARIABLE":   nil,
		"SETMINVARIABLE":   nil,
		"SETVARIABLE":      nil,
		"SHA256":           nil,
		"SIGN":             nil,
		"SIN":              nil,
		"SINH":             nil,
		"SOUNDEX":          nil,
		"SQRT":             nil,
		"STDDEV":           nil,
		"SUBSTR":           nil,
		"SUM":              nil,
		"SYSTIMESTAMP":     nil,
		"TAN":              nil,
		"TANH":             nil,
		"TIME_RANGE":       nil,
		"TO_BIGINT":        nil,
		"TO_CHAR":          nil, // Dates, Number
		"TO_DATE":          nil,
		"TO_DECIMAL":       nil,
		"TO_FLOAT":         nil,
		"TO_INTEGER":       nil,
		"TRUNC":            nil, // Dates, Number
		"UPPER":            nil,
		"VARIANCE":         nil,
	}
}

// Node is a node in the AST; nodes may be nested
type Node struct {
	Exp  string
	Args []interface{}
}

// Variable is an IDENT that is substituted during parsing
type Variable struct {
	N string // name
	T string // type
	V string // value
}

// parse a string into an AST
func parse(input []byte, vars []Variable) (node Node, err error) {
	// Tokenize
	tokens, err := scanInput(input)
	if err != nil {
		return
	}

	// Convert to AST
	nodes, endPos, err := parseExpression(tokens, vars, 0)

	if err != nil {
		return
	}

	if len(nodes) > 1 {
		err = fmt.Errorf("couldn't flatten down to one node:\n%+v", nodes)
		return
	}

	if endPos < len(tokens)-1 {
		err = fmt.Errorf("Parsing tokens ended at %d, but expected %d", endPos, len(tokens)-1)
		return
	}

	node = nodes[0]

	return
}

// parse the expression starting at the given position (recursively) into a Node
func parseExpression(tokens []*lexmachine.Token, vars []Variable, pos int) (output []Node, endPos int, err error) {
	// Here we loop through the tokens multiple times and convert each piece to a (possibly nested) Node
	// Buffers are used to hold intermediate results between passes
	// Passes should be in order of operator precedence (e.g. * before +)
	buffer := make([]interface{}, 0)

	// First pass for parentheses, functions, and variable substitution
paren: // label used to escape for loop
	for pos < len(tokens) {
		tokenType := tokenTypeName(tokens[pos])
		value := tokens[pos].Value.(string)
		switch tokenType {
		case ",": // skip commas and move the position forward
			pos++
		case ")": // we hit the end of the current parenthesis
			break paren // escape for loop
		case "(": // start a new parenthesis, nested node here
			nodes, end, cErr := parseExpression(tokens, vars, pos+1)
			if cErr != nil {
				err = cErr
				return
			}
			for _, n := range nodes {
				buffer = append(buffer, n)
			}
			pos = end + 1
		case "IDENT":
			if isFunction(value) { // nested node here
				if tokens[pos+1].Value.(string) == "(" {
					node := Node{value, make([]interface{}, 0)}
					iNodes, end, cErr := parseExpression(tokens, vars, pos+2)
					if cErr != nil {
						err = cErr
						return
					}
					for _, n := range iNodes {
						node.Args = append(node.Args, n)
					}
					buffer = append(buffer, node)
					pos = end + 1
				} else {
					err = fmt.Errorf("expected '(' after %s", value)
					return
				}
			} else { // it's a Variable
				for _, v := range vars {
					if value == v.N {
						node := Node{v.T, make([]interface{}, 0)}
						if v.T == "NUMBER" {
							f, fErr := reformatFloat(v.V)
							if fErr != nil {
								err = fErr
								return
							}
							node.Args = append(node.Args, f)
						} else {
							node.Args = append(node.Args, v.V)
						}
						buffer = append(buffer, node)
						pos++
						goto paren // continue for loop at next token
					}
				}

				err = fmt.Errorf("the identifier '%s' was not found", value)
				return
			}
		case "NUMBER": // reformat the float (e.g. 2 => 2.000000)
			f, fErr := reformatFloat(value)
			if fErr != nil {
				err = fErr
				return
			}
			node := Node{tokenType, make([]interface{}, 0)}
			node.Args = append(node.Args, f)
			buffer = append(buffer, node)
			pos++
		case "STRING":
			// Replace any params/vars with their value
			for _, v := range vars {
				if strings.HasPrefix(v.N, "$") {
					value = strings.ReplaceAll(value, v.N, v.V)
				}
			}
			node := Node{tokenType, make([]interface{}, 0)}
			node.Args = append(node.Args, value)
			buffer = append(buffer, node)
			pos++
		default:
			// misc type that doesn't need further processing right now
			node := Node{tokenType, make([]interface{}, 0)}
			node.Args = append(node.Args, value)
			buffer = append(buffer, node)
			pos++
		}
	}
	// set the endPos now since we've done a pass thru the original tokens
	// further passes will be on the buffers and so the end positions will be irrelevant to the function caller
	endPos = pos

	// "*", "/", "%"
	bufferNew := make([]interface{}, 0)
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "*", "/", "%":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "+", "-"
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "+", "-":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "||"
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "||":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "<", "<=", ">", ">="
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "<", "<=", ">", ">=":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "=", "<>", "!=", "^="
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "=", "<>", "!=", "^=":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "AND"
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "AND":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// "OR"
	buffer = nil
	buffer = append(buffer, bufferNew...)
	bufferNew = nil
	pos = 0
	for pos < len(buffer) {
		switch buffer[pos].(type) {
		case Node:
			switch buffer[pos].(Node).Exp {
			case "OR":
				node := Node{
					Exp: buffer[pos].(Node).Exp,
					Args: []interface{}{
						bufferNew[len(bufferNew)-1],
						buffer[pos+1],
					},
				}
				bufferNew[len(bufferNew)-1] = node
				pos += 2
			default:
				bufferNew = append(bufferNew, buffer[pos])
				pos++
			}
		default:
			bufferNew = append(bufferNew, buffer[pos])
			pos++
		}
	}

	// validate each output is a Node
	// TODO this check shouldn't be needed but keeping for now while in heavy dev
	for _, b := range bufferNew {
		if n, ok := b.(Node); ok {
			output = append(output, n)
		} else {
			err = fmt.Errorf("exptected a Node but got: %s\n+%v", reflect.TypeOf(bufferNew[0]), bufferNew[0])
			return
		}
	}

	return
}

// reformatFloat takes a number and rewrites it as a float (e.g. 2 => 2.000000) in string format
func reformatFloat(i string) (o string, err error) {
	f, err := strconv.ParseFloat(i, 64)
	if err != nil {
		return
	}
	o = fmt.Sprintf("%f", f)

	return
}

// Check if the identifier is a Literal
func isLiteral(ident string) bool {
	for _, l := range Literals {
		if l == ident {
			return true
		}
	}

	return false
}

// Check if the identifier is a function name
func isFunction(ident string) bool {
	if _, ok := functions[ident]; !ok {
		return false
	}

	return true
}
