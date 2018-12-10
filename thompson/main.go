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
x = [ a b name@xxxx *( d e f ) g h i ];
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
}
