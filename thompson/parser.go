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
		entries  []*ptool.TNode
		resolved bool
		data     interface{}
	}

	// TSymbolTable -
	TSymbolTable struct {
		data map[string]*TVar
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

	var traverseExpr func(*ptool.TNode) *errors.TErrors
	var traverse func(*ptool.TNode) *errors.TErrors
	traverse = func(root *ptool.TNode) *errors.TErrors {
		errors := errors.NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("an unsupported element %q", nodeType)
				continue
			case "comment":
			case "nodeDecl", "blockDecl":
				it := NewNodeIterator(node, o.parser)
				it.Next("lval")
				it.ForEach(func(i *TNodeIterator) {
					// i.Enter()
					err := symbols.Add(
						// i.Next("ident").Value(),
						// i.Next("string").Value(),
						i.Find("var", "ident").Value(),
						i.Find("var", "string").Value(),
					)
					errors.Add(err)
				})
				it.Next("options")
				it.Next("sequence")
				e := traverseExpr(it.Node())
				errors.Add(e.Get()...)
				// lval := findNode(node, "lval")
				// options := findNode(node, "options")
				// sequence := findNode(node, "sequence")
				// _ = options
				// e := traverseExpr(sequence)
				// errors.Add(e.Get()...)
				// for _, nd := range lval.Links {
				// 	v := &TVar{
				// 		name:  getNodeValue(nd, "ident"),
				// 		label: getNodeValue(nd, "string"),
				// 	}
				// 	// fmt.Printf("%q, %q\n", v.name, v.label)
				// 	if _, exist := vars[v.name]; exist {
				// 		errors.Addf("duplicate identifier %v", v.name)
				// 		continue
				// 	}
				// 	vars[v.name] = v
				// }
			} // switch nodeType
		} // for _, node := range root.Links
		o.symbols = symbols
		return errors
	}

	traverseExpr = func(root *ptool.TNode) *errors.TErrors {
		return nil
	}

	errors := traverse(o.tree)
	return errors.Get()
}
