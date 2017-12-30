//** This file was code generated by got. ***

package template

import (
	"bytes"
	"grlocal/config"

	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/util"
)

func init() {
	t := ModelTemplate{
		generator.Template{
			Overwrite: true,
			TargetDir: config.LocalDir + "/model",
		},
	}
	generator.AddTableTemplate(&t)
}

type ModelTemplate struct {
	generator.Template
}

func (n *ModelTemplate) FileName(t *db.TableDescription) string {
	return n.TargetDir + "/" + t.GoName + ".go"
}

func (n *ModelTemplate) GenerateTable(codegen generator.Codegen, dd *db.DatabaseDescription, t *db.TableDescription, buf *bytes.Buffer) {

	// The master template for the model classes

	buf.WriteString(`package model

type `)
	buf.WriteString(t.GoName)
	buf.WriteString(` struct {
	`)
	buf.WriteString(util.LcFirst(t.GoName))
	buf.WriteString(`Base
}

`)

}