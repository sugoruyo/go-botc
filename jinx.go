package botc

type Jinx struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func extractJinxes(m map[string]any) (map[string]string, error) {
	raw, ok, err := extractSlice("jinxes", m)
	if !ok {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	jinxes := make(map[string]string)

	for _, item := range raw {
		var id, reason string
		for k, v := range item.(map[string]any) {
			val := v.(string)
			switch k {
			case "id":
				id = val
			case "reason":
				reason = val
			}
		}
		jinxes[id] = reason
	}
	return jinxes, nil
}

func JinxesToMap(js []Jinx) map[string]string {
	m := make(map[string]string)
	for _, j := range js {
		m[j.Id] = j.Reason
	}
	return m
}

func JinxesFromMap(m map[string]string) []Jinx {
	js := make([]Jinx, len(m))
	i := 0
	for id, reason := range m {
		js[i] = Jinx{
			Id:     id,
			Reason: reason,
		}
		i++
	}
	return js
}
