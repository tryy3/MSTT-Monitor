package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bobziuchkovski/cue"
)

// GroupReq innehåller information för grupp requests från APIn
type GroupReq struct {
	Type     string `json:"type"`
	Name     string `json:"name"` // Group name
	ClientID int    `json:"clientid"`
}

// DelGroup tar hand om att ta bort en grupp från en klient
func DelGroup(cl *Client, group string) bool {
	cl.Lock()
	groups := cl.groups
	cl.Unlock()

	b := false

	for i, g := range groups {
		if g == group {
			cl.Lock()
			cl.groups = append(cl.groups[:i], cl.groups[i+1:]...)
			cl.Unlock()
			b = true
			break
		}
	}

	if !b {
		return b
	}

	cl.Lock()
	for i := len(cl.checks) - 1; i >= 0; i-- {
		if cl.checks[i].gruppNamn == group {
			cl.checks = append(cl.checks[:i], cl.checks[i+1:]...)
		}
	}
	cl.Unlock()

	return b
}

// AddGroup tar hand om att lägga till en grupp till en klient
func AddGroup(w http.ResponseWriter, cl *Client, group string) {
	getGroupStmt, err := db.Prepare("SELECT command_id, next_check, stop_error FROM groups")
	if err != nil {
		log.Error(err, "Can't prepare getGroupStmt")
		OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
		return
	}
	defer getGroupStmt.Close()

	getCheckStmt, err := db.Prepare("SELECT id, timestamp, checked, error, finished FROM checks WHERE command_id=? AND client_id=? ORDER BY timestamp DESC")
	if err != nil {
		log.Error(err, "Can't prepare getCheckStmt")
		OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
		return
	}
	defer getCheckStmt.Close()

	groupRows, err := getGroupStmt.Query()
	if err != nil {
		log.Error(err, "Can't query getGroupStmt")
		OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
		return
	}
	defer groupRows.Close()

	var (
		command_id int
		next_check int64
		stop_error bool

		id         int64
		timestamp  string
		checked    bool
		checkError bool
		finished   bool
	)

	cl.Lock()
	clientID := cl.clientID
	cl.Unlock()

	e := false

	for groupRows.Next() {
		err := groupRows.Scan(&command_id, &next_check, &stop_error)
		if err != nil {
			log.Error(err, "Can't scan groupRows")
			e = true
			continue
		}

		cmd, err := getCommand(command_id)
		if err != nil {
			log.Error(err, "Can't get command")
			e = true
			continue
		}

		t := time.Now()
		err = getCheckStmt.QueryRow(command_id, clientID).Scan(&id, &timestamp, &checked, &checkError, &finished)
		if err != nil {
			if err == sql.ErrNoRows {
				// Make a "fake" check
				checked = false
				checkError = false
				finished = true
				id = -1
			} else {
				log.Error(err, "Can't get last check")
				e = true
				continue
			}
		} else {
			t, err = CreateTimestamp(timestamp, next_check)
			if err != nil {
				log.Error(err, "Can't create timestamp")
				e = true
				continue
			}
		}

		check := &Check{
			command:       cmd,
			gruppNamn:     group,
			commandID:     command_id,
			nextCheck:     next_check,
			pastID:        id,
			nextTimestamp: t,
			checked:       checked,
			err:           checkError,
			failErr:       stop_error,
			done:          finished,
		}

		cl.Lock()
		cl.groups = append(cl.groups, group)
		cl.Unlock()
		cl.Add(check)
	}
	if e {
		OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
	} else {
		OutputJson(w, ErrorResp{Error: false, Message: "Inserted the group with the client"})
	}
}

// UpdateGroup tar hand om grupp releterade requests från APIn
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(cue.Fields{
		"Type":   "UpdateGroup",
		"Remote": r.RemoteAddr,
	}).Info("Got a request to internal API")

	decoder := json.NewDecoder(r.Body)
	req := GroupReq{ClientID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.WithFields(cue.Fields{
		"Type":     "UpdateGroup",
		"ReqType":  req.Type,
		"Name":     req.Name,
		"ClientID": req.ClientID,
	}).Debug("ClientReq information")

	if req.Name == "" {
		OutputJson(w, ErrorResp{Error: true, Message: "You need to supply a Name"})
	}

	switch strings.ToLower(req.Type) {
	case "insert":
		if req.ClientID == -1 {
			OutputJson(w, ErrorResp{Error: true, Message: "You need to supply a ClientID"})
			return
		}

		cl := clients.GetClientByID(req.ClientID)
		if cl == nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Can't find a client with this ID."})
			return
		}

		cl.Lock()
		groups := cl.groups
		cl.Unlock()

		for _, group := range groups {
			if group == req.Name {
				OutputJson(w, ErrorResp{Error: true, Message: "Client already belongs to this group in cache."})
				return
			}
		}
		AddGroup(w, cl, req.Name)
		break
	case "update":
		id, err := strconv.Atoi(req.Name)
		if err != nil {
			log.Error(err, "Can't convert the ID to an int.")
			OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
			return
		}

		var (
			command_id int
			next_check int64
			stop_error bool
		)

		err = db.QueryRow("SELECT command_id, next_check, stop_error FROM groups WHERE id=?", id).Scan(&command_id, &next_check, &stop_error)
		if err != nil {
			log.Error(err, "Can't query for group info")
			OutputJson(w, ErrorResp{Error: true, Message: "Internal error"})
			return
		}

		for i := clients.Length() - 1; i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}

			for j := cl.Length() - 1; i >= 0; i-- {
				ch := cl.Get(j)
				if ch == nil {
					continue
				}

				ch.Lock()
				if ch.commandID == command_id {
					ch.nextCheck = next_check
					ch.failErr = stop_error
				}
				ch.Unlock()
			}
		}
		OutputJson(w, ErrorResp{Error: false, Message: "Updated the groups in cache."})
		return
	case "delete":
		if req.ClientID != -1 {
			cl := clients.GetClientByID(req.ClientID)
			if cl == nil {
				OutputJson(w, ErrorResp{Error: true, Message: "Can't find a client with this ID."})
				return
			}

			b := DelGroup(cl, req.Name)
			if !b {
				OutputJson(w, ErrorResp{Error: true, Message: "Client does not belong to this group in cache."})
			} else {
				OutputJson(w, ErrorResp{Error: false, Message: "Group has now been removed from this client in cache."})
			}
		} else {
			for i := clients.Length(); i >= 0; i-- {
				cl := clients.Get(i)
				if cl == nil {
					continue
				}

				DelGroup(cl, req.Name)
			}
			OutputJson(w, ErrorResp{Error: false, Message: "Group has now been removed from all clients in cache"})
		}
		break
	default:
		OutputJson(w, ErrorResp{Error: true, Message: "Invalid UpdateType"})
		return
	}
}
