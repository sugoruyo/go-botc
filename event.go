package botc

type Event string

const (
	Dusk       Event = "dusk"
	MinionInfo Event = "minioninfo"
	DemonInfo  Event = "demoninfo"
	Dawn       Event = "dawn"
)

var eventNames = map[Event]string{
	Dusk:       "Dusk",
	MinionInfo: "Minion Info",
	DemonInfo:  "Demon Info",
	Dawn:       "Dawn",
}

var firstNight = map[Event]int{
	Dusk:       1,
	MinionInfo: 19,
	DemonInfo:  23,
	Dawn:       77,
}

var otherNights = map[Event]int{
	Dusk: 1,
	Dawn: 96,
}

func (e Event) GetName() string {
	return eventNames[e]
}

func (e Event) String() string {
	return string(e)
}

func (e Event) GetFirstNightOrder() int {
	order, found := firstNight[e]
	if !found {
		return -1
	}
	return order
}

func (e Event) GetOtherNightsOrder() int {
	order, found := otherNights[e]
	if !found {
		return -1
	}
	return order
}
