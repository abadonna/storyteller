package engine

import (
	"math/rand"
)

//Actor interface
type Actor interface {
	Itemer
	OnTopic(topic *Topic, action *Action, item Itemer) string
	BasicPerson() *Person
}

//Person - actor, based on item
type Person struct {
	Item
	Topics         []*Topic
	DefaultAnswers map[string][]string
	Hello          string
	NameEx         string
	IsKnown        bool
}

//Topic for interaction with Actor
type Topic struct {
	Action         string
	IsUsed         bool
	Vocab          string
	Answers        []string
	RepeatAnswers  []string
	IsItemConsumed bool
}

//OnTopic - actor's reaction on topic actions
func (person *Person) OnTopic(topic *Topic, action *Action, item Itemer) string {
	if topic == nil {
		count := len(person.DefaultAnswers[action.Name])
		if count > 0 {
			return person.DefaultAnswers[action.Name][rand.Intn(count)]
		}
		return action.DefaultTopicAnswer
	}

	if topic.IsUsed {
		count := len(topic.RepeatAnswers)
		if count > 0 {
			return topic.RepeatAnswers[rand.Intn(count)]
		}
	}

	topic.IsUsed = true
	count := len(topic.Answers)
	if count > 0 {
		return topic.Answers[rand.Intn(count)]
	}

	return "ERROR: topic has no answers!"

}

//BasicPerson - return base struct
func (person *Person) BasicPerson() *Person {
	return person
}

//NameWithArticle -
func (person *Person) NameWithArticle() string {
	if person.AName != "" {
		return person.Item.NameWithArticle()
	}
	return person.Name
}

//OnAction override for actor
func (person *Person) OnAction(action *Action, target Itemer) (string, string) {
	if action == EXAMINE {
		return person.Examine(), person.Location
	}

	return "I don't know how to " + action.Name + " " + person.Name + ".", person.Location
}

//Examine person
func (person *Person) Examine() string {
	if person.Desc != "" {
		return person.Desc
	}
	if person.NameEx != "" {
		return person.NameEx
	}

	return "You see " + person.NameWithArticle() + "."
}
