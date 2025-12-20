package botc

import (
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RoleType string

const (
	Townsfolk RoleType = "townsfolk"
	Outsider  RoleType = "outsider"
	Minion    RoleType = "minion"
	Demon     RoleType = "demon"
	Traveller RoleType = "traveller"
	Fabled    RoleType = "fabled"
	Loric     RoleType = "loric"
)

var RoleTypeOrder = []RoleType{
	Townsfolk,
	Outsider,
	Minion,
	Demon,
	Traveller,
	Fabled,
	Loric,
}

func (rt RoleType) Name() string {
	return cases.Title(language.English).String(string(rt))
}

func (rt RoleType) Id() string {
	return string(rt)
}

type Alignment string

const (
	Good   Alignment = "Good"
	Evil   Alignment = "Evil"
	Either Alignment = "Either"
	None   Alignment = "None"
)

func (a Alignment) String() string {
	return string(a)
}

type Role struct {
	Id                 string            `json:"id"`
	Name               string            `json:"name"`
	Edition            Edition           `json:"edition"`
	ImageUrls          []string          `json:"image"`
	Team               RoleType          `json:"team"`
	Ability            string            `json:"ability"`
	FirstNightOrder    int               `json:"firstNight"`
	FirstNightReminder string            `json:"firstNightReminder"`
	OtherNightOrder    int               `json:"otherNight"`
	OtherNightReminder string            `json:"otherNightReminder"`
	GlobalReminders    []string          `json:"remindersGlobal"`
	ReminderTokens     []string          `json:"reminders"`
	AltersSetup        bool              `json:"setup"`
	Flavour            string            `json:"flavor"`
	Special            []Special         `json:"special"`
	Jinxes             map[string]string `json:"jinxes"`
}

func (r *Role) Type() string {
	return r.Team.Name()
}

func (r *Role) Alignment() Alignment {
	switch r.Team {
	case Townsfolk, Outsider:
		return Good
	case Minion, Demon:
		return Evil
	case Traveller:
		return Either
	default:
		return None
	}
}

func (r *Role) ImageUrl(a Alignment) string {
	d := r.Alignment()
	// short-circuit: only select if more than 1 image or the requested isn't our default
	if len(r.ImageUrls) > 1 && a != d {
		switch d {
		// if we're a Traveller
		case Either:
			// and the requested alignment is Good or Evil grab it or fall out of the switch
			switch a {
			case Good:
				return r.ImageUrls[1]
			case Evil:
				return r.ImageUrls[2]
			}
		// if we're Good or Evil and in the switch, the requested alignment isn't our default
		case Good, Evil:
			// if it's neither Either nor None, return the second (opposite)
			if !(a == Either || a == None) {
				return r.ImageUrls[1]
			}
		}
	}
	// most cases should fall out of the switches to their default images
	return r.ImageUrls[0]
}

func (r *Role) AbilityText() string {
	pos := strings.Index(r.Ability, "[")
	return r.Ability[:pos-1]
}

func (r *Role) Setup() string {
	pos := strings.Index(r.Ability, "[")
	return r.Ability[pos+1 : len(r.Ability)-1]
}

func (r *Role) JinxWith(o *Role) (string, bool) {
	text, found := r.Jinxes[o.Id]
	return text, found
}

func (r *Role) Wiki() string {
	return fmt.Sprintf(
		"https://wiki.bloodontheclocktower.com/%s",
		url.PathEscape(strings.ReplaceAll(r.Name, " ", "_")),
	)
}

func (r *Role) ToMap() map[string]any {
	m := make(map[string]any)
	rt := reflect.TypeOf(r)
	rv := reflect.ValueOf(r)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get("json")
		n := strings.Split(t, ",")[0]
		v := rv.Field(i)
		switch n {
		case "special":
			sm := make([]map[string]any, len(r.Special))
			for i, s := range r.Special {
				sm[i] = s.ToMap()
			}
			smv := reflect.ValueOf(sm)
			m[n] = smv.Interface()
		case "jinxes":
			js := make([]map[string]string, len(r.Jinxes))
			i := 0
			for id, reason := range r.Jinxes {
				j := make(map[string]string)
				j["id"] = id
				j["reason"] = reason
				js[i] = j
				i++
			}
			m[n] = js
		default:
			m[n] = v.Interface()
		}
	}
	return m
}

