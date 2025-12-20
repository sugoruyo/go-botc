package botc

type Edition string

const (
	TB       Edition = "tb"
	SNV      Edition = "snv"
	BMR      Edition = "bmr"
	CAROUSEL Edition = "carousel"
	FABLED   Edition = "fabled"
	LORIC    Edition = "loric"
)

var EditionName = map[Edition]string{
	TB:       "Trouble Brewing",
	SNV:      "Sects & Violets",
	BMR:      "Bad Moon Rising",
	CAROUSEL: "Carousel",
	FABLED:   "Fabled",
	LORIC:    "Loric",
}

var EditionOrder = []Edition{
	TB,
	SNV,
	BMR,
	CAROUSEL,
	FABLED,
	LORIC,
}

func (e Edition) Name() string {
	return EditionName[e]
}

func (e Edition) Id() string {
	return (string(e))
}
