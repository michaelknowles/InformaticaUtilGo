// the lexer is used to extract tokens from an input string

package expression

import (
	"fmt"
	"strings"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

// Literals are tokens representing literal strings
var Literals []string

// Keywords tokens
var Keywords []string

// Tokens including both Literals and Keywords
var Tokens []string

// TokenIds map token name -> int id
var TokenIds map[string]int

// Lexer object to create a scanner
var Lexer *lexmachine.Lexer

// scanInput runs the lexer on the input and returns a slice of tokens
func scanInput(input []byte) (tokens []*lexmachine.Token, err error) {
	scanner, err := Lexer.Scanner(input)
	if err != nil {
		return
	}

	for tok, sErr, eos := scanner.Next(); !eos; tok, sErr, eos = scanner.Next() {
		if _, is := sErr.(*machines.UnconsumedInput); is {
			err = sErr
			return
		} else if sErr != nil {
			err = sErr
			return
		}
		token := tok.(*lexmachine.Token)
		tokens = append(tokens, token)
	}

	return
}

// tokenTypeName takes the token.Type and returns the corresponding string
func tokenTypeName(token *lexmachine.Token) (tokenType string) {
	tokenType = Tokens[token.Type]
	return
}

// Called only once at package initialization
func init() {
	initTokens()
	var err error
	Lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
}

func initTokens() {
	Literals = []string{
		// in order of operator precedence
		"(",
		")",
		"+",
		"-",
		"NOT",
		"*",
		"/",
		"%",
		"||",
		"<",
		"<=",
		">",
		">=",
		"=",
		"<>", // not equal
		"!=", // not equal
		"^=", // not equal
		"AND",
		"OR",
		",",
	}
	Keywords = []string{
		":EXT",
		":INFA",
		":LKP",
		":MCR",
		":SD",
		":SEQ",
		":SP",
		":TD",
		"AND",
		"DD_DELETE", // 2
		"DD_INSERT", // 0
		"DD_REJECT", // 3
		"DD_UPDATE", // 1
		"FALSE",     // 0
		"NOT",
		"NULL",
		"OR",
		"PROC_RESULT",
		"SESSTARTTIME",
		"SPOUTPUT",
		"SYSDATE",
		"TRUE", // 1
		// workflow expressions
		"WORKFLOWSTARTTIME",
		"ABORTED",
		"DISABLED",
		"FAILED",
		"NOTSTARTED",
		"STARTED",
		"STOPPED",
		"SUCCEEDED",
	}
	Tokens = []string{
		"PARAM",
		"STRING",
		"NUMBER",
		"IDENT",
	}
	Tokens = append(Tokens, Literals...)
	Tokens = append(Tokens, Keywords...)
	TokenIds = make(map[string]int)
	for i, tok := range Tokens {
		TokenIds[tok] = i
	}
}

func initLexer() (lexer *lexmachine.Lexer, err error) {
	lexer = lexmachine.NewLexer()

	// Add tokens to lexer
	// the lexer chooses which rule to use by:
	//   1. pattern which matches the longest prefix
	//   2. pattern which was defined first
	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	for _, name := range Keywords {
		lexer.Add([]byte(name), token(name))
	}

	// Tokens by regular expression
	// You can test here: https://regoio.herokuapp.com/

	// Comment
	lexer.Add([]byte(`(--|//)[^\n]*\n?`), skip)
	// Parameter/Variable
	lexer.Add([]byte(`(\$)+([a-z]|[A-Z]|[0-9]|_|\-)+`), token("PARAM"))
	// String
	// parameters can exist inside strings so we'll need to check for them later during parsing
	lexer.Add([]byte(`'`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		for tc := scan.TC; tc < len(scan.Text); tc++ {
			if scan.Text[tc] == '\'' && tc+1 < len(scan.Text) {
				m := string(scan.Text[scan.TC:tc])
				scan.TC = tc + 1 // move the scanner past the matched string
				return scan.Token(TokenIds["STRING"], m, match), nil
			}
		}
		return nil,
			fmt.Errorf("unclosed string starting at %d, (%d, %d)",
				match.TC, match.StartLine, match.StartColumn)
	})
	// Number
	lexer.Add([]byte(`-?[0-9]+(\.?[0-9]+)*`), token("NUMBER"))
	// Identifier
	// Because go doesn't support lookaheads completely, functions are also matched here
	lexer.Add([]byte(`([a-z]|[A-Z]|_|[0-9])+`), token("IDENT"))
	// Whitespace
	lexer.Add([]byte(`( |\t|\n|\r)+`), skip) // skip whitespace

	err = lexer.Compile()

	return
}

// token constructs a Token of the given token type by the token type's name
func token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (t interface{}, err error) {
		t = s.Token(TokenIds[name], string(m.Bytes), m)
		err = nil
		return
	}
}

// skip the match
func skip(*lexmachine.Scanner, *machines.Match) (t interface{}, err error) {
	return nil, nil
}
