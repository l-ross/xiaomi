package vacuum

import "encoding/json"

func (v *Vacuum) Map() (string, error) {
	rsp := &json.RawMessage{}
	err := v.do("get_map_v1", nil, rsp)
	if err != nil {
		return "", err
	}

	return string(*rsp), nil
}
