package main

import (
	"fmt"

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
		options  *tOptions
		entries  []*ptool.TNode
		resolved bool
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

// NewVar -
func NewVar(name, label string, entry *ptool.TNode) *TVar {
	return &TVar{name: name, label: label, entries: []*ptool.TNode{entry}}
}

// NewSymbolTable -
func NewSymbolTable() *TSymbolTable {
	return &TSymbolTable{data: map[string]*TVar{}}
}

// Add -
func (o *TSymbolTable) Add(name, label string) error {
	if _, exists := o.data[name]; exists {
		return fmt.Errorf("duplicate identifier %v", name)
	}
	fmt.Printf("symbol has been added: %q %q\n", name, label)
	o.data[name] = NewVar(name, label, nil)
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
		parse     func(*ptool.TNode) *errors.TErrors
		parseDecl func(*ptool.TNode) *errors.TErrors
		parseLVal func(*ptool.TNode) *errors.TErrors
		parseExpr func(*ptool.TNode) *errors.TErrors
	)
	parse = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("an unsupported element %q", nodeType)
				continue
			case "comment":
			case "nodeDecl", "blockDecl":
				errs := parseDecl(node)
				errors.Add(errs.Get()...)
			}
		} // for _, node := range root.Links
		return errors
	}

	parseDecl = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("an unsupported element %q", nodeType)
				continue
			case "comment":
			case "lval":
				errs := parseLVal(node)
				errors.Add(errs.Get()...)
			case "options":
			case "sequence":
				errs := parseExpr(node)
				errors.Add(errs.Get()...)
			}
		}
		return errors
	}

	parseLVal = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("an unsupported element %q", nodeType)
				continue
			case "comment":
			case "var":
				it := NewNodeIterator(node, o.parser)
				name := it.Accept("ident").Value()
				label := ""
				if it.Try("string") {
					label = it.Value()
				}
				err := symbols.Add(name, label)
				errors.Add(err)
			}
		}
		return errors
	}

	parseExpr = func(root *ptool.TNode) *errors.TErrors {
		return nil
	}

	errors := parse(o.tree)
	return errors.Get()
}
