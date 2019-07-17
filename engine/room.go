package engine

//Room - basic room structure
type Room struct {
	Desc  string
	North string
	South string
	East  string
	West  string

	IsVisited bool

	Items []Itemer
}

//Spacer - room interface
type Spacer interface {
	BasicRoom() *Room
	LeaveRoom(dir string) (bool, string)
	EnterRoom(name string) string
	OnAction(action *Action) string
	Look() string
}

//LeaveRoom - check if player can leave room, or display any extra message
func (room *Room) LeaveRoom(dir string) (bool, string) {
	return true, ""
}

//EnterRoom - check if player enters room
func (room *Room) EnterRoom(name string) string {
	if room.IsVisited {
		return name
	}
	room.IsVisited = true
	return room.Look()
}

//OnAction - callback
func (room *Room) OnAction(action *Action) string {

	if action.IsPredefined || action.IsItemRequired {
		return ""
	}
	return "You can't " + action.Name + " here."
}

//Look -
func (room *Room) Look() string {
	return room.Desc + notifyAboutVisibleItems(room.Items, " here")
}

//BasicRoom - provides general room data
func (room *Room) BasicRoom() *Room {
	return room
}
