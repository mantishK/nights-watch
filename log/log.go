package log

import (
	"io"
	"log"

	ae "plivo/nights-watch/apperror"
)

var errLogger *log.Logger
var accessLogger *log.Logger

func Init(outErr io.Writer, outAccess io.Writer) {
	errLogger = log.New(outErr, "Error: ", log.Ldate|log.Ltime)
	accessLogger = log.New(outAccess, "Access: ", log.Ldate|log.Ltime)
}

func Err(i interface{}) {
	if errLogger != nil {
		if aErr, ok := i.(ae.Error); ok {
			errLogger.Println(aErr.Log)
		} else {
			errLogger.Println(i)
		}
	}
}

func Access(i interface{}) {
	if accessLogger != nil {
		accessLogger.Println(i)
	}
}
