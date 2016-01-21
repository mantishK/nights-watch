package middleware

import (
	"net/http"

	al "plivo/nights-watch/log"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type Adapter func(http.Handler) http.Handler

func Adapt(hf http.HandlerFunc, adapters ...Adapter) httprouter.Handle {
	var h http.Handler
	for _, adapter := range adapters {
		if h == nil {
			h = adapter(hf)
		} else {
			h = adapter(h)
		}
	}
	handler := addParameterToContext(h)
	return handler
}

func addParameterToContext(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		context.Set(r, "params", p)
		h.ServeHTTP(w, r)
	}
}

func AccessLog() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			al.Access("URL " + r.URL.String() + " IP: " + r.RemoteAddr)
			h.ServeHTTP(w, r)
		})
	}
}

func ClearContext() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			context.Clear(r)
		})
	}
}
