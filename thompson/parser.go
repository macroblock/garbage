package main

import (
	"fmt"

	"github.com/macroblock/garbage/thompson/errors"
	"github.com/macroblock/garbage/thompson/iterator"
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
				// variable.typ = strings.TrimSuffix(nodeType, "Decl")
				it := iterator.New(node, o.parser)
				n := 0
				// fmt.Printf("x node name %q\n%v\n", it.Name(), it)
				it.FindNext("lval")
				// fmt.Printf("y node name %q\n%v\n", it.Name(), it)
				it.FindNext("lval").ForEach(func(it *iterator.TNodeIterator) {
					n++
					fmt.Printf("cycle: %v\n", n)
					it.Enter()
					// fmt.Printf("---%v\n", o.parser.ByID(it.Node().Type))
					err := symbols.Add(
						it.FindNext("ident").Value(),
						it.FindNext("string").Value(),
					)
					errors.Add(err)
					// fmt.Printf("%q %q\n", v.name, v.label)
				})
				it.FindNext("options")
				it.FindNext("sequence")
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
