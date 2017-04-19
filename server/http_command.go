package server

type CommandHandler struct{}

func (CommandHandler) Insert(r *HTTPHandler) {
	if r.Request.CommandID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply a Command ID"})
		return
	}

	if r.Request.GroupName == "" {
		r.Output(APIResponse{Error: true, Message: "You need to supply a group name"})
		return
	}

	cmd, err := r.Server.GetDatabase().GetCommand(r.Request.CommandID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	group, err := r.Server.GetDatabase().GetGroupByCommand(r.Request.GroupName, r.Request.CommandID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	command := &Command{
		Command:   cmd.Command,
		ID:        cmd.ID,
		GroupID:   group.ID,
		NextCheck: group.NextCheck,
		StopError: group.StopError,
	}

	stmt, err := r.Server.GetDatabase().Prepare("SELECT * FROM `check` WHERE `command_id`=? AND `client_id` ORDER BY `timestamp` DESC")
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}
	defer stmt.Close()

	for c := range r.Server.GetClients().IterClients() {

		getCheck := false
		for g := range c.IterGroups() {
			if g.GetName() == r.Request.GroupName {
				g.AddCommand(command)
				getCheck = true
			}
		}
		if getCheck {
			check, err := r.Server.GetDatabase().GetCheck(stmt, command.GetID(), c.GetID())
			if err != nil {
				r.Server.GetLogger().Error(err, "Internal error")
				r.Output(APIResponse{Error: true, Message: "Internal error"})
				return
			}
			ch, err := NewCheck(check, command)
			if err != nil {
				r.Server.GetLogger().Error(err, "Internal error")
				r.Output(APIResponse{Error: true, Message: "Internal error"})
				return
			}
			c.AddCheck(ch)
		}
	}
	r.Output(APIResponse{Error: false, Message: "Added the new command to the group in cache"})
}

func (CommandHandler) Update(r *HTTPHandler) {
	if r.Request.CommandID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply a Command ID"})
		return
	}

	command, err := r.Server.GetDatabase().GetCommand(r.Request.CommandID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	for c := range r.Server.GetClients().IterClients() {
		for g := range c.IterGroups() {
			for cmd := range g.IterCommands() {
				if cmd.GetID() == command.ID {
					cmd.SetCommand(command.Command)
					break
				}
			}
		}
	}

	r.Output(APIResponse{Error: false, Message: "Successfully updated the command for all clients"})
}

func (CommandHandler) Delete(r *HTTPHandler) {
	if r.Request.CommandID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply a Command ID"})
		return
	}

	for c := range r.Server.GetClients().IterClients() {
		for g := range c.IterGroups() {
			if r.Request.GroupName != "" {
				if r.Request.GroupName == g.GetName() {
					g.RemoveCommandByID(r.Request.CommandID)
				}
			} else {
				g.RemoveCommandByID(r.Request.CommandID)
			}
		}
	}
	r.Output(APIResponse{Error: false, Message: "Successfully removed the command."})
}
