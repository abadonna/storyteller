package engine

import (
	"strings"
)

//Item - basic item structure
type Item struct {
	IsDecoration      bool
	IsSurface         bool
	IsContainer       bool
	IsPickable        bool
	IsVisible         bool
	IsDisabled        bool
	IsOpen            bool
	IsLocked          bool
	IsUnbreakableName bool
	IsUseTarget       bool

	Name     string
	AName    string
	Desc     string
	Vocab    string
	Location string
	KeyName  string

	DefaultActionDesc map[string]string
	CanContainOnly    []string

	Items []Itemer
}

//Itemer - item interface
type Itemer interface {
	Basic() *Item
	OnAction(action *Action, target Itemer) (string, string)
	NameWithArticle() string
}

//Basic - provides general item data
func (item *Item) Basic() *Item {
	return item
}

//NameWithArticle - provides item full name
func (item *Item) NameWithArticle() string {
	if item.AName != "" {
		return item.AName + " " + item.Name
	}
	if strings.IndexAny(item.Name, "aeuo") == 0 {
		return "an " + item.Name
	}
	return "a " + item.Name
}

//OnAction returns result as text and item's new location
func (item *Item) OnAction(action *Action, target Itemer) (string, string) {
	switch action {
	case EXAMINE:
		return item.Examine(), item.Location
	case OPEN:
		return item.Open(), item.Location
	case CLOSE:
		return item.Close(), item.Location
	case TAKE:
		return item.Take()
	case PUT:
		return item.Put(target)
	case UNLOCK:
		return item.Unlock(target), item.Location
	case USE:
		return item.Use(target), item.Location
	}

	return "You can't " + action.Name + " " + item.NameWithArticle() + ".", item.Location
}

//Examine item
func (item *Item) Examine() string {
	var msg string
	details := item.NameWithArticle()

	if item.Desc != "" {
		msg = item.Desc
	} else {
		msg = "You see " + details + "."
	}

	if !item.IsContainer || item.IsOpen {
		for _, i := range item.Items {
			item := i.Basic()
			if !item.IsDisabled {
				item.IsVisible = true
			}
		}
	}

	if item.IsContainer {
		details = " in " + details
	} else if item.IsSurface {
		details = " on " + details
	}

	return msg + notifyAboutVisibleItems(item.Items, details)
}

//Open container item
func (item *Item) Open() string {

	msg, ok := item.DefaultActionDesc["open"]

	if !item.IsContainer {
		if ok {
			return msg
		}
		return "I don't know how to open " + item.NameWithArticle()
	}

	if item.IsOpen {
		return "It's already opened."
	}

	if item.IsLocked {
		return "It's locked."
	}

	item.IsOpen = true
	for _, i := range item.Items {
		item := i.Basic()
		item.IsVisible = true
	}

	if !ok {
		msg = "Opened."
	}

	return msg + notifyAboutVisibleItems(item.Items, " in "+item.NameWithArticle())
}

//Close container item
func (item *Item) Close() string {
	if !item.IsContainer {
		return "It can't be closed."
	}

	if !item.IsOpen {
		return "It's already closed."
	}

	item.IsOpen = false
	for _, i := range item.Items {
		item := i.Basic()
		item.IsVisible = false
	}

	msg, ok := item.DefaultActionDesc["close"]
	if !ok {
		msg = "Closed."
	}

	return msg
}

//Take item into inventory
func (item *Item) Take() (string, string) {
	if item.Location == "inventory" {
		return "You already have it.", item.Location
	}

	msg, ok := item.DefaultActionDesc["take"]
	if !item.IsPickable {
		if ok {
			return msg, item.Location
		}
		return "You can't take it.", item.Location
	}

	if ok {
		return msg, "inventory"
	}

	return "Taken.", "inventory"
}

//Put item into inventory
func (item *Item) Put(target Itemer) (string, string) {
	if item.Location != "inventory" {
		return "You are not holding " + item.NameWithArticle() + ".", item.Location
	}
	if target == nil {
		return "Where do you want to put it?", item.Location
	}
	indirect := target.Basic()
	if indirect.IsContainer {
		if indirect.IsOpen {

			cancontain := len(indirect.CanContainOnly) == 0
			for _, s := range indirect.CanContainOnly {
				if item.Name == s {
					cancontain = true
					break
				}
			}

			if !cancontain {
				return "You can't put " + item.NameWithArticle() + " in " + target.NameWithArticle() + ".", item.Location
			}

			return "You put " + item.NameWithArticle() + " in " + target.NameWithArticle() + ".", indirect.Name
		}
		return strings.Title(target.NameWithArticle() + " is closed."), item.Location
	}
	if indirect.IsSurface {
		return "You put " + item.NameWithArticle() + " on " + target.NameWithArticle() + ".", indirect.Name
	}
	return "You can't put it here.", item.Location
}

//Unlock container with key
func (item *Item) Unlock(target Itemer) string {
	if !item.IsContainer {
		return "You can't unlock it."
	}
	if !item.IsLocked {
		return "It's not locked."
	}
	if target == nil {
		return "You need a key to unlock it."
	}

	key := target.Basic()
	if key.Location != "inventory" {
		return "You are not holding " + key.NameWithArticle() + "."
	}

	item.IsLocked = false
	return "Unlocked."
}

//Use item
func (item *Item) Use(target Itemer) string {

	if item.IsUseTarget && target != nil && !target.Basic().IsUseTarget {
		msg, _ := target.OnAction(USE, item)
		return msg
	}

	if target != nil {
		box := target.Basic()
		if box.IsContainer &&
			box.IsLocked &&
			box.KeyName == item.Name {

			msg, _ := target.OnAction(UNLOCK, item)
			return msg
		}
	}

	msg, ok := item.DefaultActionDesc["use"]
	if !ok {
		return "You can't use it."
	}

	return msg
}
