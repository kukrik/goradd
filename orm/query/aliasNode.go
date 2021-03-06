package query

import (
	"log"
	"strings"
)

// An AliasNode allows reference to a prior aliased operation later in a query.

type AliasNodeI interface {
	NodeI
}

type AliasNode struct {
	Node
}

// Alias returns an AliasNode type, which allows you to refer to a prior created named alias operation.

func Alias(goName string) *AliasNode {
	return &AliasNode{
		Node: Node{
			alias: goName,
		},
	}
}

func (n *AliasNode) nodeType() NodeType {
	return ALIAS_NODE
}

func (n *AliasNode) tableName() string {
	return ""
}

func (n *AliasNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == ALIAS_NODE {
		return n.GetAlias() == n2.GetAlias()
	}

	return false
}

func (n *AliasNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "Alias: " + n.GetAlias())
}

// Return the name as a capitalized object name
func (n *AliasNode) goName() string {
	return n.GetAlias()
}
