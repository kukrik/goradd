//** This file was code generated by got. ***

package template

import (
	"bytes"
	"fmt"
	"goradd/config"

	"github.com/spekary/goradd/codegen/generator"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/util"
)

func init() {
	t := ConnectorBaseTemplate{
		generator.Template{
			Overwrite: true,
			TargetDir: config.LocalDir + "/gen",
		},
	}
	generator.AddTableTemplate(&t)
}

type ConnectorBaseTemplate struct {
	generator.Template
}

func (n *ConnectorBaseTemplate) FileName(key string, t generator.TableType) string {
	return n.TargetDir + "/" + key + "/connector/" + t.GoName + "Base.go"
}

func (n *ConnectorBaseTemplate) GenerateTable(codegen generator.Codegen, dd *db.DatabaseDescription, t generator.TableType, buf *bytes.Buffer) {
	//connector.tmpl

	// The master template for the connector classes

	var privateName = util.LcFirst(t.GoName)

	buf.WriteString(`package connector

// This file is code generated. Do not edit.

`)

	// import.tmpl

	buf.WriteString(`import (
    "context"
    "github.com/spekary/goradd/page"
	"goradd/gen/`)

	buf.WriteString(fmt.Sprintf("%v", dd.DbKey))

	buf.WriteString(`/model"
`)
	for _, imp := range t.Imports {
		if imp.Alias == "" {
			buf.WriteString(`    "`)

			buf.WriteString(imp.Path)

			buf.WriteString(`"
`)
		} else {

			buf.WriteString(`    `)

			buf.WriteString(imp.Alias)

			buf.WriteString(` "`)

			buf.WriteString(imp.Path)

			buf.WriteString(`"
`)
		}
	}

	buf.WriteString(`
)
`)

	// struct.tmpl

	buf.WriteString(`// `)

	buf.WriteString(fmt.Sprintf("%v", privateName))

	buf.WriteString(`Base is a base structure to be embedded in a "subclass" and provides the code generated
// controls and CRUD operations.

type `)

	buf.WriteString(privateName)

	buf.WriteString(`Base struct {
    ParentControl page.ControlI
    EditMode bool
    `)

	buf.WriteString(t.GoName)

	buf.WriteString(` *model.`)

	buf.WriteString(t.GoName)

	buf.WriteString(`
`)
	for _, col := range t.Columns {
		if col.Generator != nil {
			buf.WriteString(`    `)

			buf.WriteString(col.ControlName)

			buf.WriteString(` *`)

			buf.WriteString(col.Import.Namespace)

			buf.WriteString(`.`)

			buf.WriteString(col.ControlType)

			buf.WriteString(`
`)
		}
	}

	buf.WriteString(`}
`)

	buf.WriteString(`func New`)

	buf.WriteString(fmt.Sprintf("%v", t.GoName))

	buf.WriteString(`Connector(parent page.ControlI) *`)

	buf.WriteString(fmt.Sprintf("%v", t.GoName))

	buf.WriteString(` {
    c := new(`)

	buf.WriteString(t.GoName)

	buf.WriteString(`)
    c.ParentControl = parent
    return c
}

`)

	// createControl.tmpl

	buf.WriteString(`
`)
	for _, col := range t.Columns {
		if col.Generator != nil {
			buf.WriteString(`func (c *`)

			buf.WriteString(fmt.Sprintf("%v", t.GoName))

			buf.WriteString(`) New`)

			buf.WriteString(col.ControlName)

			buf.WriteString(`(id string) *`)

			buf.WriteString(col.Import.Namespace)

			buf.WriteString(`.`)

			buf.WriteString(col.ControlType)

			buf.WriteString(` {
    var ctrl *`)

			buf.WriteString(col.Import.Namespace)

			buf.WriteString(`.`)

			buf.WriteString(col.ControlType)

			buf.WriteString(`
`)

			buf.WriteString(col.Generator.GenerateCreate(col.Import.Namespace, col.ColumnDescription))

			buf.WriteString(`
    c.`)

			buf.WriteString(col.ControlName)

			buf.WriteString(` = ctrl
    return ctrl
}
`)
		}
	}

	buf.WriteString(`
`)

	// load.tmpl

	buf.WriteString(`
// Load will associate the controls with data from the given model.`)

	buf.WriteString(t.GoName)

	buf.WriteString(` object and load the controls with data.
// Generally call this after creating the controls. Otherwise, call Refresh if you Load before creating the controls.
// If you pass a nil object, it will prepare the controls for creating a new record in the database.
func (c *`)

	buf.WriteString(privateName)

	buf.WriteString(`Base) Load(ctx context.Context, modelObj *model.`)

	buf.WriteString(t.GoName)

	buf.WriteString(`) {
    c.`)

	buf.WriteString(t.GoName)

	buf.WriteString(` = modelObj
    if modelObj.PrimaryKey() == "" {
        c.EditMode = false
    } else {
        c.EditMode = true
    }
    c.Refresh()
}


// Refresh loads the controls with the cached model.`)

	buf.WriteString(t.GoName)

	buf.WriteString(` object.
func (c *`)

	buf.WriteString(privateName)

	buf.WriteString(`Base) Refresh() {
`)
	for _, col := range t.Columns {
		buf.WriteString(`
`)
		var sLoad string
		if col.Generator != nil {
			sLoad = col.Generator.GenerateGet(col.ControlName, t.GoName, col.ColumnDescription)
		}
		if sLoad != "" {
			buf.WriteString(`    if c.`)

			buf.WriteString(fmt.Sprintf("%v", col.ControlName))

			buf.WriteString(` != nil {
        `)

			buf.WriteString(sLoad)

			buf.WriteString(`
    }
`)
		}
	}

	buf.WriteString(`}
`)

	// save.tmpl

	buf.WriteString(`
// Update will update the related model.`)

	buf.WriteString(t.GoName)

	buf.WriteString(` object with data from the controls.
func (c *`)

	buf.WriteString(privateName)

	buf.WriteString(`Base) Update() {
`)
	for _, col := range t.Columns {
		buf.WriteString(`
`)
		var sUpdate string
		if col.Generator != nil {
			sUpdate = col.Generator.GeneratePut(col.ControlName, t.GoName, col.ColumnDescription)
		}
		if sUpdate != "" {
			buf.WriteString(`    if c.`)

			buf.WriteString(fmt.Sprintf("%v", col.ControlName))

			buf.WriteString(` != nil {
        `)

			buf.WriteString(sUpdate)

			buf.WriteString(`
    }
`)
		}
	}

	buf.WriteString(`}

// Save takes the data from the controls and saves it in the database.
func (c *`)

	buf.WriteString(privateName)

	buf.WriteString(`Base) Save(ctx context.Context) {
    c.Update()
    c.`)

	buf.WriteString(t.GoName)

	buf.WriteString(`.Save(ctx)
}

`)

	// delete.tmpl

	buf.WriteString(`
// Delete will deleted the related object.
func (c *`)

	buf.WriteString(privateName)

	buf.WriteString(`Base) Delete(ctx context.Context) {
    c.`)

	buf.WriteString(t.GoName)

	buf.WriteString(`.Delete(ctx)
}
`)

}

func (n *ConnectorBaseTemplate) Overwrite() bool {
	return n.Template.Overwrite
}