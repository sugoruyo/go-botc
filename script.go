package botc

const (
	Dusk       string = "dusk"
	MinionInfo string = "minioninfo"
	DemonInfo  string = "demoninfo"
	Dawn       string = "dawn"
)

type ScriptMeta struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Author     string   `json:"author"`
	Logo       string   `json:"logo"`
	HideTitle  bool     `json:"hideTitle"`
	Background string   `json:"background"`
	Almanac    string   `json:"almanac"`
	Bootlegger string   `json:"bootlegger"`
	FirstNight []string `json:"firstNight"`
	OtherNight []string `json:"otherNight"`
}

type Script struct {
	Meta       ScriptMeta
	Characters []Role
}
