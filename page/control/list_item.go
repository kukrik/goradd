package control

import (
	"fmt"
	"github.com/spekary/goradd/html"
	html2 "html"
)

type ListItemI interface {
	ItemListI
	Value() interface{}
	ID() string
	SetID(string)
	Label() string
	SetLabel(string)
	SetDisabled(bool)
	Disabled() bool
	SetIsDivider(bool)
	IsDivider() bool
	IntValue() int
	StringValue() string
	HasChildItems() bool
	Attributes() *html.Attributes
	Anchor() string
	SetAnchor(string)
	AnchorAttributes() *html.Attributes
	RenderLabel() string
}

type ItemLister interface {
	Value() interface{}
	Label() string
}

// ItemIDer is an interface to a listable object that matches orm objects
type ItemIDer interface {
	ID() interface{}
	String() string
}



type Labeler interface {
	Label() string
}

// A ListItem is an object that is a member of a list. HTML has a few different kinds of lists, and this can be a member
// of a select list (<select>), or an ordered or unordered list (<ul> or <ol>). It is up to the manager of the list to
// render the item, but this serves as a place to store options about the item. Not all options are pertinent to
// all lists.
//
// A list item generally has a value, and a label. Often, lists will have ids too, that will appear in the html output,
// but the id values are managed by the list manager and generally should not be set by you. In situations where the
// user selects a list item, you would use the id to retrieve the ListItem selected.
type ListItem struct {
	value interface{}
	id    string
	ItemList
	label      string
	attributes *html.Attributes
	shouldEscapeLabel bool
	disabled	bool
	isDivider bool
	anchorAttributes *html.Attributes
}

// NewListItem creates a new item for a list. Specify an empty value for an item that represents no selection.
func NewListItem(label string, value ...interface{}) *ListItem {
	l := &ListItem{attributes: html.NewAttributes(), label: label}
	if c := len(value); c == 1 {
		l.value = value[0]
	} else if c > 1 {
		panic("Call NewListItem with zero or one value only.")
	}

	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromItemLister creates a new item from any object that has a Value and Label method.
func NewItemFromItemLister(i ItemLister) *ListItem {
	l := &ListItem{attributes: html.NewAttributes(), value: i.Value(), label: i.Label()}
	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromLabeler creates a new item from any object that has just a Label method.
func NewItemFromLabeler(i Labeler) *ListItem {
	l := &ListItem{attributes: html.NewAttributes(), label: i.Label()}
	l.ItemList = NewItemList(l)
	return l
}

// NewItemFromStringer creates a new item from any object that has just a String method.
func NewItemFromStringer(i fmt.Stringer) *ListItem {
	l := &ListItem{attributes: html.NewAttributes(), label: i.String()}
	l.ItemList = NewItemList(l)
	return l
}


// NewItemFromItemIDer creates a new item from any object that has an ID and String method.
// Note that the ID() of the ItemIDer will become the value of the select item, and the String()
// will become the label
func NewItemFromItemIDer(i ItemIDer) *ListItem {
	l := &ListItem{attributes: html.NewAttributes(), value: i.ID(), label: i.String()}
	l.ItemList = NewItemList(l)
	return l
}


func (i *ListItem) SetValue(v interface{}) *ListItem {
	i.value = v
	return i
}

func (i *ListItem) Value() interface{} {
	return i.value
}

func (i *ListItem) IntValue() int {
	return i.value.(int)
}

func (i *ListItem) StringValue() string {
	if s, ok := i.value.(fmt.Stringer); ok {
		return s.String()
	} else {
		return i.value.(string)
	}
}

func (i *ListItem) ID() string {
	return i.id
}

// SetID should not be called by your code typically. It is exported for implementations of item lists. The IDs of an
// item list are completely managed by the list, you cannot have custom ids.
func (i *ListItem) SetID(id string) {
	i.id = id
	i.attributes.SetID(id)
	i.ItemList.reindex(0)
}

func (i *ListItem) HasChildItems() bool {
	return i.ItemList.Len() > 0
}

func (i *ListItem) Label() string {
	return i.label
}

func (i *ListItem) SetLabel(l string) {
	i.label = l
}

func (i *ListItem) SetDisabled(d bool) {
	i.disabled = d
}

func (i *ListItem) Disabled() bool {
	return i.disabled
}

func (i *ListItem) SetIsDivider(d bool) {
	i.isDivider = d
}

func (i *ListItem) IsDivider() bool {
	return i.isDivider
}

func (i *ListItem) SetAnchor(a string) {
	i.AnchorAttributes().Set("href", a)
}

func (i *ListItem) Anchor() string {
	if i.anchorAttributes == nil || !i.anchorAttributes.Has("href") {
		return ""
	}
	return i.anchorAttributes.Get("href")
}

func (i *ListItem) AnchorAttributes() *html.Attributes {
	if i.anchorAttributes == nil {
		i.anchorAttributes = html.NewAttributes()
	}
	return i.anchorAttributes
}

func (i *ListItem) SetShouldEscapeLabel(e bool) *ListItem {
	i.shouldEscapeLabel = e
	return i
}

func (i *ListItem) RenderLabel() (h string) {
	if i.shouldEscapeLabel {
		h = html2.EscapeString(i.label)
	} else {
		 h = i.label
	}
	if i.Anchor() != ""  && !i.disabled {
		h = html.RenderTag("a", i.anchorAttributes, h)
	}
	return
}


// Attributes returns a pointer to the attributes of the item for customization. You can directly set the attributes
// on the returned object.
func (i *ListItem) Attributes() *html.Attributes {
	if i.attributes == nil {
		i.attributes = html.NewAttributes()
	}
	return i.attributes
}
