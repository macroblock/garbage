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
		parser *ptool.TParser
		tree   *ptool.TNode
		vars   map[string]*TVar
	}

	// TVar -
	TVar struct {
		name     string
		label    string
		entries  []*ptool.TNode
		resolved bool
		data     interface{}
	}
)

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

func getValue(node *ptool.TNode) string {
	if node != nil {
		return node.Value
	}
	return ""
}

func getNodeValue(node *ptool.TNode, types ...string) string {
	return getValue(findNode(node, types...))
}

// Build -
func (o *TParser) Build() []error {
	vars := map[string]*TVar{}

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
				it.FindFirst("lval").ForEach(func(it *iterator.Type) {
					it.Enter().First()
					// fmt.Printf("---%v\n", o.parser.ByID(it.Node().Type))

					v := &TVar{
						name:  it.FindFirst("ident").Value(),
						label: it.FindFirst("string").Value(),
					}
					fmt.Printf("%q %q\n", v.name, v.label)
					if _, exist := vars[v.name]; exist {
						errors.Addf("duplicate identifier %v", v.name)
						return
					}
					vars[v.name] = v
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
		o.vars = vars
		return errors
	}

	traverseExpr = func(root *ptool.TNode) *errors.TErrors {
		return nil
	}

	errors := traverse(o.tree)
	return errors.Get()
}
