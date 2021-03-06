package control

import (
	"github.com/spekary/goradd/page/control"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/event"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/html"
	config2 "github.com/spekary/goradd/bootstrap/config"
	"context"
	"fmt"
)

type ModalBackdropType int

const (
	ModalBackdrop ModalBackdropType = iota
	ModalNoBackdrop
	ModalStaticBackdrop
)

type ModalI interface {
	control.DialogI
}


// Modal is a bootstrap modal dialog.
// To use a custom template in a bootstrap modal, add a Panel child element or subclass of a panel
// child element. To use the grid system, add the container-fluid class to that embedded panel.
type Modal struct {
	control.Panel
	isOpen      bool

	closeOnEscape bool
	sizeClass	string

	titleBar *TitleBar
	buttonBar *control.Panel
	backdrop ModalBackdropType

	foundRight bool	// utility for adding buttons. No need to serialize this.
}

// event codes
const (
	ButtonClick = iota + 3000
	DialogClosing
	DialogClosed
)

func NewModal(parent page.ControlI, id string) *Modal {
	d := &Modal{}
	d.Init(d, parent, id)
	return d
}

func (d *Modal) Init(self page.ControlI, parent page.ControlI, id string) {
	d.Panel.Init(self, parent, id)
	d.Tag = "div"
	d.SetShouldAutoRender(true)

	d.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.SetBlockParentValidation(true)
	d.On(event.Event("hide.bs.modal").Validate(page.ValidateNone), action.Trigger(d.ID(), event.DialogClosingEvent, nil))
	d.On(event.Event("hidden.bs.modal").Validate(page.ValidateNone), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
	d.On(event.Event("hidden.bs.modal").Validate(page.ValidateNone), action.Ajax(d.ID(), DialogClosed), action.PrivateAction{})
	config2.LoadBootstrap(d.ParentForm())

	d.AddClass("modal fade").
		SetAttribute("tabindex", -1).
		SetAttribute("role", "dialog").
		SetAttribute("aria-labelledby", d.ID() + "-title").
		SetAttribute("aria-hidden", true)
	d.titleBar = NewTitleBar(d, d.ID() + "-titlebar")
	d.buttonBar = control.NewPanel(d, d.ID() + "-btnbar")
}

func (d *Modal) this() ModalI {
	return d.Self.(ModalI)
}

func (d *Modal) SetTitle(t string) control.DialogI {
	if d.titleBar.title != t {
		d.titleBar.title = t
		d.titleBar.Refresh()
	}
	return d.this()
}

func (d *Modal) SetHasCloseBox(h bool) control.DialogI {
	if d.titleBar.hasCloseBox != h {
		d.titleBar.hasCloseBox = h
		d.titleBar.Refresh()
	}
	return d.this()
}

func (d *Modal) SetState(state int) control.DialogI {
	var class string
	switch state {
	case control.DialogStateDefault:
	class = BackgroundColorNone + " " + TextColorBody
	case control.DialogStateWarning:
		class = BackgroundColorWarning + " " + TextColorBody
	case control.DialogStateError:
		class = BackgroundColorDanger + " " + TextColorLight
	case control.DialogStateSuccess:
		class = BackgroundColorSuccess + " " + TextColorLight
	case control.DialogStateInfo:
		class = BackgroundColorInfo + " " + TextColorLight
	}
	d.titleBar.RemoveClassesWithPrefix("bg-")
	d.titleBar.RemoveClassesWithPrefix("text-")
	d.titleBar.AddClass(class)
	return d.this()
}


func (d *Modal) SetBackdrop(b ModalBackdropType) ModalI {
	d.backdrop = b
	d.Refresh()
	return d.this()
}

func (d *Modal) Title() string {
	return d.titleBar.title
}

func (d *Modal) AddTitlebarClass(class string) ModalI {
	d.titleBar.AddClass(class)
	return d.this()
}

func (d *Modal) DrawingAttributes() *html.Attributes {
	a := d.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "bs-modal")
	return a
}

