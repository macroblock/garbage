package main

import "fmt"

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
		TSequence
	}

	// TKeepValue -
	TKeepValue struct {
		TSequence
	}

	// TIdent -
	TIdent struct {
		TRepeat
		name string
	}

	// TRange -
	TRange struct {
		TRepeat
		from rune
		to   rune
	}

	// TRune -
	TRune struct {
		repeat TRepeat
		r      rune
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
func (o *TSequence) Append(v interface{}) {
	if v == nil {
		fmt.Println("sequence append: v is nil")
		return
	}
	o.elements = append(o.elements, v)
}

// Repeat -
func (o *TRepeat) Repeat(from, to int) {
	o.from = from
	o.to = to
}
