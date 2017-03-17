package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"database/sql"

	"github.com/bobziuchkovski/cue"
)

// CommandReq innehåller information för command APIn
type CommandReq struct {
	Type      string `json:"type"`
	ID        int    `json:"id"`
	GroupName string `json:"group"`
}

// getCommandBase hämtar information om ett kommand
func getCommandBase(group string, id int) (string, int64, bool, error) {
	var (
		next_check int64
		stop_error bool
	)

	err := db.QueryRow("SELECT next_check, stop_error FROM groups WHERE command_id=? AND group_name=?", id, group).Scan(&next_check, &stop_error)
	if err != nil {
		return "", -1, false, err
	}

	command, err := getCommand(id)
	if err != nil {
		return "", -1, false, err
	}
	return command, next_check, stop_error, nil
}

// getCommand hämtar ett command från databasen
func getCommand(id int) (string, error) {
	var command string
	log.WithFields(cue.Fields{
		"ID": id,
	}).Debug("Get command")
	err := db.QueryRow("SELECT command FROM commands WHERE id=?", id).Scan(&command)
	if err != nil {
		return "", nil
	}
	return command, nil
}

// getCommandInfo hämtar information releterat till ett kommand
func getCommandInfo(id int, clientID int, groups []string) (check *Check, err error) {
	command, err := getCommand(id)
	if err != nil {
		return nil, err
	}
	if command == "" {
		return nil, errors.New("Command does not exists.")
	}

	getGroupCommandStmt, err := db.Prepare("SELECT next_check, stop_error FROM groups WHERE command_id=? AND group_name=?")
	if err != nil {
		return nil, err
	}
	defer getGroupCommandStmt.Close()

	getCheckStmt, err := db.Prepare("SELECT id, timestamp, checked, error, finished FROM checks WHERE command_id=? AND client_id=? ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer getCheckStmt.Close()

	var (
		next_check int64
		stop_error bool
	)

	for _, group := range groups {
		err = getGroupCommandStmt.QueryRow(id, group).Scan(&next_check, &stop_error)
		if err != nil {
			continue
		}

		var (
			timestamp  string
			checked    bool
			checkError bool
			done       bool
			pastID     int64
		)

		err = getCheckStmt.QueryRow(id, clientID).Scan(&pastID, &timestamp, &checked, &checkError, &done)
		if err != nil {
			return &Check{
				command:   command,
				commandID: id,
				nextCheck: -1,
				pastID:    -1,
			}, nil
		}

		t, err := CreateTimestamp(timestamp, next_check)
		if err != nil {
			panic(err.Error())
		}

		return &Check{
			command:       command,
			gruppNamn:     group,
			commandID:     id,
			nextCheck:     next_check,
			pastID:        pastID,
			nextTimestamp: t,
			checked:       checked,
			err:           checkError,
			failErr:       stop_error,
			done:          done,
		}, nil
	}
	return nil, errors.New("Client does not belong to a group that have access to this command ID.")
}

// UpdateCommand tar hand om command releterade requests från APIn
func UpdateCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.WithFields(cue.Fields{
		"Type":   "UpdateCommand",
		"Remote": r.RemoteAddr,
	}).Info("Got a request to internal API")

	decoder := json.NewDecoder(r.Body)
	req := CommandReq{ID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err, "Internal error")
		return
	}
	defer r.Body.Close()

	log.WithFields(cue.Fields{
		"Type":      "UpdateCommand",
		"ReqType":   req.Type,
		"ID":        req.ID,
		"GroupName": req.GroupName,
	}).Debug("CommandReq information")

	if req.ID == -1 {
		OutputJson(w, ErrorResp{Error: true, Message: "You need to supply an ID"})
		return
	}

	switch strings.ToLower(req.Type) {
	case "delete":
		for i := clients.Length() - 1; i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}
			if req.GroupName != "" {
				cl.Lock()
				groups := cl.groups
				cl.Unlock()
				for _, g := range groups {
					if g == req.GroupName {
						cl.RemoveByCommandID(req.ID)
						break
					}
				}
			} else {
				cl.RemoveByCommandID(req.ID)
			}
		}

		OutputJson(w, ErrorResp{Error: false, Message: "Successfully removed the command."})
		return
	case "update":
		command, err := getCommand(req.ID)
		if err != nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Something is wrong with the database."})
			return
		}

		if command == "" {
			OutputJson(w, ErrorResp{Error: true, Message: "Can't find the command with this id."})
			return
		}

		for i := clients.Length() - 1; i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}

			for j := cl.Length() - 1; j >= 0; j-- {
				ch := cl.Get(j)
				if ch == nil {
					continue
				}

				ch.Lock()
				if ch.commandID == req.ID {
					ch.command = command
				}
				ch.Unlock()
			}
		}

		OutputJson(w, ErrorResp{Error: false, Message: "Successfully updated the command for all clients"})
		return
	case "insert":
		if req.GroupName == "" {
			OutputJson(w, ErrorResp{Error: true, Message: "You need to supply a group name"})
			return
		}

		command, next_check, stop_error, err := getCommandBase(req.GroupName, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				OutputJson(w, ErrorResp{Error: false, Message: "Nothing to insert."})
				return
			}
			OutputJson(w, ErrorResp{Error: true, Message: "Something is wrong with the database"})
			log.Error(err, "Internal error")
			return
		}

		if command == "" {
			OutputJson(w, ErrorResp{Error: true, Message: "Can't find the command with this id"})
			return
		}

		getCheckStmt, err := db.Prepare("SELECT id, timestamp, checked, error, finished FROM checks WHERE command_id=? AND client_id=? ORDER BY timestamp DESC")
		if err != nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Something is wrong with the database"})
			return
		}
		defer getCheckStmt.Close()

		var (
			timestamp  string
			checked    bool
			checkError bool
			done       bool
			pastID     int64
		)

		for i := clients.Length(); i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}

			cl.Lock()
			groups := cl.groups
			id := cl.clientID
			cl.Unlock()

			for _, g := range groups {
				if g == req.GroupName {
					err = getCheckStmt.QueryRow(req.ID, id).Scan(&pastID, &timestamp, &checked, &checkError, &done)
					if err != nil {
						cl.Add(&Check{
							command:   command,
							nextCheck: next_check,
							failErr:   stop_error,
							gruppNamn: g,
							commandID: req.ID,
						})
					} else {
						t, err := CreateTimestamp(timestamp, next_check)
						if err != nil {
							log.Error(err, "Failed creating a timestamp.")
						}
						cl.Add(&Check{
							command:       command,
							nextCheck:     next_check,
							failErr:       stop_error,
							gruppNamn:     g,
							commandID:     req.ID,
							pastID:        pastID,
							nextTimestamp: t,
							checked:       checked,
							err:           checkError,
							done:          done,
						})
					}
				}
			}

			OutputJson(w, ErrorResp{Error: false, Message: "Added the new command to the client."})
			return
		}

		OutputJson(w, ErrorResp{Error: true, Message: "Can't find the ClientID"})
		break
	default:
		OutputJson(w, ErrorResp{Error: true, Message: "Invalid UpdateType"})
		break
	}
}
