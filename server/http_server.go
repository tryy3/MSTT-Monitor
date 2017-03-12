package server

import (
	"encoding/json"
	"net/http"
)

// ErrorResp innehåller information om response för webbsidan
type ErrorResp struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// StartWebServer startar Web APIn
func StartWebServer() {
	http.HandleFunc("/update/command", UpdateCommand)
	http.HandleFunc("/update/client", UpdateClient)
	//http.HandleFunc("/update/group", UpdateGroup)
	//http.HandleFunc("/check", UpdateHandler)

	err := http.ListenAndServe(config.APIAdress, nil)
	if err != nil {
		log.Error(err, "Can't start Web API")
	}
}

// OutputJson är en enkel funktion som printar ut en interface i json format
func OutputJson(w http.ResponseWriter, i interface{}) {
	out, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
