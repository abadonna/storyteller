package engine

import (
	"strings"
)

//Game - basic game structure
type Game struct {
	Location   string
	Inventory  []Itemer
	Rooms      map[string]Spacer
	Actions    []Action
	IsFinished bool
	Input      string
}

//Adventurer interface for the game context
type Adventurer interface {
	BasicGame() *Game
	Intro() string
	Help() string
	CurrentRoom() Spacer
	DoItemAction(words []string, action *Action) string
	DoActorAction(words []string, action *Action) string
	Navigate(location string, dir string) string
	ShowInventory() string
}

//BasicGame - returns basic game struct
func (game *Game) BasicGame() *Game {
	return game
}

//Intro for the game
func (game *Game) Intro() string {
	room := game.Rooms[game.Location]
	return "\n\nWELCOME!\n\n" + room.EnterRoom(game.Location) + "\n"
}

//CurrentRoom - returns current location
func (game *Game) CurrentRoom() Spacer {
	room := game.Rooms[game.Location]
	return room
}

//Process command on game context
func Process(game Adventurer, command string) string {
	command = strings.Replace(strings.TrimSpace(strings.ToLower(command)), "  ", " ", -1)

	if command == "" {
		return "I beg your pardon?"
	}

	if game.BasicGame().IsFinished {
		return "Game is finished, but you can restart it."
	}

	msg := executeCommand(game, command)
	return msg
}

func executeCommand(game Adventurer, command string) string {
	game.BasicGame().Input = command

	room := game.CurrentRoom()
	base := room.BasicRoom()
	words := strings.Split(command, " ")

	for i, word := range words {
		switch word {
		case "go", "the", "a", "an", "from":
			//skip
		case "north", "n":
			return game.Navigate(base.North, "n")
		case "south", "s":
			return game.Navigate(base.South, "s")
		case "east", "e":
			return game.Navigate(base.East, "e")
		case "west", "w":
			return game.Navigate(base.West, "w")
		case "look", "l":
			return room.Look() + room.OnAction(LOOK)
		case "examine", "x", "search":
			return game.DoItemAction(words[i+1:], EXAMINE)
		case "open":
			return game.DoItemAction(words[i+1:], OPEN)
		case "close":
			return game.DoItemAction(words[i+1:], CLOSE)
		case "unlock":
			return game.DoItemAction(words[i+1:], UNLOCK)
		case "take", "pick":
			return game.DoItemAction(words[i+1:], TAKE)
		case "put", "drop":
			return game.DoItemAction(words[i+1:], PUT)
		case "use":
			return game.DoItemAction(words[i+1:], USE)
		case "ask", "tell", "talk":
			return game.DoActorAction(words[i+1:], ASK) + room.OnAction(ASK)

		case "give", "show":
			return game.DoItemAction(words[i+1:], GIVE)

		case "help":
			return game.Help()

		case "inventory", "i":
			return game.ShowInventory() + room.OnAction(INVENTORY)
		default:
			//check custom actions
			for _, action := range game.BasicGame().Actions {
				if word == action.Name {
					if action.IsItemRequired {
						return game.DoItemAction(words[i+1:], &action)
					}
					if action.IsActorRequired {
						return game.DoActorAction(words[i+1:], &action)
					}
					return room.OnAction(&action)
				}
			}
			return "I don't know the word \"" + word + "\"."
		}
	}

	return ""
}

//DoActorAction - generic actor action processor
func (game *Game) DoActorAction(words []string, action *Action) string {
	if len(words) == 0 {
		return strings.Title(action.Name) + " who?"
	}

	room := game.Rooms[game.Location]
	items := visibleItems(room.BasicRoom().Items, false, true)
	actors := filterActors(items)

	if len(actors) == 0 {
		return "There is nobody to " + action.Name + "."
	}

	actor := actors[0]
	if len(actors) > 1 {
		var msg string
		msg, actor, words = findActor(words, actors)
		if msg != "" {
			return msg
		}
	}

	if action.IsTopicRequired {
		for _, topic := range findTopics(words, actor) {
			if strings.Contains(topic.Action, action.Name) {
				return actor.OnTopic(topic, action, nil)
			}
		}
		return actor.OnTopic(nil, action, nil)
	}

	if action.Syntax != "" {
		words, _ = findSyntax(words, action.Syntax)
	}

	if action.IsTargetRequired && len(words) == 0 && action.Syntax != "" {
		return strings.Title(action.Name) + " " + action.Syntax + " what?"
	}

	msg, target, _ := findTarget(words, items, !action.IsTargetRequired)

	if !action.IsTargetRequired {
		msg, _ := actor.OnAction(action, target)
		return msg
	}

	if msg != "" {
		return msg
	}

	if target == nil {
		return strings.Title(action.Name) + " " + action.Syntax + " what?"
	}

	msg, _ = actor.OnAction(action, target)
	return msg
}

