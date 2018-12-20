package main

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	// ISeq -
	ISeq interface {
		Append(interface{})
		Repeat(int, int)
	}

	// TSequence -
	TSequence struct {
		TRepeat
		elements []interface{}
	}

	// TSplit -
	TSplit struct {
		TRepeat
		elements []interface{}
	}

	// TKeepValue -
	TKeepValue struct {
		TRepeat
		elements []interface{}
	}

	// TIdent -
	TIdent struct {
		TRepeat
		name    string
		element interface{}
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
	}
)

// Append -
func (o *TSequence) Append(v interface{}) { o.elements = appendElement(o.elements, v) }

// Append -
func (o *TSplit) Append(v interface{}) { o.elements = appendElement(o.elements, v) }

// Append -
func (o *TKeepValue) Append(v interface{}) { o.elements = appendElement(o.elements, v) }

// Repeat -
func (o *TRepeat) Repeat(from, to int) {
	o.from = from
	o.to = to
}

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

func appendElement(elements []interface{}, v interface{}) []interface{} {
	if v == nil {
		fmt.Println("sequence append: v is nil")
		return elements
	}
	return append(elements, v)
}

func elemsToStr(elements []interface{}) string {
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
	if o.from == o.to {
		if o.from == 1 {
			return ""
		}
		return from
	}
	return fmt.Sprintf("%v-%v", from, repeatXtoStr(o.to))
}
