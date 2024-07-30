package handlers

import (
	"net/http"
)


func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " message\n"))
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " stats\n"))
}
