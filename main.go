package main

import (
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./Statics"))
	http.Handle("/statics/", http.StripPrefix("/statics", fileServer))
	http.HandleFunc("/", Index)
	http.HandleFunc("/api/NewToken", NewToken)
	http.HandleFunc("/api/GetInfo", GetInfo)
	http.HandleFunc("/image/", Image)
	http.ListenAndServe(":80", nil)
}
