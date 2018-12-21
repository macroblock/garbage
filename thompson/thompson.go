package main

import (
	o "github.com/golang/freetype/example/round"
	"github.com/macroblock/garbage/thompson/errors"
)

type tOp int

const (
	_ = iota
	opMatchState
	opRune
	opAnd
	opOr
	opMaybe
	opStar
	opPlus
)

// TState -
type TState struct {
	cmd      tOp
	r        rune
	out      []*TState
	lastList int
}

// TFrag -
type TFrag struct {
	in  *TState
	out []**TState
}

func appendState(frag *TFrag, state *TState) {
	for i := range frag.out {
		frag.out[i] = state
	}
}

// Thompson -
func Thompson(element interface{}) (*TFrag, []error) {
	errors := errors.NewErrors()
	frag := &TFrag{}
	switch t := elem.(type) {
	default:
		errors.Addf("thompson: an unsupported element type %T", t)
	case *TSequence:
		if len(t.elements) == 0 {
			return nil, nil
		}
		fr, errs := Thompson(t.elements[0])
		errors.Add(errs...)
		if err != nil || in == nil {
			return err
		}
		fr.cmd = opAnd
		frag.in = in
		for _, v := range t.elements[1:] {
			fr, errs = Thompson(v)
			errors.Add(errs...)
			if errs != nil || fr == nil {
				continue
			}
			fr.cmd = opAnd
		}
	case *TSplit:
		for _, v := range t.elements {
			errs := o.resolveNodes(v)
			errors.Add(errs...)
		}
	case *TKeepValue:
		for _, v := range t.elements {
			errs := o.resolveNodes(v)
			errors.Add(errs...)
		}
	case *TIdent:
		if t.element == nil {
			v, err := o.symbols.Get(t.name)
			if err != nil {
				errors.Add(err)
				break
			}
			if v.element == nil {
				errors.Addf("something went wrong var %q", v.name)
				break
			}
			// errors.Addf("ident %q", v.name)
			t.element = v.element
			if !v.resolved {
				v.resolved = true
				errs := o.resolveNodes(v.element)
				errors.Add(errs...)
			}
		}
	case *TRange:
	case *TRune:
	case *TString:
	case *TKeepNode:
		errors.Addf("KeepNode is not supported yet")
	}
	return errors.Get()

	return nil, nil
}
