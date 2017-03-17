package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bobziuchkovski/cue"
)

// ClientReq innehåller information för klient requests från APIn
type ClientReq struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

// UpdateClient tar hand om klient releterade requests från APIn
func UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(cue.Fields{
		"Type":   "UpdateClient",
		"Remote": r.RemoteAddr,
	}).Info("Got a request to internal API")

	decoder := json.NewDecoder(r.Body)
	req := ClientReq{ID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err, "Internal error")
		return
	}
	defer r.Body.Close()

	log.WithFields(cue.Fields{
		"Type":    "UpdateClient",
		"ReqType": req.Type,
		"ID":      req.ID,
	}).Debug("ClientReq information")

	if req.ID == -1 {
		OutputJson(w, ErrorResp{Error: true, Message: "You need to supply an ID"})
		return
	}

	switch strings.ToLower(req.Type) {
	case "update":
		var ip string
		query := fmt.Sprintf("SELECT ip FROM clients WHERE id=%d", req.ID)

		err := db.QueryRow(query).Scan(&ip)
		if err != nil {
			if err == sql.ErrNoRows {
				OutputJson(w, ErrorResp{Error: true, Message: "ID does not exists"})
				return
			} else {
				OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
				log.Error(err, "Error getting information from database")
				return
			}
		}

		cl := clients.GetClientByID(req.ID)

		log.WithFields(cue.Fields{
			"Type":   "UpdateCommand",
			"Client": cl == nil,
		}).Debug("Client exists")

		if cl == nil {
			newCL := &Client{ip: ip, clientID: req.ID}
			clients.Add(newCL)
			log.Debug("Client added")
			OutputJson(w, ErrorResp{Error: false, Message: "Client added"})
			break
		}

		cl.Lock()
		cl.ip = ip
		cl.Unlock()
		log.Debug("Client updated")
		OutputJson(w, ErrorResp{Error: false, Message: "Client updated"})
		break
	case "delete":
		rem := clients.RemoveByClientID(req.ID)
		if !rem {
			OutputJson(w, ErrorResp{Error: true, Message: "ClientID does not exist in cache"})
			return
		}

		OutputJson(w, ErrorResp{Error: false, Message: "Removed the client from cache"})
		break
	default:
		OutputJson(w, ErrorResp{Error: true, Message: "Invalid UpdateType"})
		break
	}
}
