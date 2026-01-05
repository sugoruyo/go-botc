package botc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"
)

type NightOrdered interface {
	GetName() string
	GetFirstNightOrder() int
	GetOtherNightsOrder() int
}

type ScriptMeta struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Author     string   `json:"author"`
	Logo       string   `json:"logo"`
	HideTitle  bool     `json:"hideTitle"`
	Background string   `json:"background"`
	Almanac    string   `json:"almanac"`
	Bootlegger []string `json:"bootlegger"`
	FirstNight []string `json:"firstNight"`
	OtherNight []string `json:"otherNight"`
}

type Script struct {
	Meta                 ScriptMeta       `json:"meta"`
	CustomCharacters     []Role           `json:"custom"`
	OriginalCharacterIds []string         `json:"original"`
	Index                map[string]*Role `json:"index"`
}

func (s *Script) Author() string {
	if s.Meta.Author != "" {
		return s.Meta.Author
	} else {
		return "Unknown"
	}
}

func (s *Script) PopulateIndex(r Roster) []string {
	notFound := make([]string, 0)
	for _, o := range s.OriginalCharacterIds {
		orig, ok := r.CharacterIndex[o]
		if !ok {
			notFound = append(notFound, o)
		}
		s.Index[o] = orig
	}
	for _, c := range s.CustomCharacters {
		s.Index[c.Id] = &c
	}
	return notFound
}

func (s *Script) GetCharacter(id string) Role {
	return *s.Index[id]
}

func (s *Script) FirstNight() []NightOrdered {
	order := make([]NightOrdered, 0)
	for e := range firstNight {
		order = append(order, e)
	}
	for _, role := range s.Index {
		if role.FirstNightOrder != -1 {
			order = append(order, role)
		}
	}
	slices.SortFunc(order, func(a, b NightOrdered) int {
		return a.GetFirstNightOrder() - b.GetFirstNightOrder()
	})
	return order
}

func (s *Script) OtherNights() []NightOrdered {
	order := make([]NightOrdered, 0)
	for e := range otherNights {
		order = append(order, e)
	}
	for _, role := range s.Index {
		if role.OtherNightOrder != -1 {
			order = append(order, role)
		}
	}
	slices.SortFunc(order, func(a, b NightOrdered) int {
		return a.GetOtherNightsOrder() - b.GetOtherNightsOrder()
	})
	return order
}

func (s *Script) OfficialToolUrl() string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	jdata, err := json.Marshal(s)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := gz.Write(jdata); err != nil {
		log.Fatalln(err)
	}
	gz.Close()
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	esc := url.QueryEscape(enc)
	scriptUrl := fmt.Sprintf("https://script.bloodontheclocktower.com/?script=%s", esc)
	return scriptUrl
}

func (s *Script) UnmarshalJSON(data []byte) error {
	raw := make([]any, 0)
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	for _, v := range raw {
		switch vt := v.(type) {
		case string:
			s.OriginalCharacterIds = append(s.OriginalCharacterIds, strings.ReplaceAll(vt, "_", ""))
		case map[string]any:
			if vt["id"] == "_meta" {
				s.Meta.Id = vt["id"].(string)
				s.Meta.Name = checkNilString(vt["name"])
				s.Meta.Author = checkNilString(vt["author"])
				s.Meta.Almanac = checkNilString(vt["almanac"])
				s.Meta.Bootlegger = checkNilStringS(vt["bootlegger"])
				s.Meta.FirstNight = checkNilStringS(vt["firstNight"])
				s.Meta.OtherNight = checkNilStringS(vt["otherNight"])
			} else {
				if len(vt) == 1 {
					s.OriginalCharacterIds = append(s.OriginalCharacterIds, normaliseId(vt["id"].(string)))
				} else {
					role, err := NewRole(vt)
					if err != nil {
						return err
					}
					s.CustomCharacters = append(s.CustomCharacters, role)
				}
			}
		}
	}
	s.Index = make(map[string]*Role)
	return nil
}

func (s *Script) MarshalJSON() ([]byte, error) {
	raw := make([]any, 0)
	meta := make(map[string]any)
	meta["id"] = s.Meta.Id
	meta["name"] = s.Meta.Name
	if s.Meta.Author != "" {
		meta["author"] = s.Meta.Author
	}
	if s.Meta.Logo != "" {
		meta["logo"] = s.Meta.Logo
	}
	meta["hideTitle"] = s.Meta.HideTitle
	if s.Meta.Background != "" {
		meta["background"] = s.Meta.Background
	}
	if s.Meta.Almanac != "" {
		meta["almanac"] = s.Meta.Almanac
	}
	if len(s.Meta.Bootlegger) > 0 {
		meta["bootlegger"] = s.Meta.Bootlegger
	}
	meta["firstNight"] = s.Meta.FirstNight
	meta["otherNight"] = s.Meta.OtherNight
	raw = append(raw, meta)
	for _, o := range s.OriginalCharacterIds {
		raw = append(raw, o)
	}
	for _, c := range s.CustomCharacters {
		raw = append(raw, c.ToMap())
	}
	var bytes []byte
	bytes, err := json.Marshal(raw)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func checkNilString(x any) string {
	if x != nil {
		return x.(string)
	} else {
		return ""
	}
}

func checkNilStringS(x any) []string {
	if x != nil {
		return convertStringS(x)
	} else {
		return make([]string, 0)
	}
}

func convertStringS(x any) []string {
	xs := x.([]any)
	r := make([]string, len(xs))
	for i, s := range xs {
		r[i] = s.(string)
	}
	return r
}
