package control

import (
	"context"
	"github.com/spekary/goradd/html"
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/page/action"
	"github.com/spekary/goradd/page/control/control_base"
	"github.com/spekary/goradd/page/event"
)

// event codes
const (
	ButtonClick = iota + 3000
	DialogClose
)

const DialogButtonEvent = "gr-dlgbtn"

const (
	DialogStateDefault = iota
	DialogStateError
	DialogStateWarning
	DialogStateInfo
	DialogStateSuccess
)

/*
DialogI defines the publicly coconsumable api that the QCubed framework uses to interact with a dialog.

More and more CSS and javascript frameworks are coming out with their own forms of dialog, which is usually a
combination of html tag(s), css and javascript widget. QCubed has many ways of potentially interacting with
* dialogs, but to be able to inject a dialog into the framework, we need a consistent interface for all to use.
*
* This particular interface has been implemented in both JQuery UI dialogs and Bootstrap dialogs. As more needs arise,
* we can modify the interface to accomodate as many frammeworks as possible.
*
* Dialogs should descend from the Panel control. Dialogs should be able to be a member of a form or control object
* and appear with an Open call, but they should also be able to be instantiated on the fly. The framework has hooks for
* both, and if you are creating a dialog implementation, see the current JQuery UI and Bootstrap implementations for more
* direction.
*
* Feel free to implement more than just the function listed. These are the minimal set to allow your dialog to be used
* by the default QCubed framework.
*/
type DialogI interface {
	control_base.PanelI

	SetTitle(string) DialogI
	SetState(state int) DialogI
}

// Our own implementation of a dialog. This works cooperatively with javascript in qcubed.js to create a minimal
// implementation of the dialog interface.
type Dialog struct {
	control_base.Panel
	buttonBar   *Panel
	titleBar    *Panel
	closeBox    *Button
	isOpen      bool
	dialogState int
	title       string
	//validators map[string]bool
}

// DialogButtonOptions are optional additional items you can add to a dialog button.
type DialogButtonOptions struct {
	// Set Validates to true to indicate that this button will validate the dialog
	Validates bool
	// Set IsPrimary to true to make this a submit button so the user can press enter to activate it
	IsPrimary bool
	// ConfirmationMessage will appear with a yes/no box making sure the user wants the action. This is usually used
	// when the action could be destructive, like a Delete button.
	ConfirmationMessage string
	// PushLeft pushes this button to the left side of the dialog. Buttons are typically aligned right. This is helpful to separate particular
	// buttons from the main grouping of buttons.
	PushLeft bool
	// IsClose will set the button up to automatically close the dialog. Detect closes with the DialogCloseEvent if needed.
	// The button will not send a DialogButton event.
	IsClose bool
	// Options are additional options specific to the dialog implementation you are using.
	Options map[string]interface{}
}

func NewDialog(parent page.ControlI, id string) *Dialog {
	d := &Dialog{}

	d.Init(d, parent, id) // parent is always the overlay
	return d
}

func (d *Dialog) Init(self DialogI, parent page.ControlI, id string) {
	// We add the dialog to the overlay. The overlay acts as a dialog controller/container too.
	overlay := parent.Page().GetControl("groverlay")

	if overlay == nil {
		overlay = NewPanel(parent.ParentForm(), "groverlay")
		overlay.SetShouldAutoRender(true)
	} else {
		overlay.SetVisible(true)
	}

	d.Panel.Init(self, overlay, id)
	d.Tag = "div"

	d.titleBar = NewPanel(d, d.ID() + "_title")
	d.titleBar.AddClass("gr-dialog-title")

	d.buttonBar = NewPanel(d, d.ID() + "_buttons")
	d.buttonBar.AddClass("gr-dialog-buttons")
	d.SetValidationType(page.ValidateChildrenOnly) // allows sub items to validate and have validation stop here
	d.On(event.DialogClosed(), action.Ajax(d.ID(), DialogClose), action.PrivateAction{})

	//d.FormBase().AddStyleSheetFile(config.GORADD_FONT_AWESOME_CSS, nil)
}

func (d *Dialog) SetTitle(t string) DialogI {
	d.titleBar.SetText(t)
	return d
}

func (d *Dialog) Title() string {
	return d.titleBar.Text()
}

func (d *Dialog) SetState(state int) DialogI {
	return d
}

func (d *Dialog) DrawingAttributes() *html.Attributes {
	a := d.Panel.DrawingAttributes()
	a.SetDataAttribute("grctl", "dialog")
	return a
}

