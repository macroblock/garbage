package main

import (
	"fmt"
	"sort"
	"strconv"
	"unicode/utf8"

	"github.com/macroblock/garbage/thompson/errors"
	"github.com/macroblock/imed/pkg/ptool"
)

type (
	// TParser -
	TParser struct {
		parser  *ptool.TParser
		tree    *ptool.TNode
		symbols *TSymbolTable

		start   IElem
		runners []*TRunner
	}

	// TVar -
	TVar struct {
		name      string
		label     string
		element   IElem
		seqNode   *ptool.TNode
		options   *tOptions
		entries   []*ptool.TNode
		resolved  bool
		defined   bool
		inUse     bool
		processed bool
		data      interface{}
	}

	// TSymbolTable -
	TSymbolTable struct {
		data map[string]*TVar
	}

	tOptions struct {
		skipSpace bool
		runeSize  int
		alwayKeep bool
	}

	tSequence struct {
	}
)

// NewSymbolTable -
func NewSymbolTable() *TSymbolTable {
	return &TSymbolTable{data: map[string]*TVar{}}
}

// Get -
func (o *TSymbolTable) Get(name string) (*TVar, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("get: variable name is empty")
	}
	ret := o.data[name]
	if ret == nil {
		return nil, fmt.Errorf("get: unknown symbol name %q", name)
	}
	return ret, nil
}

// Update -
func (o *TSymbolTable) Update(v *TVar) error {
	if v == nil {
		return fmt.Errorf("variable is nil")
	}
	if len(v.name) == 0 {
		return fmt.Errorf("variable name is empty")
	}
	inUse := "-"
	if v.inUse {
		inUse = "+"
	}
	if symbol, exists := o.data[v.name]; exists {
		if symbol.defined && v.defined {
			return fmt.Errorf("duplicate identifier %v %5v", inUse, v.name)
		}
		if !v.defined {
			fmt.Printf("skiped %v %5v %q %q\n", inUse, v.defined, v.name, v.label)
			return nil
		}
	}
	fmt.Printf("symbol has been added: %v %5v %q %q\n", inUse, v.defined, v.name, v.label)
	o.data[v.name] = v
	return nil
}

// NewParser -
func NewParser() *TParser {
	return &TParser{parser: parser}
}

// Tree -
func (o *TParser) Tree() string {
	if o.parser == nil || o.tree == nil {
		return "!!! not initialized !!!"
	}
	return ptool.TreeToString(o.tree, o.parser.ByID)
}

// Parse -
func (o *TParser) Parse(src string) []error {
	o.tree = nil
	tree, err := o.parser.Parse(testProg)
	if err != nil {
		return []error{err}
	}
	o.tree = tree
	return nil
}

func decodeRuneOrCode(s string) (rune, error) {
	ret := utf8.RuneError
	err := error(nil)
	switch utf8.RuneCountInString(s) {
	default:
		err = fmt.Errorf("decodeRuneOrCode: something went wrong")
	case 1:
		ret, _ = utf8.DecodeRuneInString(s)
	case 2:
		x, err := strconv.ParseInt(s, 16, 16)
		return rune(x), err
	}
	if ret == utf8.RuneError {
		err = fmt.Errorf("decodeRuneOrCode: rune error")
	}
	return ret, err
}

