package expression

import (
	"strings"
	"testing"
)

func TestTokenTypeName(t *testing.T) {
	input := strings.Join(Tokens, "")

	tokens, err := scanInput([]byte(input))
	if err != nil {
		t.Error(err)
	}

	for _, token := range tokens {
		if len(tokenTypeName(token)) == 0 {
			t.Errorf("Couldn't get type of token: %s", token.Lexeme)
		}
	}
}

func TestLiterals(t *testing.T) {
	input := strings.Join(Literals, " ")

	tokens, err := scanInput([]byte(input))
	if err != nil {
		t.Error(err)
	}

	expect := Literals

	if len(expect) != len(tokens) {
		t.Errorf("Got different number of tokens: %d instead of %d", len(tokens), len(expect))
	}

	for i, token := range tokens {
		if expect[i] != token.Value {
			t.Errorf("Expected: %s, got: %s", expect[i], token.Value)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := strings.Join(Keywords, " ")
	tokens, err := scanInput([]byte(input))
	if err != nil {
		t.Error(err)
	}

	expect := Keywords

	if len(expect) != len(tokens) {
		t.Errorf("Got different number of tokens: %d instead of %d", len(tokens), len(expect))
	}

	for i, token := range tokens {
		if expect[i] != token.Value {
			t.Errorf("Expected: %s, got: %s", expect[i], token.Value)
		}
	}

}

func TestComments(t *testing.T) {
	testCases := []string{
		`-- abscasdf`,
		`// ;laksjdf`,
	}

	for _, tc := range testCases {
		tokens, err := scanInput([]byte(tc))
		if err != nil {
			t.Error(err)
		}
		if len(tokens) > 0 {
			t.Error("Expected COMMENT to be skipped")
		}
	}
}

func TestSimpleExpression(t *testing.T) {
	input := []byte(`IIF(101 > 100.99, 'YES')`)

	tokens, err := scanInput(input)
	if err != nil {
		t.Error(err)
	}

	expect := []string{
		"IDENT",
		"(",
		"NUMBER",
		">",
		"NUMBER",
		",",
		"STRING",
		")",
	}

	for i, token := range tokens {
		tokenType := tokenTypeName(token)
		if tokenType != expect[i] {
			t.Errorf("Expected: %s, got: %s", expect[i], tokenType)
		}
	}
}
