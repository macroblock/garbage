package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

type (
	// TState -
	TState struct {
		op byte

		links    []*TState
		lastList int
	}

	// TFragment -
	TFragment struct {
		in  *TState
		out []*TState
	}

	// TStringStack -
	TStringStack []*string
)

// NewStringStack -
func NewStringStack() *TStringStack {
	return &TStringStack{}
}

// Push -
func (o *TStringStack) Push(val *string) {
	*o = append(*o, val)
}

// Pop -
func (o *TStringStack) Pop() *string {
	i := len(*o) - 1
	ret := (*o)[i]
	*o = (*o)[:i]
	return ret
}

var testProg = `
x += a b @c +(d y z) [ e f g ];
// comment here
z += asdf;
a 'test1', x 'test2', z 'test3'+= a b c;
`

func main() {
	if parser == nil {
		return
	}
	tree, err := parser.Parse(testProg)
	if err != nil {
		fmt.Println("\n*TParser.Parse error: ", err)
		return
	}
	fmt.Println(ptool.TreeToString(tree, parser.ByID))

	parser := NewParser()
	errors := parser.Parse(testProg)
	if errors != nil {
		print("parse error(s):", errors...)
		return
	}
	errors = parser.Build()
	if errors != nil {
		print("build error(s):", errors...)
		return
	}
}

func print(str string, errors ...error) {
	fmt.Println(str)
	for _, v := range errors {
		fmt.Println(v)
	}
}
