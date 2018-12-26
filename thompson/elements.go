package main

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	// IElem -
	IElem interface {
		String() string
		Name() string
	}

	// ISeq -
	ISeq interface {
		Append(IElem)
		Repeat(int, int, bool)
	}

	// TSequence -
	TSequence struct {
		TRepeat
		elements []IElem
	}

	// TSplit -
	TSplit struct {
		TRepeat
		elements []IElem
	}

	// TKeepValue -
	TKeepValue struct {
		TRepeat
		elements []IElem
	}

	// TIdent -
	TIdent struct {
		TRepeat
		name    string
		element IElem
	}

	// TRange -
	TRange struct {
		TRepeat
		from rune
		to   rune
	}

	// TRune -
	TRune struct {
		TRepeat
		r rune
	}

	// TString -
	TString struct {
		TRepeat
		str string
	}

	// TKeepNode -
	TKeepNode struct {
		name string
	}

	// TRepeat -
	TRepeat struct {
		from, to int
		lazy     bool
	}
)

// Append -
func (o *TSequence) Append(v IElem) { o.elements = appendElement(o.elements, v) }

// Append -
func (o *TSplit) Append(v IElem) { o.elements = appendElement(o.elements, v) }

// Append -
func (o *TKeepValue) Append(v IElem) { o.elements = appendElement(o.elements, v) }

// Repeat -
func (o *TRepeat) Repeat(from, to int, lazy bool) {
	o.from = from
	o.to = to
	o.lazy = lazy
}

// Name -
func (o *TSequence) Name() string {
	if o == nil {
		return fmt.Sprint("sequence nil")
	}
	return fmt.Sprintf("%v(%v)", o.TRepeat, "sequence")
}

// Name -
func (o *TSplit) Name() string {
	if o == nil {
		return fmt.Sprint("split nil")
	}
	return fmt.Sprintf("%v[%v]", o.TRepeat, "split")
}

// Name -
func (o *TKeepValue) Name() string {
	if o == nil {
		return fmt.Sprint("keep value nil")
	}
	return fmt.Sprintf("%v<%v>", o.TRepeat, "keep value")
}

// Name -
func (o *TIdent) Name() string { return fmt.Sprintf("%v%v", o.TRepeat, o.name) }

// Name -
func (o *TRange) Name() string {
	if o.from == o.to {
		return fmt.Sprintf("%v%q", o.TRepeat, o.from)
	}
	return fmt.Sprintf("%v%q-%q", o.TRepeat, o.from, o.to)
}

// Name -
func (o *TRune) Name() string { return fmt.Sprintf("%v%q", o.TRepeat, o.r) }

// Name -
func (o *TString) Name() string { return fmt.Sprintf("%v\"%v\"", o.TRepeat, o.str) }

// Name -
func (o *TKeepNode) Name() string { return fmt.Sprintf("@%v", o.name) }

// -----------------------------------------------------------------------

// String -
func (o *TSequence) String() string { return fmt.Sprintf("%v(%v)", o.TRepeat, elemsToStr(o.elements)) }

// String -
func (o *TSplit) String() string { return fmt.Sprintf("%v[%v]", o.TRepeat, elemsToStr(o.elements)) }

// String -
func (o *TKeepValue) String() string { return fmt.Sprintf("%v<%v>", o.TRepeat, elemsToStr(o.elements)) }

// String -
func (o *TIdent) String() string { return fmt.Sprintf("%v%v:%v", o.TRepeat, o.name, o.element) }

// String -
func (o *TRange) String() string {
	if o.from == o.to {
		return fmt.Sprintf("%v%q", o.TRepeat, o.from)
	}
	return fmt.Sprintf("%v%q-%q", o.TRepeat, o.from, o.to)
}

// String -
func (o *TRune) String() string { return fmt.Sprintf("%v%q", o.TRepeat, o.r) }

// String -
func (o *TString) String() string { return fmt.Sprintf("%v\"%v\"", o.TRepeat, o.str) }

// String -
func (o *TKeepNode) String() string { return fmt.Sprintf("@%v", o.name) }

// -----------------------------------------------------------------------

func appendElement(elements []IElem, v IElem) []IElem {
	if v == nil {
		fmt.Println("sequence append: v is nil")
		return elements
	}
	return append(elements, v)
}

func elemsToStr(elements []IElem) string {
	slice := []string{}
	for _, v := range elements {
		slice = append(slice, fmt.Sprintf("%v", v))
	}
	return strings.Join(slice, " ")
}

func repeatXtoStr(x int) string {
	ret := "inf"
	if x >= 0 {
		ret = strconv.Itoa(x)
	}
	return ret
}
func (o TRepeat) String() string {
	from := repeatXtoStr(o.from)
	lazy := "!"
	if o.lazy {
		lazy = "?"
	}
	if o.from == o.to {
		if o.from == 1 {
			return ""
		}
		return from + lazy
	}
	return fmt.Sprintf("%v-%v%v", from, repeatXtoStr(o.to), lazy)
}
