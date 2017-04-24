package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

type APIRequest struct {
	Type      string `json:"type"`
	Command   string `json:"command"`
	GroupName string `json:"group_name"`
	ID        int64  `json:"id"`
	CommandID int64  `json:"command_id"`
	GroupID   int64  `json:"group_id"`
	Save      bool   `json:"save"`
}

type HTTPHandler struct {
	Request APIRequest
	Server  *Server
	w       http.ResponseWriter
	r       *http.Request
}

func (h HTTPHandler) Output(a APIResponse) {
	out, err := json.Marshal(a)
	if err != nil {
		h.Server.GetLogger().Error(err, "Internal error")
		http.Error(h.w, "Something went wrong internal, contact IT support", http.StatusInternalServerError)
		return
	}
	h.w.Header().Set("Content-Type", "application/json")
	h.w.Write(out)
}

type APIHandler interface {
	Insert(*HTTPHandler)
	Update(*HTTPHandler)
	Delete(*HTTPHandler)
}

type APIResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type HTTPServer struct {
	Handlers map[string]APIHandler
	Check    CheckHandler
	Server   *Server
}

func (h *HTTPServer) Start() {
	h.Handlers["command"] = CommandHandler{}
	h.Handlers["client"] = ClientHandler{}
	h.Handlers["group"] = GroupHandler{}
	h.Handlers["alert"] = GroupHandler{}
	h.Check = CheckHandler{}

	http.Handle("/", h)
	err := http.ListenAndServe(h.Server.GetConfig().APIAdress, nil)
	if err != nil {
		h.Server.GetLogger().Error(err, "Can't start Web API")
	}
}

func (h HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiHandler := &HTTPHandler{
		Server: h.Server,
		w:      w,
		r:      r,
	}

	decoder := json.NewDecoder(r.Body)
	req := APIRequest{ID: -1, CommandID: -1, GroupID: -1}
	err := decoder.Decode(&req)
	if err != nil {
		h.Server.GetLogger().Error(err, "Internal error")
		apiHandler.Output(APIResponse{Error: true, Message: "Internal error when parsing request body"})
		return
	}
	defer r.Body.Close()
	apiHandler.Request = req

	if strings.HasPrefix(r.URL.Path, "/check") {
		h.Check.Serve(apiHandler)
	}
	if strings.HasPrefix(r.URL.Path, "/update") {
		path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(path) < 2 {
			apiHandler.Output(APIResponse{Error: true, Message: "Invalid Path"})
			return
		}
		handler, ok := h.Handlers[path[1]]
		if !ok {
			apiHandler.Output(APIResponse{Error: true, Message: "Invalid Path"})
			return
		}

		switch strings.ToLower(req.Type) {
		case "insert":
			handler.Insert(apiHandler)
			return
		case "update":
			handler.Update(apiHandler)
			return
		case "delete":
			handler.Delete(apiHandler)
			return
		default:
			apiHandler.Output(APIResponse{Error: true, Message: "Invalid Request type"})
			return
		}
	}
}
