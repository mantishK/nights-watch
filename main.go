package main

import (
	"net/http"
	"os"

	"plivo/nights-watch/handler"
	appLog "plivo/nights-watch/log"
	m "plivo/nights-watch/middleware"

	"github.com/julienschmidt/httprouter"
)

func main() {
	server := GetHandler()
	http.ListenAndServe(":8080", server)
}

func GetHandler() http.Handler {
	router := httprouter.New()

	appLog.Init(os.Stdout, os.Stdout)
	apiNS := "/api"

	router.Handle("POST", apiNS+"/outbound/sms", m.Adapt(handler.SendSMS, m.AccessLog(), m.ClearContext()))
	router.Handle("POST", apiNS+"/inbound/sms", m.Adapt(handler.StopNo, m.AccessLog(), m.ClearContext()))

	return router

}
