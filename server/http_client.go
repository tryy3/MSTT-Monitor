package server

type ClientHandler struct{}

func (ClientHandler) Insert(r *HTTPHandler) {
	r.Output(APIResponse{Error: true, Message: "Unsupported method, please refer to the documentation"})
}

func (ClientHandler) Update(r *HTTPHandler) {
	if r.Request.ID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply an ID"})
		return
	}

	c, err := r.Server.GetDatabase().GetClient(r.Request.ID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	cl := r.Server.GetHandler().GetClientByID(r.Request.ID)

	if cl == nil {
		client := NewClient(c)
		r.Server.GetHandler().AddClient(client)
		r.Output(APIResponse{Error: false, Message: "Client added"})
	} else {
		cl.SetIP(c.IP)
		r.Output(APIResponse{Error: false, Message: "Client updated"})
	}
}

func (ClientHandler) Delete(r *HTTPHandler) {
	if r.Request.ID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply an ID"})
		return
	}

	b := r.Server.GetHandler().RemoveClientByID(r.Request.ID)
	if !b {
		r.Output(APIResponse{Error: true, Message: "Client ID does not exists in cache"})
		return
	} else {
		r.Output(APIResponse{Error: false, Message: "Removed the client from cache"})
	}
}
