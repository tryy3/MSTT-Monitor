package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type CommandReq struct {
	Type     string `json:"type"`
	ID       int    `json:"id"`
	ClientID int    `json:"clientid"`
}

func getCommand(id int) (string, error) {
	getCommandStmt, err := db.Prepare("SELECT command FROM commands WHERE id=?")
	if err != nil {
		return "", err
	}
	defer getCommandStmt.Close()
	var command string
	err = getCommandStmt.QueryRow(id).Scan(&command)
	if err != nil {
		return "", nil
	}
	return command, nil
}

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
		panic(err.Error())
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

func UpdateCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	req := CommandReq{ID: -1, ClientID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	switch strings.ToLower(req.Type) {
	case "delete":
		// Kolla om kommandot existerar, för säkerhetskull.
		command, err := getCommand(req.ID)
		if err != nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Something is wrong with the database."})
			return
		}
		if command != "" {
			OutputJson(w, ErrorResp{Error: true, Message: "You can't remove existing command."})
			return
		}

		for i := clients.Length(); i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}
			if req.ClientID != -1 {
				cl.Lock()
				cid := cl.clientID
				cl.Unlock()
				if cid == req.ClientID {
					cl.RemoveByCommandID(req.ID)
					break
				}
			}
			cl.RemoveByCommandID(req.ID)
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

		for i := clients.Length(); i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}
			for j := cl.Length(); j >= 0; j-- {
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
		OutputJson(w, ErrorResp{Error: false, Message: "Successfully updated the command for all clients."})
		return
	case "insert":
		if req.ClientID == -1 {
			OutputJson(w, ErrorResp{Error: true, Message: "You need to supply a ClientID"})
			return
		}
		command, err := getCommand(req.ID)
		if err != nil {
			OutputJson(w, ErrorResp{Error: true, Message: "Something is wrong with the database."})
			return
		}
		if command == "" {
			OutputJson(w, ErrorResp{Error: true, Message: "Can't find the command with this id."})
			return
		}
		for i := clients.Length(); i >= 0; i-- {
			cl := clients.Get(i)
			if cl == nil {
				continue
			}

			cl.Lock()
			clientID := cl.clientID
			cl.Unlock()

			if clientID == req.ClientID {
				cl.Lock()
				groups := cl.groups
				cl.Unlock()

				check, err := getCommandInfo(req.ID, req.ClientID, groups)
				if err != nil {
					OutputJson(w, ErrorResp{Error: true, Message: err.Error()})
				}
				cl.Add(check)
				OutputJson(w, ErrorResp{Error: false, Message: "Added the new command to the client."})
				return
			}
		}
		OutputJson(w, ErrorResp{Error: true, Message: "Can't find the ClientID"})
		return
	default:
		OutputJson(w, ErrorResp{Error: true, Message: "Invalid UpdateType"})
		return
	}
}
