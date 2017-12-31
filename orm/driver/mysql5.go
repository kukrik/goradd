package driver

import (
	sqldb "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
	//"goradd/orm/query"
	"github.com/knq/snaker"
	//"github.com/spekary/goradd/util"
	"context"
	"strconv"
	"github.com/spekary/goradd/orm/db"
)

// goradd additions to the mysql database driver
// call NewMysql5, but afterwards, work through the DB parent interface

type Mysql5 struct {
	SqlDb
	description *db.DatabaseDescription
	config      *mysql.Config
}

// New Mysql5 creates a new Mysql5 object and returns its matching interface
func NewMysql5(dbKey string, params string, config *mysql.Config) *Mysql5I {
	var err error

	m := Mysql5{
		SqlDb: NewSqlDb(dbKey),
	}

	if params == "" && config == nil {
		panic("Must specify how to connect to the database.")
	}
	if params == "" {
		params = config.FormatDSN()
		m.config = config
	} else {
		m.config, err = mysql.ParseDSN(params)
		if err != nil {
			panic("Could not parse the connection string.")
		}
	}
	m.db, err = sqldb.Open("mysql", params)
	if err != nil {
		panic("Could not open database: " + err.Error())
	}
	err = m.db.Ping()
	if err != nil {
		panic("Could not ping database: " + err.Error())
	}
	m.loadDescription()
}


func (m *Mysql5) Describe() *db.DatabaseDescription {
	return m.description
}

