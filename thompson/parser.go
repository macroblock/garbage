package main

import (
	"fmt"
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
	if symbol, exists := o.data[v.name]; exists {
		if symbol.defined && v.defined {
			return fmt.Errorf("duplicate identifier %5v", v.name)
		}
		if !v.defined {
			fmt.Printf("skiped %5v %q %q\n", v.defined, v.name, v.label)
			return nil
		}
	}
	fmt.Printf("symbol has been added: %5v %q %q\n", v.defined, v.name, v.label)
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

func findNode(node *ptool.TNode, types ...string) *ptool.TNode {
	if node == nil {
		return nil
	}
	for _, typ := range types {
		tid := parser.ByName(typ)
		nextLevel := node
		node = nil
		for _, v := range nextLevel.Links {
			if v.Type == tid {
				node = v
				break
			}
		}
		if node == nil {
			return nil
		}
		nextLevel = node
	}
	return node
}

// Build -
func (o *TParser) Build() []error {
	symbols := NewSymbolTable()

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
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("@parseDecl: an unsupported element %q", nodeType)
				continue
				// case "comment":
			case "lval":
				variable = &TVar{defined: true}
				variable.name = node.Links[0].Value // "ident"
				if len(node.Links) > 1 {
					variable.label = node.Links[0].Value // "string"
				}
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

	for _, v := range symbols.data {
		fmt.Printf("--> %v %q, defined:%v\n", v.name, v.label, v.defined)
		if v.element != nil {
			fmt.Printf("element: %v\n", v.element)
		}
	}

	return errors.Get()
}
