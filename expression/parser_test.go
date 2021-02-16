package expression

import (
	"reflect"
	"testing"
)

func TestIsFunction(t *testing.T) {
	testCases := []struct {
		value  string
		expect bool
	}{
		{"IS_DATE", true},
		{"CEIL", true},
		{"abc", false},
		{"is_date", false},
	}

	for _, tc := range testCases {
		result := isFunction(tc.value)
		if result != tc.expect {
			t.Errorf("For %s, expected %v but got %v", tc.value, tc.expect, result)
		}
	}
}

func TestIsLiteral(t *testing.T) {
	testCases := []struct {
		value  string
		expect bool
	}{
		{"+", true},
		{"*", true},
		{"plus", false},
	}

	for _, tc := range testCases {
		result := isLiteral(tc.value)
		if result != tc.expect {
			t.Errorf("For %s, expected %v but got %v", tc.value, tc.expect, result)
		}
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		input  string
		expect Node
	}{
		{
			input: `ABS(1)`,
			expect: Node{
				Exp: "ABS",
				Args: []interface{}{
					Node{"NUMBER", []interface{}{"1.000000"}},
				},
			},
		},
		{
			input: `1`,
			expect: Node{
				Exp:  "NUMBER",
				Args: []interface{}{"1.000000"},
			},
		},
		{
			input: `1 + 2`,
			expect: Node{
				Exp: "+",
				Args: []interface{}{
					Node{"NUMBER", []interface{}{"1.000000"}},
					Node{"NUMBER", []interface{}{"2.000000"}},
				},
			},
		},
		{
			input: `1 + 2 + 3`,
			expect: Node{
				Exp: "+",
				Args: []interface{}{
					Node{
						Exp: "+",
						Args: []interface{}{
							Node{"NUMBER", []interface{}{"1.000000"}},
							Node{"NUMBER", []interface{}{"2.000000"}},
						},
					},
					Node{"NUMBER", []interface{}{"3.000000"}},
				},
			},
		},
		{
			input: `1 + 2 * 3`,
			expect: Node{
				Exp: "+",
				Args: []interface{}{
					Node{"NUMBER", []interface{}{"1.000000"}},
					Node{
						Exp: "*",
						Args: []interface{}{
							Node{"NUMBER", []interface{}{"2.000000"}},
							Node{"NUMBER", []interface{}{"3.000000"}},
						},
					},
				},
			},
		},
		{
			input: `8 + (5 - 2) * 8`,
			expect: Node{
				Exp: "+",
				Args: []interface{}{
					Node{"NUMBER", []interface{}{"8.000000"}},
					Node{
						Exp: "*",
						Args: []interface{}{
							Node{
								Exp: "-",
								Args: []interface{}{
									Node{"NUMBER", []interface{}{"5.000000"}},
									Node{"NUMBER", []interface{}{"2.000000"}},
								},
							},
							Node{"NUMBER", []interface{}{"8.000000"}},
						},
					},
				},
			},
		},
		{
			input: `1 + ABS(-1)`,
			expect: Node{
				Exp: "+",
				Args: []interface{}{
					Node{"NUMBER", []interface{}{"1.000000"}},
					Node{
						Exp: "ABS",
						Args: []interface{}{
							Node{"NUMBER", []interface{}{"-1.000000"}},
						},
					},
				},
			},
		},
		{
			input: `RTRIM(' a ', ' ')`,
			expect: Node{
				Exp: "RTRIM",
				Args: []interface{}{
					Node{"STRING", []interface{}{` a `}},
					Node{"STRING", []interface{}{` `}},
				},
			},
		},
		{
			input: `IIF(1 + 2, ABS(5), 'b')`,
			expect: Node{
				Exp: "IIF",
				Args: []interface{}{
					Node{
						Exp: "+",
						Args: []interface{}{
							Node{"NUMBER", []interface{}{"1.000000"}},
							Node{"NUMBER", []interface{}{"2.000000"}},
						},
					},
					Node{
						Exp: "ABS",
						Args: []interface{}{
							Node{"NUMBER", []interface{}{"5.000000"}},
						},
					},
					Node{"STRING", []interface{}{"b"}},
				},
			},
		},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		node, err := parse([]byte(tc.input), vars)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.expect, node) {
			t.Errorf("Unexpected output\nExpected: \n%v\nGot: \n%v", tc.expect, node)
		}
	}
}

func TestParseVars(t *testing.T) {
	testCases := []struct {
		input  string
		vars   []Variable
		expect Node
	}{
		{
			input: `SUBSTR(LTRIM(RTRIM(in_AMT)), 1, 7)`,
			vars: []Variable{
				{"in_AMT", "NUMBER", "2"},
			},
			expect: Node{
				Exp: "SUBSTR",
				Args: []interface{}{
					Node{
						Exp: "LTRIM",
						Args: []interface{}{
							Node{
								Exp: "RTRIM",
								Args: []interface{}{
									Node{"NUMBER", []interface{}{"2.000000"}},
								},
							},
						},
					},
					Node{"NUMBER", []interface{}{"1.000000"}},
					Node{"NUMBER", []interface{}{"7.000000"}},
				},
			},
		},
	}

	for _, tc := range testCases {
		node, err := parse([]byte(tc.input), tc.vars)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.expect, node) {
			t.Errorf("Unexpected output\nExpected: \n%v\nGot: \n%v", tc.expect, node)
		}
	}
}