func NewRole(m map[string]any) (Role, error) {
	var r Role
	id, err := extractRoleId(m)
	if err != nil {
		return r, err
	}
	r.Id = id

	name, err := extractRequiredString("name", m)
	if err != nil {
		return r, err
	}
	r.Name = name

	edition, err := extractRoleEdition(m)
	if err != nil {
		return r, err
	}
	r.Edition = edition

	image, err := extractRoleImageUrls(m)
	if err != nil {
		return r, err
	}
	r.ImageUrls = image

	team, err := extractRoleTeam(m)
	if err != nil {
		return r, err
	}
	r.Team = team

	ability, err := extractRequiredString("ability", m)
	if err != nil {
		return r, err
	}
	r.Ability = ability

	firstNight, _, err := extractInt("firstNight", m)
	if err != nil {
		return r, err
	}
	r.FirstNightOrder = firstNight

	firstNightReminder, _, err := extractString("firstNightReminder", m)
	if err != nil {
		return r, err
	}
	r.FirstNightReminder = firstNightReminder

	otherNight, _, err := extractInt("otherNight", m)
	if err != nil {
		return r, err
	}
	r.OtherNightOrder = otherNight

	otherNightReminder, _, err := extractString("otherNightReminder", m)
	if err != nil {
		return r, err
	}
	r.OtherNightReminder = otherNightReminder

	globalReminders, _, err := extractStringSlice("remindersGlobal", m)
	if err != nil {
		return r, err
	}
	r.GlobalReminders = globalReminders

	reminders, _, err := extractStringSlice("reminders", m)
	if err != nil {
		return r, err
	}
	r.ReminderTokens = reminders

	alters, err := extractRequiredBool("setup", m)
	if err != nil {
		return r, err
	}
	r.AltersSetup = alters

	flavour, _, err := extractString("flavor", m)
	if err != nil {
		return r, err
	}
	r.Flavour = flavour

	specials, err := extractSpecial(m)
	if err != nil {
		return r, err
	}
	r.Special = specials

	jinxes, err := extractJinxes(m)
	if err != nil {
		return r, err
	}
	r.Jinxes = jinxes

	return r, nil
}

func extractRoleId(m map[string]any) (string, error) {
	key := "id"
	id, err := extractRequiredString(key, m)
	if err != nil {
		return "", err
	}
	pos := strings.Index(id, "_")
	if pos > 0 {
		return id[:pos], nil
	} else {
		return id, nil
	}
}

func extractRoleEdition(m map[string]any) (Edition, error) {
	key := "edition"
	ed, err := extractRequiredString(key, m)
	if err != nil {
		return Edition(""), err
	}
	edition := Edition(ed)
	if slices.Contains(EditionOrder, edition) {
		return edition, nil
	} else {
		return Edition(""), NewIllegalValueForEnumError(key, edition, EditionOrder)
	}
}

func extractRoleImageUrls(m map[string]any) ([]string, error) {
	key := "image"
	urls, ok, err := extractStringSlice(key, m)
	if !ok {
		return []string{}, &RequiredFieldMissingError{key}
	}
	if err != nil {
		return urls, err
	}
	for _, u := range urls {
		_, err := url.Parse(u)
		if err != nil {
			return []string{}, fmt.Errorf("invalid url: %w", err)
		}
	}
	return urls, nil
}

func extractRoleTeam(m map[string]any) (RoleType, error) {
	key := "team"
	team, err := extractRequiredString(key, m)
	if err != nil {
		return RoleType(""), err
	}
	rt := RoleType(team)
	if slices.Contains(RoleTypeOrder, rt) {
		return rt, nil
	} else {
		return RoleType(""), NewIllegalValueForEnumError(key, rt, RoleTypeOrder)
	}
}
