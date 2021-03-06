package query

import (
	"log"
	"strings"
)

// A ReferenceNode is a forward-pointing foreignKey relationship, and can define a one-to-one or one-to-many relationship,
// with this side always being the one side. It also sets up a many-to-one relationship, if it is not unique, when viewed
// from the other side of the relationship.
type ReferenceNode struct {
	Node

	// Which database in the global list of databases does the node belong to
	dbKey string
	// Name of table in the database we point to
	dbTable string
	// The name of the table that is the foreign key
	dbColumn string
	// The name of the table related to this reference
	goColumnName string
	// The name of the function used to access the property as a node or ORM item
	goPropName string
	// The name of the variable in the model structure used to hold the object
	goVarName string

	// Is this pointing to a type table item?
	isTypeTable bool

	// The name of the table we are joining to
	refTable string
	// If a forward reference and NoSQL, the name of the table that will contain the reference or references backwards to us. If SQL, the Pk of the RefTable
	refColumn string
}

// NewReferenceNode creates a forward reference node.
func NewReferenceNode(
	dbKey string,
	dbTableName string,
	dbColumnName string,
	goColumnName string,
	goName string,
	refTableName string,
	refColumn string, // only used in NoSQL situation
	isType bool,
) *ReferenceNode {
	n := &ReferenceNode{
		dbKey:        dbKey,
		dbTable:      dbTableName,
		dbColumn:     dbColumnName,
		goColumnName: goColumnName,
		goPropName:   goName,
		refTable:     refTableName,
		refColumn:    refColumn,
		isTypeTable:  isType,
	}
	return n
}

func (n *ReferenceNode) nodeType() NodeType {
	return REFERENCE_NODE
}

func (n *ReferenceNode) Equals(n2 NodeI) bool {
	if n2.nodeType() == REFERENCE_NODE {
		cn := n2.(TableNodeI).EmbeddedNode_().(*ReferenceNode)
		return cn.dbTable == n.dbTable &&
			cn.goPropName == n.goPropName &&
			(cn.alias == "" || n.alias == "" || cn.alias == n.alias)

	}
	return false
}

func (n *ReferenceNode) tableName() string {
	return n.refTable
}

func (n *ReferenceNode) setCondition(condition NodeI) {
	n.condition = condition
}

func (n *ReferenceNode) getCondition() NodeI {
	return n.condition
}

func (n *ReferenceNode) log(level int) {
	tabs := strings.Repeat("\t", level)
	log.Print(tabs + "R: " + n.dbTable + "." + n.dbColumn + "." + n.refTable + " AS " + n.GetAlias())
}

// Return the name as a capitalized object name
func (n *ReferenceNode) goName() string {
	return n.goPropName
}

// Return a node for the table that is the foreign key
func (n *ReferenceNode) relatedColumnNode() *ColumnNode {
	n2 := NewColumnNode(n.dbKey, n.dbTable, n.dbColumn, n.goColumnName, ColTypeString)
	SetParentNode(n2, n.parentNode)
	return n2
}

func RelatedColumnNode(n NodeI) NodeI {
	if tn, _ := n.(TableNodeI); tn != nil {
		if rn, _ := tn.EmbeddedNode_().(*ReferenceNode); rn != nil {
			return rn.relatedColumnNode()
		}
	}
	return nil
}

func ReferenceNodeRefTable(n *ReferenceNode) string {
	return n.refTable
}

func ReferenceNodeRefColumn(n *ReferenceNode) string {
	return n.refColumn
}

func ReferenceNodeDbColumnName(n *ReferenceNode) string {
	return n.dbColumn
}