//DoItemAction - generic item action processor
func (game *Game) DoItemAction(words []string, action *Action) string {
	if len(words) == 0 {
		return strings.Title(action.Name) + " what?"
	}

	room := game.Rooms[game.Location]
	items := visibleItems(append(game.Inventory, room.BasicRoom().Items...), false, true)

	msg, item, words := findTarget(words, items, false)

	if msg != "" {
		return msg
	}

	isSyntaxFound := false
	if action.Syntax != "" {
		words, isSyntaxFound = findSyntax(words, action.Syntax)
	}

	if action.IsTargetRequired && len(words) == 0 && action.Syntax != "" {
		return strings.Title(action.Name) + " " + action.Syntax + " what?"
	}

	if action.IsActorTarget {
		actors := filterActors(items)
		var actor Actor
		var msg string

		if len(actors) == 1 && !isSyntaxFound && action.Syntax != "" {
			actor = actors[0]
		} else {
			msg, actor, _ = findActor(words, actors)
			if msg != "" {
				return msg
			}
		}

		if action.IsTopicRequired {
			for _, topic := range findTopics(strings.Split(item.Basic().Name, " "), actor) {
				if strings.Contains(topic.Action, action.Name) {
					if topic.IsItemConsumed {
						game.ChangeParent(item, "")
					}
					msg = actor.OnTopic(topic, action, item)
					return msg
				}
			}
			return actor.OnTopic(nil, action, item)
		}

		return game.finalizeItemAction(item, actor, action)
	}

	msg, target, _ := findTarget(words, items, !action.IsTargetRequired)

	if !action.IsTargetRequired {
		return game.finalizeItemAction(item, target, action)
	}

	if msg != "" {
		return msg
	}

	if target == nil {
		return strings.Title(action.Name) + " " + action.Syntax + " what?"
	}

	return game.finalizeItemAction(item, target, action)
}

func (game *Game) finalizeItemAction(item Itemer, target Itemer, action *Action) string {
	msg, parent := item.OnAction(action, target)

	if parent != item.Basic().Location {
		game.ChangeParent(item, parent)
	}

	return msg + game.Rooms[game.Location].OnAction(action)
}

//ChangeParent - move item to the new owner
func (game *Game) ChangeParent(item Itemer, parentName string) {
	room := game.CurrentRoom().BasicRoom()

	//remove from previous owner
	if item.Basic().Location == "inventory" {
		for idx, check := range game.Inventory {
			if item.Basic().Name == check.Basic().Name {
				game.Inventory = append(game.Inventory[:idx], game.Inventory[idx+1:]...)
				break
			}
		}
	} else {
		parent, idx := findParent(item, append(room.Items, game.Inventory...))

		if parent != nil {
			parent.Items = append(parent.Items[:idx], parent.Items[idx+1:]...)
		} else if idx > -1 {
			room := game.CurrentRoom().BasicRoom()
			room.Items = append(room.Items[:idx], room.Items[idx+1:]...)
		} //or it was an object without parent
	}

	//add to new owner
	if parentName == "inventory" {
		game.Inventory = append(game.Inventory, item)
	} else if parentName != "" {
		_, parent, _ := findTarget([]string{parentName}, append(visibleItems(room.Items, false, true), game.Inventory...), true)
		if parent != nil {
			parent.Basic().Items = append(parent.Basic().Items, item)
		} else {
			room := game.Rooms[parentName]
			if room != nil {
				room.BasicRoom().Items = append(room.BasicRoom().Items, item)
			}
		}
	}

	item.Basic().Location = parentName
}

//Navigate -
func (game *Game) Navigate(location string, dir string) string {
	room := game.CurrentRoom()

	if location != "" {
		canLeave, msg := room.LeaveRoom(dir)
		if !canLeave {
			if msg == "" {
				return "You can't go that way."
			}
			return msg
		}
		room = game.Rooms[location]
		if room.BasicRoom().Locked != "" {
			return room.BasicRoom().Locked
		}
		game.Location = location

		return msg + room.EnterRoom(location)
	}
	return "You can't go that way."
}

//ShowInventory -
func (game *Game) ShowInventory() string {
	msg := []string{}

	if len(game.Inventory) == 0 {
		return "You have nothing."
	}

	for i, item := range game.Inventory {
		if i == 0 {
			msg = append(msg, "You have ")
		} else if i == len(game.Inventory)-1 {
			msg = append(msg, " and ")
		} else {
			msg = append(msg, ", ")
		}
		msg = append(msg, item.NameWithArticle())
	}

	return strings.Join(msg, "")
}

//Help - displays keywords
func (game *Game) Help() string {
	return `Navigation: (n)orth, (s)outh, (e)ast, (w)est.
Useful verbs: (l)ook, e(x)amine, take, open, close, put _ on _, unlock _ with _
Characters: ask _ about _, give _ to _`
}
