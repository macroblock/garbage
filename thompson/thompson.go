package main

import (
	"fmt"

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

	opIdent
	opRange
	opString
)

// TState -
type TState struct {
	// cmd      tOp
	element  IElem
	out      []*TState
	lastList int
}

// TFrag -
type TFrag struct {
	first *TState
	last  []*TState
}

// String -
func (o *TState) String() string {
	// ret := fmt.Sprintf("%T:%v\n", o.element, o.element)
	// if len(o.out) == 0 {
	// 	ret += "<finish>"
	// }
	// for _, v := range o.out {
	// 	ret += fmt.Sprintf("%v\n<cr>", v)
	// }
	ret := fmt.Sprintf("%v", o.Name())
	// ret := fmt.Sprintf("%v:%v", o.Name(), len(o.out))
	return ret
}

// Name -
func (o *TState) Name() string {
	// ret := fmt.Sprintf("%v:%v", o.element.Name(), len(o.out))
	ret := fmt.Sprintf("%v", o.element.Name())
	// if len(o.out) > 0 {
	// 	ret = fmt.Sprintf("%v>%v", ret, o.out[0].element.Name())
	// }
	return ret
}

// NewState -
func NewState(elem IElem) *TState {
	ret := &TState{element: elem}
	return ret
}

// NewFrag -
func NewFrag(state *TState) *TFrag {
	ret := &TFrag{}
	ret.first = state
	ret.last = append(ret.last, state)
	return ret
}

// ResetFrag -
func ResetFrag(frag *TFrag, state *TState) {
	frag.first = state
	frag.last = append(frag.last, state)
}

func appendFrag(frag *TFrag, fr *TFrag) *TFrag {
	for _, state := range frag.last {
		log.Errorf(len(state.out) != 0, "append1 frag")
	}
	for _, state := range fr.last {
		log.Errorf(len(state.out) != 0, "append2 frag")
	}

	for _, state := range frag.last {
		state.out = append(state.out, fr.first)
	}
	frag.last = fr.last
	return frag
}

func expandFrag(frag *TFrag, frs []*TFrag) *TFrag {
	for _, state := range frag.last {
		log.Errorf(len(state.out) != 0, "expand1 frag")
	}
	for _, fr := range frs {
		for _, state := range fr.last {
			log.Errorf(len(state.out) != 0, "expand2 frag")
		}
	}

	last := []*TState{}
	for _, state := range frag.last {
		for _, fr := range frs {
			state.out = append(state.out, fr.first)
			last = append(last, fr.last...)
		}
	}
	frag.last = last
	return frag
}

// Thompson -
func Thompson(element interface{}) (*TFrag, []error) {
	errors := errors.NewErrors()
	frag := NewFrag(NewState(nil))
	switch t := element.(type) {
	default:
		errors.Addf("thompson: an unsupported element type %T", t)
	case *TSequence:
		frag.first.element = (*TSequence)(nil)
		if len(t.elements) == 0 {
			log.Warningf(true, "sequence: len 0")
			return nil, nil
		}
		for _, v := range t.elements {
			fmt.Println("#### trying to append:", v)
			fr, errs := Thompson(v)
			errors.Add(errs...)
			if errs != nil || fr == nil {
				continue
			}
			frag = appendFrag(frag, fr)
		}
	case *TSplit:
		frag.first.element = (*TSplit)(nil)
		if len(t.elements) == 0 {
			log.Warningf(true, "split: len 0")
			return nil, nil
		}
		frs := []*TFrag{}
		for _, v := range t.elements {
			fr, errs := Thompson(v)
			errors.Add(errs...)
			if errs != nil || fr == nil {
				continue
			}
			frs = append(frs, fr)
		}
		frag = expandFrag(frag, frs)
	case *TKeepValue:
		errors.Addf("KeepValue is not supported yet")
	case *TIdent:
		ResetFrag(frag, NewState(t))
	case *TRange:
		ResetFrag(frag, NewState(t))
	case *TRune:
		ResetFrag(frag, NewState(t))
	case *TString:
		ResetFrag(frag, NewState(t))
	case *TKeepNode:
		errors.Addf("KeepNode is not supported yet")
	}
	return frag, errors.Get()

}
