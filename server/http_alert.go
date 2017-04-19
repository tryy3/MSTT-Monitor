package server

type AlertHandler struct{}

func (AlertHandler) Insert(r *HTTPHandler) {
	if r.Request.ID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply an alert option ID"})
		return
	}

	alert, err := r.Server.GetDatabase().GetAlertOptionsByID(r.Request.ID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	for cl := range r.Server.GetClients().IterClients() {
		for ch := range cl.IterChecks() {
			for al := range ch.IterAlerts() {
				if al.GetID() == r.Request.ID {
					al.Update(alert)
				}
			}
		}
	}
	for cmd := range groups[0].IterCommands() {
		check, err := r.Server.GetDatabase().GetCheck(stmt, cl.GetID(), cmd.GetID())
		if err != nil {
			r.Server.GetLogger().Error(err, "Internal error")
			r.Output(APIResponse{Error: true, Message: "Internal error"})
			return
		}

		ch, err := NewCheck(check, cmd)
		if err != nil {
			r.Server.GetLogger().Error(err, "Internal error")
			r.Output(APIResponse{Error: true, Message: "Internal error"})
			return
		}
		cl.AddCheck(ch)
	}
	cl.AddGroup(groups[0])
	r.Output(APIResponse{Error: false, Message: "Added the group to the client"})
}

func (GroupHandler) Update(r *HTTPHandler) {
	if r.Request.GroupID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply a group name"})
		return
	}
	group, err := r.Server.GetDatabase().GetGroupByID(r.Request.GroupID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	for cl := range r.Server.GetClients().IterClients() {
		for g := range cl.IterGroups() {
			for c := range g.IterCommands() {
				if c.GetGroupID() == group.ID {
					c.SetNextCheck(group.NextCheck)
					c.SetStopError(group.StopError)
				}
			}
		}
	}
	r.Output(APIResponse{Error: false, Message: "Updated the group in cache"})
}

func (GroupHandler) Delete(r *HTTPHandler) {
	if r.Request.GroupName == "" {
		r.Output(APIResponse{Error: true, Message: "You need to supply a group name."})
		return
	}

	if r.Request.ID != -1 {
		cl := r.Server.GetClients().GetClientByID(r.Request.ID)
		if cl == nil {
			r.Output(APIResponse{Error: true, Message: "Can't find a client with this ID"})
			return
		}

		ok := cl.RemoveGroupsByName(r.Request.GroupName)
		if !ok {
			r.Output(APIResponse{Error: true, Message: "Client does not belong to this group in cache"})
		} else {
			r.Output(APIResponse{Error: false, Message: "Group has now been removed from the client in cache"})
		}
		return
	} else {
		for cl := range r.Server.GetClients().IterClients() {
			cl.RemoveGroupsByName(r.Request.GroupName)
		}
		r.Output(APIResponse{Error: false, Message: "Group has now been removed from all clients in cache"})
	}
}
