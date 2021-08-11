package httpUtils

import (
	"encoding/json"
	"github.com/danielgom/bookstore_utils-go/errors"
	"net/http"
)

func ResponseJson(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	res, e := json.Marshal(body)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, e = w.Write(res)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(status)
}

func ResponseError(w http.ResponseWriter, err errors.RestErr) {
	ResponseJson(w, err.Status(), err)
}