// Build -
func (o *TParser) Build() (IElem, []error) {
	if o.tree == nil {
		return nil, []error{fmt.Errorf("parse tree is <nil>")}
	}

	symbols := NewSymbolTable()

	useMode := false

	var (
		parse        func(*ptool.TNode) *errors.TErrors
		parseDecl    func(*ptool.TNode) *errors.TErrors
		parsSequence func(ISeq, *ptool.TNode) *errors.TErrors
		// parseLVal     func(*ptool.TNode) *errors.TErrors
	)
	parse = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("@parse: an unsupported element %q", nodeType)
				continue
			// case "comment":
			case "useBelowOn":
				useMode = true
			case "useBelowOff":
				useMode = false
			case "nodeDecl", "blockDecl":
				errs := parseDecl(node)
				errors.Add(errs.Get()...)
			}
		}
		return errors
	}

	parseDecl = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		variable := (*TVar)(nil)
		use := useMode
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("@parseDecl: an unsupported element %q", nodeType)
				continue
			// case "comment":
			case "useOn":
				use = true
			case "useOff":
				use = false
			case "useExclude":
				use = !use
			case "lval":
				variable = &TVar{defined: true, inUse: use}
				variable.name = node.Links[0].Value // "ident"
				if len(node.Links) > 1 {
					variable.label = node.Links[1].Value // "string"
				}
			case "string1", "string2":
			case "sequence":
				variable.seqNode = node
				elem := &TSequence{}
				elem.Repeat(1, 1, false)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				variable.element = elem
			}
		}
		err := symbols.Update(variable)
		errors.Add(err)
		return errors
	}

	parsSequence = func(element ISeq, root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		repeat := TRepeat{1, 1, false}
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("@parseVarsInSequence an unsupported element %q", nodeType)
				repeat.Repeat(1, 1, false)
				continue
			case "comment":
			case "repeat_01":
				repeat = TRepeat{0, 1, false}
			case "repeat_0f":
				repeat = TRepeat{0, -1, len(node.Links) != 0}
			case "repeat_1f":
				repeat = TRepeat{1, -1, len(node.Links) != 0}
			case "repeat_xy":
				from, err1 := strconv.Atoi(node.Links[0].Value)
				to, err2 := strconv.Atoi(node.Links[1].Value)
				errors.Add(err1, err2)
				repeat = TRepeat{from, to, len(node.Links) != 0}
			case "repeat_xf":
				from, err := strconv.Atoi(node.Links[0].Value)
				errors.Add(err)
				repeat = TRepeat{from, -1, len(node.Links) != 0}
			case "repeat_x":
				from, err := strconv.Atoi(node.Links[0].Value)
				errors.Add(err)
				repeat = TRepeat{from, from, len(node.Links) != 0}
			case "sequence":
				elem := &TSequence{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "split":
				elem := &TSplit{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "keepValue":
				elem := &TKeepValue{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "keepNode":
				elem := &TKeepNode{}
				elem.name = node.Links[0].Value
				element.Append(elem)
			case "rune", "runeCode":
				elem := &TRune{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				r, err := decodeRuneOrCode(node.Value)
				errors.Add(err)
				elem.r = r
				element.Append(elem)
			case "range":
				elem := &TRange{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				r1, err1 := decodeRuneOrCode(node.Links[0].Value)
				r2, err2 := decodeRuneOrCode(node.Links[1].Value)
				errors.Add(err1, err2)
				elem.from = r1
				elem.to = r2
				element.Append(elem)
			case "string1", "string2":
				elem := &TString{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				elem.str = node.Value
				element.Append(elem)
			case "ident":
				err := symbols.Update(&TVar{name: node.Value, defined: false})
				errors.Add(err)
				elem := &TIdent{TRepeat: repeat}
				repeat.Repeat(1, 1, false)
				elem.name = node.Value
				element.Append(elem)
			}
		}
		return errors
	}

	errors := parse(o.tree)

	slice := []string{}
	for _, v := range symbols.data {
		slice = append(slice, v.name)
	}
	sort.Strings(slice)

	o.symbols = symbols

	entries := []*TVar{}
	for _, name := range slice {
		v := o.symbols.data[name]
		fmt.Printf("--> %v %q, defined:%v\n", v.name, v.label, v.defined)
		if v.element != nil {
			fmt.Printf("element: %v\n", v.element)
		}
		if !v.defined {
			errors.Addf("udefined variable %v, %q", v.name, v.label)
		}
		if v.inUse {
			entries = append(entries, v)
			v.resolved = true
			errs := o.resolveNodes(v.element)
			if len(errs) == 0 {
				fmt.Printf("resolved: %v\n", v.element)
			}
			errors.Add(errs...)
		}
	}

	ret := IElem(nil)
	switch len(entries) {
	case 0:
	case 1:
		ret = entries[0].element
	default:
		split := &TSplit{}
		split.Repeat(1, 1, false)
		for _, v := range entries {
			split.Append(v.element)
		}
		ret = split
	}

	o.start = ret

	return ret, errors.Get()
}

func (o *TParser) resolveNodes(elem interface{}) []error {
	errors := errors.NewErrors()
	switch t := elem.(type) {
	default:
		errors.Addf("resolveNodes: an unsupported element type %T", t)
	case *TSequence:
		for _, v := range t.elements {
			errs := o.resolveNodes(v)
			errors.Add(errs...)
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
}
