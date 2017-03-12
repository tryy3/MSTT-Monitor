package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ClientReq innehåller information för klient requests från APIn
type ClientReq struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

// UpdateClient tar hand om klient releterade requests från APIn
func UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	req := ClientReq{ID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	switch strings.ToLower(req.Type) {
	case "update":
		var ip string
		query := fmt.Sprintf("SELECT ip FROM clients WHERE id=%d", req.ID)

		err := db.QueryRow(query).Scan(&ip)
		if err != nil {
			if err == sql.ErrNoRows {
				OutputJson(w, ErrorResp{Error: true, Message: "ID does not exists."})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		cl := clients.GetClientByID(req.ID)
		if cl == nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Can't find client in cache."})
		}

		cl.Lock()
		cl.ip = ip
		cl.Unlock()
	case "insert":
	case "delete":
		var exists bool
		query := fmt.Sprintf("SELECT exists (SELECT namn FROM clients WHERE id=%d)", req.ID)

		err := db.QueryRow(query).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if !exists {
			OutputJson(w, ErrorResp{Error: true, Message: "ClientID does not exist in database"})
		}

		rem := clients.RemoveByClientID(req.ID)
		if !rem {
			OutputJson(w, ErrorResp{Error: true, Message: "ClientID does not exist in cache"})
		}

		OutputJson(w, ErrorResp{Error: false, Message: "Removed the client from cache"})
	default:
		OutputJson(w, ErrorResp{Error: true, Message: "Invalid UpdateType"})
		return
	}
}
