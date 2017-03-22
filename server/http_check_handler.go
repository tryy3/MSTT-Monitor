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
		pingError := false
		if strings.Contains(req.Command, "-error") {
			pingError = true
		}
		re := regexp.MustCompile("-port=\"?([\\d,-]+)\"?")
		p := re.FindStringSubmatch(req.Command)
		if len(p) >= 2 {
			ports = p[1]
		}

		cl.Lock()
		ip := cl.ip
		cl.Unlock()

		resp := Ping(ip, ports, pingError)
		m, err := json.Marshal(resp)
		if err != nil {
			log.Error(err, "Something went wrong when marshaling struct.")
			OutputJson(w, ErrorResp{Error: true, Message: "Something went wrong when turning struct into string, check logs."})
			return
		}
		OutputJson(w, ErrorResp{Error: false, Message: string(m)})

		if req.Save {
			if req.CommandID == -1 {
				log.Error(nil, "Missing command ID.")
				OutputJson(w, ErrorResp{Error: true, Message: "Missing Command ID for saving to mysql."})
				return
			}

			stmt, err := db.Prepare("INSERT INTO checks(command_id, client_id, response, checked, error, finished) VALUES (?,?,?,1,?,1)")
			if err != nil {
				log.Error(err, "Mysql error when inserting a check.")
				OutputJson(w, ErrorResp{Error: true, Message: "SQL error."})
				return
			}
			e := false
			if resp.Error != "" {
				e = true
			}

			_, err = stmt.Exec(req.CommandID, req.ID, string(m), e)
			if err != nil {
				log.Error(err, "Mysql error when inserting a check.")
				OutputJson(w, ErrorResp{Error: true, Message: "SQL error."})
				return
			}
			OutputJson(w, ErrorResp{Error: true, Message: "Succesffully saved the response."})
		}
		return
	}
}
