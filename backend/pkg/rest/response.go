package rest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	gzPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(ioutil.Discard)
		},
	}
)

type paging struct {
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}

//APIResponse model
type APIResponse struct {
	w http.ResponseWriter

	gzip      bool
	iserror   bool
	usepaging bool
	paging    paging
	message   string
	errors    interface{}
	data      interface{}
}

//NewAPIResponse new instance of APIResponse
func NewAPIResponse(w http.ResponseWriter, r *http.Request) *APIResponse {
	return &APIResponse{
		w:    w,
		gzip: r != nil && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip"),
	}
}

//Error flag as error
func (res *APIResponse) Error(message string, err interface{}) *APIResponse {
	res.iserror = true
	res.errors = err
	res.message = message
	return res
}

//Paging flag as data list
func (res *APIResponse) Paging(totaldata, totalpage int64) *APIResponse {
	res.usepaging = true
	res.paging = paging{TotalData: totaldata, TotalPage: totalpage}
	return res
}

//Payload add data payload to resposne
func (res *APIResponse) Payload(payload interface{}) *APIResponse {
	res.data = payload
	return res
}

//Respond to HTTP client
func (res *APIResponse) Respond(statusCode int) {
	d, _ := res.MarshalJSON()
	res.w.Header().Set("Content-Type", "application/json")
	if res.iserror {
		fmt.Println(res.message, res.errors)
	}

	if res.gzip {
		res.w.Header().Set("Content-Encoding", "gzip") //Set encoding to gzip
		res.w.WriteHeader(statusCode)
		gzwriter := gzPool.Get().(*gzip.Writer)
		gzwriter.Reset(res.w)
		gzwriter.Write(d)

		gzwriter.Close()
		gzPool.Put(gzwriter)
	} else {
		res.w.Header().Set("Content-Length", strconv.Itoa(len(d)))
		res.w.WriteHeader(statusCode)
		res.w.Write(d)
	}
}

//MarshalJSON serialize APIResponse to JSON
func (res *APIResponse) MarshalJSON() ([]byte, error) {
	if res.iserror {
		return json.Marshal(struct {
			Message string      `json:"message"`
			Errors  interface{} `json:"errors,omitempty"`
		}{
			Message: res.message,
			Errors:  res.errors,
		})
	}

	if res.usepaging {
		return json.Marshal(struct {
			Paging paging      `json:"paging"`
			Data   interface{} `json:"data"`
		}{
			Paging: res.paging,
			Data:   res.data,
		})
	}

	if res.data != nil {
		return json.Marshal(struct {
			Data interface{} `json:"data,omitempty"`
		}{Data: res.data})
	}

	return nil, nil
}

//UnmarshalJSON deserialize JSON into APIResponse
func (res *APIResponse) UnmarshalJSON(jsonbyte []byte) error {
	js := make(map[string]interface{})
	err := json.Unmarshal(jsonbyte, &js)
	if err != nil {
		return err
	}

	pg, ok := js["paging"]
	if ok {
		paging := pg.(map[string]interface{})
		res.usepaging = true
		res.paging.TotalData = int64(paging["total_data"].(float64))
		res.paging.TotalPage = int64(paging["total_page"].(float64))
	}

	message, ok := js["message"]
	if ok && message != "" {
		res.iserror = true
		res.message = message.(string)
	}

	errors, ok := js["errors"]
	if ok {
		res.errors = errors
	}

	data, ok := js["data"]
	if ok {
		res.data = data
	}

	return nil
}
