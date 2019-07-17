package game

import (
	"math/rand"
	"storyteller/engine"
	"strings"
)

//SampleGame just sample
type SampleGame struct {
	engine.Game
}

///////////////////////////////CUSTOM ROOM SAMPLE///////////////////////////////
type outerRoom struct {
	engine.Room
}

func (room *outerRoom) LeaveRoom(dir string) (bool, string) {
	if dir == "s" {
		return false, "It's not a good idea to walk that way."
	}
	return true, "You entered the darkness...\n"
}

func (room *outerRoom) OnAction(action *engine.Action) string {
	if action.Name == "sleep" {
		return "Zzzz..."
	}
	events := [3]string{
		"\nYou hear some noizes from the cave.",
		"\nHot winds are blowing from the desert.",
		""}
	return room.BasicRoom().OnAction(action) + events[rand.Intn(len(events))]
}

///////////////////////////////CUSTOM ITEM SAMPLE///////////////////////////////
type bottle struct {
	engine.Item
	isEmpty bool
}

func (item *bottle) OnAction(action *engine.Action, target engine.Itemer) (string, string) {
	if action.Name == "drink" {
		if item.isEmpty {
			return "It's empty.", item.Location
		}
		item.isEmpty = true
		item.Name = "empty bottle"
		return "You drink some water.", item.Location
	}
	return item.Basic().OnAction(action, target)
}

///////////////////////////////////////////////////////////////////////////////
type skull struct {
	engine.Item
	game *engine.Game
}

func (item *skull) OnAction(action *engine.Action, target engine.Itemer) (string, string) {
	if action == engine.TAKE {
		if len(item.game.Rooms["Cave"].BasicRoom().Items[0].Basic().Items) < 2 {
			item.game.IsFinished = true
			return "As you lift the skull, a volley of poisonous arrows is shot from the walls! You try to dodge the arrows, but they take you by surprise!\nYou are dead.", ""
		}
	}
	return item.Basic().OnAction(action, target)
}

///////////////////////////////CUSTOM ACTOR SAMPLE//////////////////////////////
type witch struct {
	engine.Person
	game *engine.Game
}

func (person *witch) OnTopic(topic *engine.Topic, action *engine.Action) string {

	if action.Name == "give" {
		if topic != nil && topic.Vocab == "gold skull" {
			person.game.IsFinished = true
			return "\"Yes! Thank you!\"\nGame Over."
		}
	}

	msg := ""

	if action.Name == "ask" {

		if topic != nil && topic.Vocab == "name" {
			person.NameEx = "Mellisa the witch"
			person.Name = "Melissa"
		} else if topic != nil && strings.Contains(topic.Vocab, "box") {
			item := engine.Item{
				Name:       "key",
				Desc:       "A small key that should help to unlock something.",
				IsPickable: true,
				IsVisible:  true}
			person.game.ChangeParent(&item, "inventory")
			msg = "\nYou obtained a small key!"
		}
	}

	return person.BasicPerson().OnTopic(topic, action) + msg
}

///////////////////////////////////////////////////////////////////////////////

//Sample game instance
func Sample() *SampleGame {
	context := new(SampleGame)
	context.Game.Rooms = make(map[string]engine.Spacer)
	context.Game.Actions = []engine.Action{
		{Name: "sleep"},
		{Name: "drink", IsItemRequired: true}}

	context.Game.Rooms["Outside cave"] = &outerRoom{
		engine.Room{
			Desc:  "[[img=https://i.imgur.com/ar18tWi.jpg]]You're standing in the bright sunlight just outside of a large, dark, foreboding cave, which lies to the north. Desert lies to the south.",
			North: "Cave",
			South: "Desert",
			Items: []engine.Itemer{
				&witch{
					game: &context.Game,
					Person: engine.Person{
						Topics: []*engine.Topic{
							{
								Action:  "ask",
								Vocab:   "pedestal",
								Answers: []string{"\"Yes, examine it. The skull should be somewere on in.\""}},
							{
								Action:        "ask",
								Vocab:         "box key lock",
								Answers:       []string{"\"Ah, box... Here, take the key.\""},
								RepeatAnswers: []string{"\"You have the key, right?\""}},
							{
								Action: "give",
								Vocab:  "gold skull"},
							{
								Action:  "give",
								Vocab:   "bottle",
								Answers: []string{"She drinks the water.\n\"Nice, thanks. But I need the skull.\""}},

							{
								Action:  "ask",
								Vocab:   "help skull",
								Answers: []string{"\"Bring me the skull from the cave! But be careful, be sure to put something on the pedestal before taking the skull!\""}},
							{
								Action:  "ask",
								Vocab:   "name",
								Answers: []string{"\"I'm Melissa, the local witch. And I need your help.\""}}},
						Item: engine.Item{
							Name:      "mysterious woman",
							Vocab:     "girl woman witch melissa",
							Desc:      "[[img=https://images.sex.com/images/pinporn/2016/10/20/620/16760796.jpg]]You see a mysterious woman in dark clothes.\n\"Hey, can we talk?\", she asks.",
							IsVisible: true},
						NameEx: "Mysterious beautiful woman"}},
				&engine.Item{
					Name:         "cave",
					Vocab:        "dark large foreboding",
					Desc:         "It's a very dark cave.",
					IsVisible:    true,
					IsDecoration: true}}}}

	context.Game.Rooms["Cave"] = &engine.Room{
		Desc:  "You're inside a dark and musty cave. Sunlight pours in from a passage to the south.",
		South: "Outside cave",
		Items: []engine.Itemer{
			&engine.Item{
				Name:      "pedestal",
				Desc:      "There is an ancient pedestal inside the cave.",
				IsSurface: true,
				IsVisible: true,
				Location:  "Cave",
				Items: []engine.Itemer{
					&skull{
						game: &context.Game,
						Item: engine.Item{
							Name:       "gold skull",
							IsPickable: true,
							Location:   "pedestal"}}}},
			//-------------------------------------//
			&engine.Item{
				Name:        "box",
				Desc:        "Old wooden box.",
				IsContainer: true,
				IsVisible:   true,
				IsLocked:    true,
				KeyName:     "key",
				Location:    "Cave",
				Items: []engine.Itemer{
					&bottle{
						Item: engine.Item{
							Name:       "bottle",
							IsPickable: true,
							Location:   "box"}},
					&engine.Item{
						Name:       "steel sword",
						IsPickable: true,
						Location:   "box"},
					&engine.Item{
						Name:       "silver sword",
						IsPickable: true,
						Location:   "box"}}}}}

	context.Game.Location = "Outside cave"

	//var _ engine.Itemer = &engine.Person{}
	return context
}

//Help - sample help
func (sample *SampleGame) Help() string {
	return sample.Game.Help() + `
	
	Sample custom verbs: sleep, drink`
}
