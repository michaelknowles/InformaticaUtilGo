# InformaticaUtilGo

This is a library for Informatica in Go. It is planned to have several packages to help working with Informatica, such
as expression testing, log parsing, and parameter searching.

**This is a work in progress**

## Expression

Install using:
```
go get -u "github.com/michaelknowles/informaticautilgo"
```

Import using:
```
import infa "github.com/michaelknowles/informaticautilgo/expression"
```

### Usage

This package is used to parse, validate, and evaluate Informatica expressions. You can use the Evaluate function:

```go
package main

import infa "github.com/michaelknowles/informaticautilgo/expression"

func main() {
	// You can use your parameters
	vars := []infa.Variable{
		{
			N: "$$Y",
			T: "STRING",
			V: "YES",
		},
		{
			N: "$$N",
			T: "STRING",
			V: "NO",
		},
	}
	result, err := infa.Evaluate("IIF(1 < 2, $$Y, $$N)", vars)
	
	if err != nil {
		fmt.Println("There was an error: ", err)
    } else {
		fmt.Println(result) // "YES"
    }
}

```