func (m *Mysql5) generateSelectSql(b *sqlBuilder) (sql string, args []interface{}) {
	var s string
	var a []interface{}

	if b.distinct {
		sql = "SELECT DISTINCT\n"
	} else {
		sql = "SELECT\n"
	}

	s, a = m.I().(Mysql5I).generateColumnListWithAliases(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateFromSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateWhereSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateGroupBySql(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateHaving(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateOrderBySql(b)
	sql += s
	args = append(args, a...)

	sql += m.I().(Mysql5I).generateLimitSql(b)

	return
}

func (m *Mysql5) generateDeleteSql(b *sqlBuilder) (sql string, args []interface{}) {
	var s string
	var a []interface{}

	n := b.rootNode

	sql = "DELETE " + n.GetAlias() + " "

	s, a = m.I().(Mysql5I).generateFromSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateWhereSql(b)
	sql += s
	args = append(args, a...)

	s, a = m.I().(Mysql5I).generateOrderBySql(b)
	sql += s
	args = append(args, a...)

	sql += m.I().(Mysql5I).generateLimitSql(b)

	return
}

func (m *Mysql5) generateColumnListWithAliases(b *sqlBuilder) (sql string, args []interface{}) {
	b.columnAliases.Range(func(key string, v interface{}) bool {
		node := v.(*ColumnNode)
		sql += m.I().(Mysql5I).generateColumnNodeSql(node, false) + " AS `" + node.GetAlias() + "`,\n"
		return true
	})

	b.aliasNodes.Range(func(key string, v interface{}) bool {
		node := v.(NodeI)
		s, a := m.I().(Mysql5I).generateNodeSql(node, false)
		sql += s + " AS `" + node.GetAlias() + "`,\n"
		args = append(args, a...)
		return true
	})

	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n"
	return
}

func (m *Mysql5) generateFromSql(b *sqlBuilder) (sql string, args []interface{}) {
	var s string
	var a []interface{}

	sql = "FROM\n"

	n := b.rootNode
	sql += "`" + n.tableName() + "` AS `" + n.GetAlias() + "`\n"

	var childNodes []NodeI
	var cn NodeI
	if childNodes = n.getChildNodes(); childNodes != nil {
		for _, cn = range childNodes {
			s, a = m.I().(Mysql5I).generateJoinSql(b, cn)
			sql += s
			args = append(args, a...)
		}
	}
	return
}

func (m *Mysql5) generateJoinSql(b *sqlBuilder, n NodeI) (sql string, args []interface{}) {
	var tn TableNodeI
	var ok bool

	if tn, ok = n.(TableNodeI); !ok {
		return
	}

	switch node := tn.EmbeddedNode_().(type) {
	case *ReferenceNode:
		sql = "LEFT JOIN "
		sql += "`" + node.refTable + "` AS `" + node.GetAlias() + "` ON `" + node.getParentNode().GetAlias() + "`.`" + node.dbColumn + "` = `" + node.GetAlias() + "`.`" + node.refColumn + "`"
		if condition := node.getCondition(); condition != nil {
			s, a := m.I().(Mysql5I).generateNodeSql(condition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ReverseReferenceNode:
		if b.limitInfo != nil {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "
		sql += "`" + node.refTable + "` AS `" + node.GetAlias() + "` ON `" + node.getParentNode().GetAlias() + "`.`" + node.dbColumn + "` = `" + node.GetAlias() + "`.`" + node.refColumn + "`"
		if condition := node.getCondition(); condition != nil {
			s, a := m.I().(Mysql5I).generateNodeSql(condition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	case *ManyManyNode:
		if b.limitInfo != nil {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "

		var pk string
		if node.isTypeTable {
			pk = snaker.CamelToSnake(m.I().(Mysql5I).Describe().TypeTableDescription(node.refTable).PkField)
		} else {
			pk = m.I().(Mysql5I).Describe().TableDescription(node.refTable).PrimaryKeyColumn.DbName
		}

		sql += "`" + node.dbTable + "` AS `" + node.GetAlias() + "a` ON `" + node.getParentNode().GetAlias() + "`.`" + node.getParentNode().(TableNodeI).PrimaryKeyNode_().name() +
			"` = `" + node.GetAlias() + "a`.`" + node.dbColumn + "`\n"
		sql += "LEFT JOIN `" + node.refTable + "` AS `" + node.GetAlias() + "` ON `" + node.GetAlias() + "a`.`" + node.refColumn +
			"` = `" + node.GetAlias() + "`.`" + pk + "`"

		if condition := node.getCondition(); condition != nil {
			s, a := m.I().(Mysql5I).generateNodeSql(condition, false)
			sql += " AND " + s
			args = append(args, a...)
		}
	default:
		return
	}
	sql += "\n"
	if childNodes := n.getChildNodes(); childNodes != nil {
		for _, cn := range childNodes {
			s, a := m.I().(Mysql5I).generateJoinSql(b, cn)
			sql += s
			args = append(args, a...)

		}
	}
	return
}

func (m *Mysql5) generateNodeSql(n NodeI, useAlias bool) (sql string, args []interface{}) {
	switch node := n.(type) {
	case *ValueNode:
		sql = "?"
		args = append(args, node.value)
	case *OperationNode:
		sql, args = m.I().(Mysql5I).generateOperationSql(node, useAlias)
	case *ColumnNode:
		sql = m.I().(Mysql5I).generateColumnNodeSql(node, useAlias)
	case *aliasNode:
		sql = "`" + node.GetAlias() + "`"
	case *SubqueryNode:
		sql, args = m.I().(Mysql5I).generateSubquerySql(node)
	default:
		if tn, ok := n.(TableNodeI); ok {
			sql = m.I().(Mysql5I).generateColumnNodeSql(tn.PrimaryKeyNode_(), false)
		} else {
			panic("Can't generate sql from m.I().(Mysql5I) node type.")
		}

	}
	return
}

func (m *Mysql5) generateSubquerySql(node *SubqueryNode) (sql string, args []interface{}) {
	sql, args = m.I().(Mysql5I).generateSelectSql(node.b.(*sqlBuilder))
	sql = "(" + sql + ")"
	return
}

func (m *Mysql5) generateOperationSql(n *OperationNode, useAlias bool) (sql string, args []interface{}) {
	if useAlias && n.GetAlias() != "" {
		sql = n.GetAlias()
		return
	}
	switch n.op {
	case OpFunc:
		if len(n.operands) > 0 {
			for _, o := range n.operands {
				s, a := m.I().(Mysql5I).generateNodeSql(o, useAlias)
				sql += s + ","
				args = append(args, a...)
			}
			sql = sql[:len(sql)-1]
		} else {
			if n.functionName == "COUNT" {
				sql = "*"
			}
		}

		if n.distinct {
			sql = "DISTINCT " + sql
		}
		sql = n.functionName + "(" + sql + ") "

	case OpNull:
		fallthrough
	case OpNotNull:
		s, a := m.I().(Mysql5I).generateNodeSql(n.operands[0], useAlias)
		sql = s + " IS " + n.op.String()
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpNot:
		s, a := m.I().(Mysql5I).generateNodeSql(n.operands[0], useAlias)
		sql = n.op.String() + " " + s
		args = append(args, a...)
		sql = "(" + sql + ") "

	case OpIn:
		fallthrough
	case OpNotIn:
		s, a := m.I().(Mysql5I).generateNodeSql(n.operands[0], useAlias)
		sql = s + " " + n.op.String() + " ("
		args = append(args, a...)

		for _, o := range n.operands[1].(*ValueNode).value.([]NodeI) {
			s, a = m.I().(Mysql5I).generateNodeSql(o, useAlias)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",") + ") "

	default:
		for _, o := range n.operands {
			s, a := m.I().(Mysql5I).generateNodeSql(o, useAlias)
			sql += s + " " + n.op.String() + " "
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, " "+n.op.String()+" ")
		sql = "(" + sql + ") "

	}
	return
}

func (m *Mysql5) generateColumnNodeSql(n *ColumnNode, useAlias bool) (sql string) {
	if useAlias {
		sql = "`" + n.GetAlias() + "`"
	} else {
		sql = "`" + n.getParentNode().GetAlias() + "`.`" + n.name() + "`"
	}
	return
}

func (m *Mysql5) generateNodeListSql(nodes []NodeI, useAlias bool) (sql string, args []interface{}) {
	for _, node := range nodes {
		s, a := m.I().(Mysql5I).generateNodeSql(node, useAlias)
		sql += s + ","
		args = append(args, a...)
	}
	sql = strings.TrimSuffix(sql, ",")
	return
}

func (m *Mysql5) generateOrderBySql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.orderBys != nil && len(b.orderBys) > 0 {
		sql = "ORDER BY "
		for _, n := range b.orderBys {
			s, a := m.I().(Mysql5I).generateNodeSql(n, true)
			if sorter, ok := n.(NodeSorter); ok {
				if sorter.sortDesc() {
					s += " DESC"
				}
			}
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *Mysql5) generateGroupBySql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.groupBys != nil && len(b.groupBys) > 0 {
		sql = "GROUP BY "
		for _, n := range b.groupBys {
			s, a := m.I().(Mysql5I).generateNodeSql(n, true)
			sql += s + ","
			args = append(args, a...)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

func (m *Mysql5) generateWhereSql(b *sqlBuilder) (sql string, args []interface{}) {
	if b.condition != nil {
		sql = "WHERE "
		var s string
		s, args = m.I().(Mysql5I).generateNodeSql(b.condition, false)
		sql += s + "\n"
	}
	return
}

func (m *Mysql5) generateHaving(b *sqlBuilder) (sql string, args []interface{}) {
	if b.having != nil {
		sql = "HAVING "
		var s string
		s, args = m.I().(Mysql5I).generateNodeSql(b.having, false)
		sql += s + "\n"
	}
	return
}

func (m *Mysql5) generateLimitSql(b *sqlBuilder) (sql string) {
	if b.limitInfo == nil {
		return ""
	}
	if b.limitInfo.offset > 0 {
		sql = strconv.FormatInt(b.limitInfo.offset, 10) + ","
	}

	if b.limitInfo.maxRowCount > -1 {
		sql += strconv.FormatInt(b.limitInfo.maxRowCount, 10)
	}

	if sql != "" {
		sql = "LIMIT " + sql + "\n"
	}

	return
}

func (m *Mysql5) Update(ctx context.Context, table string, fields map[string]interface{}, pkName string, pkValue string) {
	var sql = "UPDATE " + table + "\n"
	var args = []interface{}{}
	s, a := m.I().(Mysql5I).makeSetSql(fields)
	sql += s
	args = append(args, a...)

	sql += "WHERE " + pkName + " = ?"
	args = append(args, pkValue)
	_, e := m.I().(Mysql5I).Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

func (m *Mysql5) Insert(ctx context.Context, table string, fields map[string]interface{}) string {
	var sql = "INSERT " + table + "\n"
	var args = []interface{}{}
	s, a := m.I().(Mysql5I).makeSetSql(fields)
	sql += s
	args = append(args, a...)

	if r, err := m.I().(Mysql5I).Exec(ctx, sql, args...); err != nil {
		panic(err.Error())
	} else {
		if id, err := r.LastInsertId(); err != nil {
			panic(err.Error())
			return ""
		} else {
			return fmt.Sprint(id)
		}
	}
}

func (m *Mysql5) Delete(ctx context.Context, table string, pkName string, pkValue interface{}) {
	var sql = "DELETE FROM " + table + "\n"
	var args = []interface{}{}
	sql += "WHERE " + pkName + " = ?"
	args = append(args, pkValue)
	_, e := m.I().(Mysql5I).Exec(ctx, sql, args...)
	if e != nil {
		panic(e.Error())
	}
}

func (m *Mysql5) makeSetSql(fields map[string]interface{}) (sql string, args []interface{}) {
	if len(fields) == 0 {
		panic("No fields to set")
	}
	sql = "SET "
	for k, v := range fields {
		sql += fmt.Sprintf("%s=?, ", k)
		args = append(args, v)
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += "\n"
	return
}

func (m *Mysql5) IsA(className string) bool {
	if className == "Mysql5" {
		return true
	}
	return m.SqlDb.IsA(className)
}

func (m *Mysql5) Class() string {
	return "Mysql5"
}
