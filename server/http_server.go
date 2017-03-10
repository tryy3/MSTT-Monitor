package server

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type ErrorResp struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func StartWebServer() {
	http.HandleFunc("/update/command", UpdateCommand)
	http.HandleFunc("/update/client", UpdateClient)
	http.HandleFunc("/update/group", UpdateGroup)
	http.HandleFunc("/check", UpdateHandler)

	http.ListenAndServe(":8080", nil)
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func OutputJson(w http.ResponseWriter, i interface{}) {
	out, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
