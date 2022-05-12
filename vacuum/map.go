package vacuum

import "encoding/json"

func (v *Vacuum) Map() (string, error) {
	m := &json.RawMessage{}
	err := v.do("get_map_v1", nil, m)
	if err != nil {
		return "", err
	}

	return string(*m), nil
}
