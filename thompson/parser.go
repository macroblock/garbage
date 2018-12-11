package main

import (
	"fmt"

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

// TErrors -
type TErrors struct {
	errors []error
}

// NewErrors -
func NewErrors() *TErrors {
	return &TErrors{}
}

// Add -
func (o *TErrors) Add(err ...error) {
	o.errors = append(o.errors, err...)
}

// Addf -
func (o *TErrors) Addf(format string, vals ...interface{}) {
	err := fmt.Errorf(format, vals...)
	o.errors = append(o.errors, err)
}

// Get -
func (o *TErrors) Get() []error {
	if len(o.errors) == 0 {
		return nil
	}
	return []error(o.errors)
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

func getNode(node *ptool.TNode, typ string) *ptool.TNode {
	if node == nil {
		return nil
	}
	tid := parser.ByName(typ)
	for _, v := range node.Links {
		if v.Type == tid {
			return v
		}
	}
	return nil
}

func getValue(node *ptool.TNode) string {
	if node != nil {
		return node.Value
	}
	return ""
}

func getNodeValue(node *ptool.TNode, typ string) string {
	return getValue(getNode(node, typ))
}

// Build -
func (o *TParser) Build() []error {
	vars := map[string]*TVar{}

	var traverse func(*ptool.TNode) *TErrors
	traverse = func(root *ptool.TNode) *TErrors {
		errors := NewErrors()
		for _, node := range root.Links {
			nodeType := o.parser.ByID(node.Type)
			switch nodeType {
			default:
				errors.Addf("an unsupported element %q", nodeType)
				continue
			case "comment":
			case "nodeDecl", "blockDecl":
				// variable.typ = strings.TrimSuffix(nodeType, "Decl")
				lval := getNode(node, "lval")
				for _, nd := range lval.Links {
					v := &TVar{
						name:  getNodeValue(nd, "ident"),
						label: getNodeValue(nd, "string"),
					}
					// fmt.Printf("%q, %q\n", v.name, v.label)
					if _, exist := vars[v.name]; exist {
						errors.Addf("duplicate identifier %v", v.name)
						continue
					}
					vars[v.name] = v
				}
			} // switch nodeType
		} // for _, node := range root.Links
		o.vars = vars
		return errors
	}

	errors := traverse(o.tree)
	return errors.Get()
}
