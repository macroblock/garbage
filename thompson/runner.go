package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/macroblock/imed/pkg/ptool"
)

const (
	// RuneEOF -
	RuneEOF = ptool.RuneEOF
)

type (
	// TNode -
	TNode struct {
		Links  []*TNode
		Values []string
		Type   int
		Label  string
	}
	// TRunner -
	TRunner struct {
		counters   []int
		elem       IElem
		prev       *TRunner
		nodeToKeep *TNode
		appendTo   *TNode
	}
)

// NewRunner -
func NewRunner(elem IElem, prev *TRunner, appendTo *TNode) *TRunner {
	ret := &TRunner{elem: elem, prev: prev, appendTo: appendTo}
	return ret
}

// Step -
func (o *TParser) Step(r rune) bool {
	// elems := o.curr.Step(&counters, r)
	// switch len(elems) {
	// case 0:
	// 	return nil
	// case 1:
	// 	o.curr = elems[0]
	// 	return []*TRunner{o}
	// }
	// split
	fmt.Printf("%q\n", r)
	ret := false
	for _, runner := range o.runners {
		switch t := runner.elem.(type) {
		default:
			log.Errorf(true, "unsupported element type %T", t)
		case *TSequence:
			runner.counters = append(runner.counters, 0)
			ret = o.Step(r)
		}
	}

	return ret
}

func nextRune(s string, pos int) (rune, int) {
	if pos > len(s)-1 {
		return RuneEOF, 2
	}
	return utf8.DecodeRuneInString(s[pos:])
}

// Run -
func (o *TParser) Run(text string) {
	o.runners = []*TRunner{NewRunner(o.start, nil, nil)}
	pos := 0
	for {
		r, w := nextRune(text, pos)
		if !o.Step(r) {
			break
		}

		pos += w
		if pos > len(text) {
			break
		}
	}
}
