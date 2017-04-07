package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CheckHandler struct{}

func (c CheckHandler) Serve(r *HTTPHandler) {
	if r.Request.ID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply an ID"})
		return
	}

	cl := r.Server.GetClients().GetClientByID(r.Request.ID)
	if cl == nil {
		r.Output(APIResponse{Error: true, Message: "Can't find a client in cache with this ID"})
		return
	}

	if r.Request.CommandID != -1 {
		if r.Request.Save {
			resp := []string{}
			for _, che := range cl.GetChecksByCommandID(r.Request.CommandID) {
				resp = append(resp, cl.SendCheck(r.Server, che))
			}
			b, _ := json.Marshal(resp)
			r.Output(APIResponse{Error: false, Message: fmt.Sprintf("%s", b)})
			return
		} else {
			resp := []string{}
			for _, che := range cl.GetChecksByCommandID(r.Request.CommandID) {
				response, err := c.SendCheck(r, cl, che.GetCommand().GetCommand())
				if err != nil {
					r.Output(APIResponse{Error: true, Message: response})
				}
				resp = append(resp, response)
			}
			b, _ := json.Marshal(resp)
			r.Output(APIResponse{Error: false, Message: fmt.Sprintf("%s", b)})
		}
	} else {
		if r.Request.Command == "" {
			r.Output(APIResponse{Error: true, Message: "You need to supply a command"})
			return
		}
		resp, err := c.SendCheck(r, cl, r.Request.Command)
		if err != nil {
			r.Output(APIResponse{Error: true, Message: resp})
		}
		r.Output(APIResponse{Error: false, Message: resp})
	}
}

func (c CheckHandler) SendCheck(r *HTTPHandler, cl *Client, command string) (string, error) {
	if strings.HasPrefix(command, "ping") {
		ports := "3333"
		p := re.FindStringSubmatch(r.Request.Command)
		if len(p) >= 2 {
			ports = p[1]
		}

		pingResponse, err := cl.Ping(ports)
		b, _ := json.Marshal(pingResponse)
		return string(b), err
	}

	resp, err := cl.SendMessage(command)
	if err != nil {
		return "Can't connect to client", err
	}
	return resp, nil
}
