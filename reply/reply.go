package reply

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	ae "plivo/nights-watch/apperror"
	"plivo/nights-watch/log"
)

type public interface {
	Public() interface{}
}

//Reply with 200 ok
func OK(w http.ResponseWriter, data interface{}) {
	if obj, ok := data.(public); ok {
		data = obj.Public()
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		Err(w, ae.JsonEncode("", err))
	}
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, string(jsonData))
}

//Reply with given error
func Err(w http.ResponseWriter, e *ae.Error) {
	log.Err(*e)
	log.Err(string(debug.Stack()))
	jsonData, err := json.Marshal(e)
	if err != nil {
		fmt.Fprint(w, "Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.HttpStatus)
	fmt.Fprintln(w, string(jsonData))
}
