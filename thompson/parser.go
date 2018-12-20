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
	}

	// TVar -
	TVar struct {
		name     string
		label    string
		element  interface{}
		seqNode  *ptool.TNode
		options  *tOptions
		entries  []*ptool.TNode
		resolved bool
		defined  bool
		inUse    bool
		data     interface{}
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

// Build -
func (o *TParser) Build() []error {
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
			case "string":
			case "sequence":
				variable.seqNode = node
				elem := &TSequence{}
				elem.Repeat(1, 1)
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
		repeat := TRepeat{1, 1}
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("@parseVarsInSequence an unsupported element %q", nodeType)
				repeat.Repeat(1, 1)
				continue
			case "comment":
			case "repeat_01":
				repeat = TRepeat{0, 1}
			case "repeat_0f":
				repeat = TRepeat{0, -1}
			case "repeat_1f":
				repeat = TRepeat{1, -1}
			case "repeat_xy":
				from, err1 := strconv.Atoi(node.Links[0].Value)
				to, err2 := strconv.Atoi(node.Links[1].Value)
				errors.Add(err1, err2)
				repeat = TRepeat{from, to}
			case "repeat_xf":
				from, err := strconv.Atoi(node.Links[0].Value)
				errors.Add(err)
				repeat = TRepeat{from, -1}
			case "repeat_x":
				from, err := strconv.Atoi(node.Links[0].Value)
				errors.Add(err)
				repeat = TRepeat{from, from}
			case "sequence":
				elem := &TSequence{TRepeat: repeat}
				repeat.Repeat(1, 1)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "split":
				elem := &TSplit{TRepeat: repeat}
				repeat.Repeat(1, 1)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "keepValue":
				elem := &TKeepValue{TRepeat: repeat}
				repeat.Repeat(1, 1)
				errs := parsSequence(elem, node)
				errors.Add(errs.Get()...)
				element.Append(elem)
			case "keepNode":
				elem := &TKeepNode{}
				elem.name = node.Links[0].Value
				element.Append(elem)
			case "rune":
				elem := &TRune{TRepeat: repeat}
				repeat.Repeat(1, 1)
				elem.r, _ = utf8.DecodeRuneInString(node.Value)
				element.Append(elem)
			case "range":
				elem := &TRange{TRepeat: repeat}
				repeat.Repeat(1, 1)
				elem.from, _ = utf8.DecodeRuneInString(node.Links[0].Value)
				elem.to, _ = utf8.DecodeRuneInString(node.Links[1].Value)
				element.Append(elem)
			case "string":
				elem := &TString{TRepeat: repeat}
				repeat.Repeat(1, 1)
				elem.str = node.Value
				element.Append(elem)
			case "ident":
				err := symbols.Update(&TVar{name: node.Value, defined: false})
				errors.Add(err)
				elem := &TIdent{TRepeat: repeat}
				repeat.Repeat(1, 1)
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
	for _, name := range slice {
		v := symbols.data[name]
		fmt.Printf("--> %v %q, defined:%v\n", v.name, v.label, v.defined)
		if v.element != nil {
			fmt.Printf("element: %v\n", v.element)
		}
		if !v.defined {
			errors.Addf("udefined variable %v, %q", v.name, v.label)
		}
	}

	return errors.Get()
}
