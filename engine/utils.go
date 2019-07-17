package engine

import (
	"regexp"
	"strings"
)

var ignore = map[string]bool{
	"a":     true,
	"an":    true,
	"the":   true,
	"in":    true,
	"into":  true,
	"on":    true,
	"onto":  true,
	"upon":  true,
	"from":  true,
	"to":    true,
	"about": true,
	"with":  true}

func notifyAboutVisibleItems(items []Itemer, location string) string {
	msg := []string{}
	actors := []string{}
	visible := visibleItems(items, true, false)

	for i, item := range visible {
		actor, ok := item.(Actor)
		if ok { //let's keep actors separatly from items
			actors = append(actors, "\n"+actor.BasicPerson().NameEx+" is here.")
			continue
		}

		if i == 0 {
			msg = append(msg, "\nYou see ")
		} else if i == len(visible)-1 {
			msg = append(msg, " and ")
		} else {
			msg = append(msg, ", ")
		}
		msg = append(msg, item.NameWithArticle())
	}

	if len(msg) > 0 {
		msg = append(msg, location+".")
	}

	return strings.Join(msg, "") + strings.Join(actors, "")
}

func visibleItems(items []Itemer, ignoreDecoration bool, nested bool) []Itemer {
	result := []Itemer{}
	for _, i := range items {
		item := i.Basic()
		if ignoreDecoration && item.IsDecoration {
			continue
		}
		if item.IsVisible && !item.IsDisabled {
			result = append(result, i)
			if nested {
				result = append(result, visibleItems(item.Items, ignoreDecoration, nested)...)
			}
		}
	}
	return result
}

func findItemsInList(words []string, items []Itemer) ([]Itemer, []string) {

	target := []string{}

	for _, word := range words {
		if ignore[word] {
			continue
		}

		isEnd := true
		possible := []Itemer{}
		for _, i := range items {
			item := i.Basic()
			match, _ := regexp.MatchString("(^| )"+word+"( |$)",
				strings.ToLower(item.Vocab+" "+item.Name))

			if match {
				if isEnd {
					target = append(target, word)
				}
				isEnd = false
				possible = append(possible, i)
			}
		}

		if isEnd {
			break
		}
		items = possible
	}

	if len(target) == 0 {
		return nil, target
	}

	return items, target
}

func findTarget(words []string, items []Itemer, optional bool) (string, Itemer, []string) {
	possible, object := findItemsInList(words, items)

	if len(object) == 0 {
		if optional {
			return "", nil, nil
		}
		return "You don't see any " + strings.Join(words, " ") + " here.", nil, nil
	}

	if len(possible) > 1 {
		input := strings.Join(object, " ")
		list := ""
		for i, item := range possible {
			if i > 0 && i < len(possible)-1 {
				list += ", "
			} else if i == len(possible)-1 {
				list += " or "
			}
			list += item.Basic().Name
		}
		return "What " + input + " do you mean: " + list + "?", nil, nil
	}

	return "", possible[0], words[len(object):]
}

func findParent(item Itemer, items []Itemer) (*Item, int) {
	for i, p := range items {
		if p == item {
			return nil, i
		}
		parent := p.Basic()
		for idx, check := range parent.Items {
			if check == item {
				return parent, idx
			}
			nested, idx := findParent(item, check.Basic().Items)
			if nested != nil {
				return nested, idx
			}
		}

	}
	return nil, -1
}

func filterActors(items []Itemer) []Actor {
	actors := []Actor{}
	for _, item := range items {
		var i interface{} = item
		actor, ok := i.(Actor)
		if ok {
			actors = append(actors, actor)
		}
	}
	return actors
}

func findActor(words []string, actors []Actor) (string, Actor, []string) {
	items := []Itemer{}
	for _, actor := range actors {
		items = append(items, actor)
	}

	possible, object := findItemsInList(words, items)

	if len(object) == 0 {
		return "You don't see this person here.", nil, nil
	}

	if len(possible) > 1 {
		list := ""
		for i, item := range possible {
			if i > 0 && i < len(possible)-1 {
				list += ", "
			} else if i == len(possible)-1 {
				list += " or "
			}
			list += item.Basic().Name
		}
		return "Whom do you mean: " + list + "?", nil, nil
	}

	return "", possible[0].(Actor), words[len(object):]
}

func findTopics(words []string, actor Actor) []*Topic {
	result := []*Topic{}
	topics := actor.BasicPerson().Topics

	for _, word := range words {
		if ignore[word] {
			continue
		}

		isEnd := true
		for _, topic := range topics {
			match, _ := regexp.MatchString("(^| )"+word+"( |$)", topic.Vocab)

			if match {
				isEnd = false
				result = append(result, topic)
			}
		}

		if len(result) > 0 {
			if isEnd {
				break
			}
			topics = result
		}
	}

	return result
}

func findSyntax(words []string, syntax string) ([]string, bool) {

	for i, word := range words {
		if word == syntax {
			return words[i+1:], true
		}
	}

	return nil, false
}
