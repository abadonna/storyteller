package engine

//Action -
type Action struct {
	Name             string
	IsItemRequired   bool
	IsTargetRequired bool
	IsActorRequired  bool
	IsActorTarget    bool
	IsTopicRequired  bool

	IsPredefined       bool //should be false for any user defined actions!
	Syntax             string
	DefaultTopicAnswer string
}

//LOOK action
var LOOK = &Action{
	Name:         "look",
	IsPredefined: true}

//INVENTORY action
var INVENTORY = &Action{
	Name:         "inventory",
	IsPredefined: true}

//EXAMINE action
var EXAMINE = &Action{
	Name:           "examine",
	IsItemRequired: true,
	IsPredefined:   true}

//OPEN action
var OPEN = &Action{
	Name:           "open",
	IsItemRequired: true,
	IsPredefined:   true}

//CLOSE action
var CLOSE = &Action{
	Name:           "close",
	IsItemRequired: true,
	IsPredefined:   true}

//TAKE action
var TAKE = &Action{
	Name:           "take",
	IsItemRequired: true,
	IsPredefined:   true}

//PUT action
var PUT = &Action{
	Name:           "put",
	IsItemRequired: true,
	IsPredefined:   true}

//USE action
var USE = &Action{
	Name:           "use",
	Syntax:         "on",
	IsItemRequired: true,
	IsPredefined:   true}

//ASK action
var ASK = &Action{
	Name:               "ask",
	IsActorRequired:    true,
	IsTopicRequired:    true,
	DefaultTopicAnswer: "\"I don't know much about that.\"",
	IsPredefined:       true}

//UNLOCK action
var UNLOCK = &Action{
	Name:             "unlock",
	Syntax:           "with",
	IsItemRequired:   true,
	IsTargetRequired: true,
	IsPredefined:     true}

//GIVE action
var GIVE = &Action{
	Name:               "give",
	Syntax:             "to",
	IsItemRequired:     true,
	IsActorTarget:      true,
	IsTopicRequired:    true,
	DefaultTopicAnswer: "\"I don't don't need it.\"",
	IsPredefined:       true}
