package control

import (
	"github.com/spekary/goradd/page/control_base"
	"github.com/spekary/goradd/page"
)


// Panel is a Goradd control that is a basic "div" wrapper. Use it to style and listen to events on a div. It
// can also be used as the basis for more advanced javascript controls.
type Panel struct {
	control_base.Panel
}

func NewPanel(parent page.ControlI) *Panel {
	p := &Panel{}
	p.Tag = "div"
	p.Init(p, parent)
	return p
}
