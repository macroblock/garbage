package main

import (
	"github.com/macroblock/garbage/thompson/errors"
)

type tOp int

const (
	_ = iota
	opTemp
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
	element  IElem
	out      []*TState
	lastList int
}

// TFrag -
type TFrag struct {
	first *TState
	last  []*TState
}

// NewState -
func NewState(op tOp) *TState {
	ret := &TState{cmd: op, out: []*TState{nil}}
	return ret
}

// NewFrag -
func NewFrag(state *TState) *TFrag {
	ret := &TFrag{first: state}
	ret.last = append(ret.last, state)
	return ret
}

func appendFrag(frag *TFrag, fr *TFrag) *TFrag {
	// if frag == nil {
	// 	return fr
	// }
	for _, state := range frag.last {
		for i := range state.out {
			state.out[i] = fr.first
		}
	}
	frag.last = fr.last
	// frag.out = frag.out[:0]
	// for i := range fr.out {
	// 	frag.out = append(frag.out, fr.out[i])
	// }
	return frag
}

func expandFrag(frag *TFrag, fr []*TFrag) *TFrag {
	outList := []**TState{}
	frag.last = make([]**TState, len(fr))
	for i := range fr {
		// frag.out = append(frag.out, nil)
		*frag.last[i] = fr[i].first
		outList = append(outList, fr[i].last...)
	}
	frag.last = outList
	return frag
}

// Thompson -
func Thompson(element interface{}) (*TFrag, []error) {
	errors := errors.NewErrors()
	frag := NewFrag(NewState(opTemp))
	switch t := element.(type) {
	default:
		errors.Addf("thompson: an unsupported element type %T", t)
	case *TSequence:
		if len(t.elements) == 0 {
			return nil, nil
		}
		for _, v := range t.elements {
			fr, errs := Thompson(v)
			errors.Add(errs...)
			if errs != nil || fr == nil {
				continue
			}
			frag = appendFrag(frag, fr)
		}
	case *TSplit:
		frags := []*TFrag{}
		for _, v := range t.elements {
			fr, errs := Thompson(v)
			errors.Add(errs...)
			if errs != nil || fr == nil {
				continue
			}
			frags = append(frags, fr)
		}
		frag = expandFrag(frag, frags)
	case *TKeepValue:
		errors.Addf("KeepValue is not supported yet")
	case *TIdent:
		state := NewState(opAnd)
		frag.first = state
	case *TRange:
	case *TRune:
	case *TString:
	case *TKeepNode:
		errors.Addf("KeepNode is not supported yet")
	}
	return nil, errors.Get()

}