func (d *Dialog) AddButton(
	label string,
	id string,
	options *DialogButtonOptions,
) page.ControlI {
	if label == "" {
		id = label
	}
	btn := NewButton(d.buttonBar, id)
	btn.SetLabel(label)

	if options != nil {
		if options.IsPrimary {
			btn.SetIsPrimary(true)
		}

		if options.Validates {
			//d.validators[id] = true
			btn.SetValidationType(page.ValidateContainer)
		}

		if options.PushLeft {
			btn.AddClass("push-left")
		}

		if options.ConfirmationMessage == "" {
			btn.On(event.Click(), action.Trigger(d.ID(), DialogButtonEvent, id))
		} else {
			btn.On(event.Click(),
				action.Confirm(options.ConfirmationMessage),
				action.Trigger(d.ID(), DialogButtonEvent, id),
			)
		}
	}

	d.Refresh()
	return btn
}

func (d *Dialog) RemoveButton(id string) {
	d.buttonBar.RemoveChild(id)
	d.Refresh()
	//delete(d.validators, id)

}

func (d *Dialog) RemoveAllButtons() {
	d.buttonBar.RemoveChildren()
	d.Refresh()
	//delete(d.validators, id)
}

func (d *Dialog) SetButtonVisible(id string, visible bool) {
	if ctrl := d.buttonBar.Child(id); ctrl != nil {
		ctrl.SetVisible(false)
	}
}

// SetButtonStyle sets css styles on a button that is already in the dialog
func (d *Dialog) SetButtonStyles(id string, a *html.Style) {
	if ctrl := d.buttonBar.Child(id); ctrl != nil {
		ctrl.SetStyles(a)
	}
}

func (d *Dialog) HasCloseBox() page.ControlI {
	d.addCloseBox()
	return d
}

func (d *Dialog) addCloseBox() {
	d.closeBox = NewButton(d.titleBar, d.ID() + "_closebox")
	d.closeBox.AddClass("gr-dialog-close")
	d.closeBox.SetText(`<i class="fa fa-times"></i>`)
	d.closeBox.SetEscapeText(false)
	d.closeBox.On(event.Click(), action.Ajax(d.ID(), DialogClose))
}

// AddCloseButton adds a button to the list of buttons with the given label, but this button will trigger the DialogCloseEvent
// instead of the DialogButtonEvent. The button will also close the dialog (by hiding it).
func (d *Dialog) AddCloseButton(label string, id string) {
	btn := NewButton(d.buttonBar, id)
	btn.SetLabel(label)
	btn.On(event.Click(), action.Trigger(d.ID(), event.DialogClosedEvent, nil))
	// Note: We will also do the public doAction with a DialogCloseEvent
}

func (d *Dialog) Action(ctx context.Context, a page.ActionParams) {
	switch a.ID {
	case DialogClose:
		d.Close()
	}
}

func (d *Dialog) Open() {
	d.SetVisible(true)
	d.isOpen = true
}

func (d *Dialog) Close() {
	d.SetVisible(false)
	d.isOpen = false
	parent := d.Parent()
	if len(parent.Children()) == 1 {
		parent.SetVisible(false)
	}
	d.Remove()
}

func (d *Dialog) SetDialogState(s int) *Dialog {
	d.dialogState = s
	d.Refresh()
	return d
}

/**
Alert creates a message dialog.

If you specify no buttons, a close box in the corner will be created that will just close the dialog. If you
specify just a string in buttons, or just one string as a slice of strings, one button will be shown that will just close the message.

If you specify more than one button, the first button will be the default button (the one pressed if the user presses the return key). In
this case, you will need to detect the button by adding a On(event.DialogButton(), action) to the dialog returned.
You will also be responsible for calling "Close()" on the dialog after detecting a button in this case.

Call RegisterAlertFunc to register a different alert function for the framework to use.

*/

func Alert(form page.FormI, message string, buttons interface{}) DialogI {
	return alertFunc(form, message, buttons)
}

func DefaultAlert(form page.FormI, message string, buttons interface{}) DialogI {
	dlg := NewDialog(form, "")
	dlg.SetText(message)
	if buttons != nil {
		switch b := buttons.(type) {
		case string:
			dlg.AddCloseButton(b,"")
		case []string:
			if len(b) == 1 {
				dlg.AddCloseButton(b[0],"")
			} else {
				dlg.AddButton(b[0], "", &DialogButtonOptions{IsPrimary: true})
				for _, l := range b[1:] {
					dlg.AddButton(l, "", nil)
				}
			}
		}
	} else {
		dlg.HasCloseBox()
	}
	dlg.Open()
	return dlg
}

type AlertFuncType func(form page.FormI, message string, buttons interface{}) DialogI


var alertFunc AlertFuncType = DefaultAlert // default to our built in one

func RegisterAlertFunc(f AlertFuncType) {
	alertFunc = f
}