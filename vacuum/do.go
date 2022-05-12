package vacuum

import "encoding/json"

type request struct {
	ID     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type rawResponse struct {
	PartnerID string          `json:"partner_id"`
	ID        int             `json:"id"`
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Result    json.RawMessage `json:"Result"`
}

func (v *Vacuum) do(method string, params json.RawMessage, rsp interface{}) error {
	v.id++

	req := &request{
		Method: method,
		ID:     v.id,
		Params: params,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	rspBytes, err := v.client.Send(reqBytes)
	if err != nil {
		return err
	}

	rr := &rawResponse{}

	err = json.Unmarshal(rspBytes, rr)
	if err != nil {
		return err
	}

	// TODO: Check rawResponse

	err = json.Unmarshal(rr.Result, rsp)
	if err != nil {
		return err
	}

	return nil
}
