package botc

import (
	"encoding/json"
	"log"
)

type Roster struct {
	Author         string `json:"author"`
	Name           string `json:"name"`
	Almanac        string `json:"almanac"`
	Characters     []*Role
	CharacterIndex map[string]*Role
}

func (r *Roster) UnmarshalJSON(b []byte) error {
	r.CharacterIndex = make(map[string]*Role)
	var items []any
	err := json.Unmarshal(b, &items)
	if err != nil {
		return err
	}
	for _, i := range items {
		switch I := i.(type) {
		case map[string]any:
			switch I["id"] {
			case "_meta":
				r.Author = I["author"].(string)
				r.Name = I["name"].(string)
				r.Almanac = I["almanac"].(string)
			default:
				role, err := NewRole(I)
				if err != nil {
					log.Printf("%v", role)
					return err
				}
				r.Characters = append(r.Characters, &role)
				r.CharacterIndex[role.Id] = &role
			}
		}
	}
	return nil
}

func (r *Roster) MarshalJSON() ([]byte, error) {
	items := make([]any, len(r.Characters)+1)
	meta := make(map[string]string)
	meta["author"] = r.Author
	meta["name"] = r.Name
	meta["almanac"] = r.Almanac
	items[0] = meta
	for i, c := range r.Characters {
		items[i+1] = c.ToMap()
	}
	var bytes []byte
	bytes, err := json.Marshal(items)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}
