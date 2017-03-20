package server

import (
	"encoding/json"
	"net/http"

	"regexp"
	"strings"

	"github.com/bobziuchkovski/cue"
)

// CheckReq innehåller information för klient requests från APIn
type CheckReq struct {
	ID        int    `json:"id"`
	CommandID int    `json:"command_id"`
	Command   string `json:"command"`
	Save      bool   `json:"save"`
}

// UpdateClient tar hand om klient releterade requests från APIn
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(cue.Fields{
		"Type":   "CheckHandler",
		"Remote": r.RemoteAddr,
	}).Info("Got a request to internal API")

	decoder := json.NewDecoder(r.Body)
	req := CheckReq{ID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err, "Internal error")
		return
	}
	defer r.Body.Close()

	log.WithFields(cue.Fields{
		"Type":    "CheckHandler",
		"ID":      req.ID,
		"Command": req.Command,
		"Save":    req.Save,
	}).Debug("ClientReq information")

	if req.ID == -1 {
		OutputJson(w, ErrorResp{Error: true, Message: "You need to supply an ID"})
		return
	}

	cl := clients.GetClientByID(req.ID)
	if cl == nil {
		OutputJson(w, ErrorResp{Error: true, Message: "Can't find a client in cache with this ID"})
		return
	}

	if strings.HasPrefix(req.Command, "ping") {
		ports := "3333"
		re := regexp.MustCompile("-port=\"?([\\d,-]+)\"?")
		p := re.FindStringSubmatch(req.Command)
		if len(p) >= 2 {
			ports = p[1]
		}

		cl.Lock()
		ip := cl.ip
		cl.Unlock()

		resp := Ping(ip, ports)
		OutputJson(w, resp)

		if req.Save {

		}
		return
	}
}
