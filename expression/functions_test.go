package expression

import "testing"

func TestABS(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{`ABS(NULL)`, `NULL`},
		{`ABS(250)`, `250.000000`},
		{`ABS(-250)`, `250.000000`},
		{`ABS(1.1)`, `1.100000`},
		{`ABS(-1.1)`, `1.100000`},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		result, err := Evaluate(tc.input, vars)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expect {
			t.Errorf("Input: %s\nExpected: `%s`, got `%s`", tc.input, tc.expect, result)
		}
	}
}

func TestCHR(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{`CHR(NULL)`, `NULL`},
		{`CHR(39)`, `'`},
		{`CHR(84)`, `T`},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		result, err := Evaluate(tc.input, vars)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expect {
			t.Errorf("Input: %s\nExpected: `%s`, got `%s`", tc.input, tc.expect, result)
		}
	}
}

func TestCONCAT(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{`CONCAT('mike', 'knowles')`, `mikeknowles`},
		{`CONCAT(NULL, 'knowles')`, `knowles`},
		{`CONCAT('mike', NULL)`, `mike`},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		result, err := Evaluate(tc.input, vars)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expect {
			t.Errorf("Input: %s\nExpected: `%s`, got `%s`", tc.input, tc.expect, result)
		}
	}
}

func TestLTRIM(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{`LTRIM('H. Bender', 'S.')`, `H. Bender`},
		{`LTRIM(NULL)`, `NULL`},
		{`LTRIM(RTRIM(' a ')`, `a`},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		result, err := Evaluate(tc.input, vars)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expect {
			t.Errorf("Input: %s\nExpected: `%s`, got `%s`", tc.input, tc.expect, result)
		}
	}
}

func TestRTRIM(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{`RTRIM(' a ')`, ` a`},
		{`RTRIM(NULL)`, `NULL`},
	}

	vars := make([]Variable, 0)

	for _, tc := range testCases {
		result, err := Evaluate(tc.input, vars)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expect {
			t.Errorf("Input: %s\nExpected: `%s`, got `%s`", tc.input, tc.expect, result)
		}
	}
}
