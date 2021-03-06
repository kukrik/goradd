//** This file was code generated by got. ***

package page

import (
	"bytes"
	"context"

	"github.com/spekary/goradd/html"
)

func ErrorTmpl(ctx context.Context, ctrl ControlI, h string, buf *bytes.Buffer) {

	buf.WriteString(`<div id="`)

	buf.WriteString(ctrl.ID())

	buf.WriteString(`_ctl" `)

	buf.WriteString(ctrl.WrapperAttributes().String())

	buf.WriteString(` >
`)

	buf.WriteString(html.Indent(h))

	buf.WriteString(`
  <div id="`)

	buf.WriteString(ctrl.ID())

	buf.WriteString(`_err" class="goradd-error">`)

	buf.WriteString(ctrl.ValidationMessage())

	buf.WriteString(`</div>
`)
	if ctrl.Instructions() != "" {
		buf.WriteString(`  <div class="goradd-instructions">`)

		buf.WriteString(ctrl.Instructions())

		buf.WriteString(`</div>
`)
	}

	buf.WriteString(`</div>
`)

	return

}
