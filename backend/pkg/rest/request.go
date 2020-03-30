package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type apirequest struct {
	Data json.RawMessage `json:"data"`
}

// ParseBody of HTTP request into obj struct
// Payload must be inside {"data": ...}
func ParseBody(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return nil
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	req := apirequest{}
	err = json.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	return json.Unmarshal(req.Data, obj)
}