// AddButton adds a button to the modal. Buttons should be added in the order to appear.
// Styling options you can include in options.Options:
//  style - ButtonStyle value
//  size - ButtonSize value
func (d *Modal) AddButton(
	label string,
	id string,
	options *control.DialogButtonOptions,
)  ModalI {
	if id == "" {
		id = label
	}
	btn := NewButton(d.buttonBar, d.ID() + "-btn-" + id)
	btn.SetLabel(label)

	if options != nil {
		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}

		if options.Validates {
			btn.SetValidationType(page.ValidateContainer)
		}

		if !options.PushLeft && !d.foundRight {
			btn.AddClass("ml-auto")
			d.foundRight = true
		}

		if options.IsClose {
			btn.SetAttribute("data-dismiss", "modal") // make it a close button
		} else if options.ConfirmationMessage == "" {
			btn.On(event.Click(), action.Trigger(d.ID(), event.DialogButtonEvent, id))
		} else {
			btn.On(event.Click(),
				action.Confirm(options.ConfirmationMessage),
				action.Trigger(d.ID(), event.DialogButtonEvent, id),
			)
		}

		if options.Options != nil && len(options.Options) > 0 {
			if _,ok := options.Options["style"]; ok {
				btn.SetButtonStyle(options.Options["style"].(ButtonStyle))
			}
			if _,ok := options.Options["size"]; ok {
				btn.SetButtonSize(options.Options["size"].(ButtonSize))
			}
		}
	}

	d.buttonBar.Refresh()
	return d.this()
}

func (d *Modal) RemoveButton(id string) {
	d.buttonBar.RemoveChild(d.ID() + "-btn-" + id)
	d.buttonBar.Refresh()
}

func (d *Modal) RemoveAllButtons() {
	d.buttonBar.RemoveChildren()
	d.Refresh()
}

func (d *Modal) SetButtonVisible(id string, visible bool) ModalI {
	if ctrl := d.buttonBar.Child(d.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetVisible(visible)
	}

	return d.this()
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Modal) SetButtonStyle(id string, a *html.Style) ModalI {
	if ctrl := d.buttonBar.Child(d.ID() + "-btn-" + id); ctrl != nil {
		ctrl.SetStyles(a)
	}
	return d.this()
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (d *Modal) AddCloseButton(label string) ModalI {
	d.AddButton(label,"", &control.DialogButtonOptions{IsClose:true})
	return d.this()
}

func (d *Modal) PrivateAction(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClosed:
		d.closed()
	}
}

func (d *Modal) Open() {
	if d.Parent() == nil {
		d.SetParent(d.ParentForm()) // This is a saved modal which has previously been created and removed. Insert it back into the form.
	}
	d.SetVisible(true)
	d.isOpen = true
	//d.Refresh()
	d.AddRenderScript("modal", "show")
}

func (d *Modal) Close() {
	d.ParentForm().Response().ExecuteControlCommand(d.ID(), "modal", page.PriorityLow, "hide")
}


func (d *Modal) closed() {
	d.isOpen = false
	//d.Remove()
	d.SetVisible(false)
}

func (d *Modal) PutCustomScript(ctx context.Context, response *page.Response) {
	var backdrop interface{}

	switch d.backdrop {
	case ModalBackdrop:
		backdrop = true
	case ModalNoBackdrop:
	backdrop = false
	case ModalStaticBackdrop:
		backdrop = "static"
	}

	script := fmt.Sprintf (
`$j("#%s").modal({backdrop: %#v, keyboard: %t, focus: true, show: %t});`,
			d.ID(), backdrop, d.closeOnEscape, d.isOpen)
	response.ExecuteJavaScript(script, page.PriorityStandard)
}


/**
Alert creates a message dialog.

If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
specify just a string in buttons, or just one string as a slice of strings, one button will be shown that will just close the message.

If you specify more than one button, the first button will be the default button (the one pressed if the user presses the return key). In
this case, you will need to detect the button by adding a On(event.DialogButton(), action) to the dialog returned.
You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.
*/
func BootstrapAlert(form page.FormI, message string, buttons interface{}) control.DialogI {
	dlg := NewModal(form, "")
	dlg.SetText(message)
	if buttons != nil {
		switch b := buttons.(type) {
		case string:
			dlg.AddCloseButton(b)
		case []string:
			if len(b) == 1 {
				dlg.AddCloseButton(b[0])
			} else {
				dlg.AddButton(b[0], "", &control.DialogButtonOptions{IsPrimary: true})
				for _, l := range b[1:] {
					dlg.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dlg.SetHasCloseBox(true)
	}
	dlg.Open()
	return dlg
}




type TitleBar struct {
	control.Panel
	hasCloseBox bool
	title string
}

func NewTitleBar(parent page.ControlI, id string) *TitleBar {
	d := &TitleBar{}
	d.Panel.Init(d, parent, id)
	return d
}

func init() {
	control.RegisterAlertFunc(BootstrapAlert)
}