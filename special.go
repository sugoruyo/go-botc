package botc

import (
	"reflect"
	"slices"
	"strings"
)

var specialTypeEnum []string = []string{
	"selection",
	"ability",
	"signal",
	"vote",
	"reveal",
	"player",
}

var specialNameEnum []string = []string{
	"grimoire",
	"pointing",
	"ghost-votes",
	"distribute-roles",
	"bag-disabled",
	"bag-duplicate",
	"multiplier",
	"hidden",
	"replace-character",
	"player",
	"card",
	"open-eyes",
}

var specialTimeEnum []string = []string{
	"pregame",
	"day",
	"night",
	"firstNight",
	"firstDay",
	"otherNight",
	"otherDay",
}

var specialGlobalEnum []string = []string{
	"townsfolk",
	"outsider",
	"minion",
	"demon",
	"traveller",
	"dead",
}

type Special struct {
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Time   string `json:"time,omitempty"`
	Global string `json:"global,omitempty"`
	Value  any    `json:"value,omitempty"`
}

func (s *Special) StringValue() (string, bool) {
	switch v := s.Value.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}

func (s *Special) NumValue() (int, bool) {
	switch v := s.Value.(type) {
	case int:
		return v, true
	default:
		return 0, false
	}
}

func (s *Special) ToMap() map[string]any {
	m := make(map[string]any)
	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(s)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t := f.Tag.Get("json")
		n := strings.Split(t, ",")[0]
		v := rv.Field(i)
		m[n] = v.Interface()
	}
	return m
}

func extractSpecial(m map[string]any) ([]Special, error) {
	raw, ok, err := extractSlice("special", m)
	if !ok {
		return []Special{}, nil
	}
	if err != nil {
		return []Special{}, err
	}
	specials := make([]Special, 0)

	for _, item := range raw {
		var special Special

		for k, v := range item.(map[string]any) {
			switch k {
			case "type":
				if slices.Contains(specialTypeEnum, v.(string)) {
					special.Type = v.(string)
				}
			case "name":
				if slices.Contains(specialNameEnum, v.(string)) {
					special.Name = v.(string)
				}
			case "time":
				if slices.Contains(specialTimeEnum, v.(string)) {
					special.Time = v.(string)
				}
			case "global":
				if slices.Contains(specialGlobalEnum, v.(string)) {
					special.Global = v.(string)
				}
			case "value":
				special.Value = v
			}
		}

		specials = append(specials, special)
	}
	return specials, nil
}
