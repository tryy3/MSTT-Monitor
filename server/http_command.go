package server

type CommandHandler struct{}

func (CommandHandler) Insert(r *HTTPHandler) {
	// Check if Command ID and Group Name is set.
	if r.Request.CommandID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply a Command ID"})
		return
	}

	if r.Request.GroupName == "" {
		r.Output(APIResponse{Error: true, Message: "You need to supply a group name"})
		return
	}

	// Get the newly added command from the database
	cmd, err := r.Server.GetDatabase().GetCommand(r.Request.CommandID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	// Get the group settings for the command from the database
	group, err := r.Server.GetDatabase().GetGroupByCommand(r.Request.GroupName, r.Request.CommandID)
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}

	// Create a new command in cache
	command := &Command{
		Command:   cmd.Command,
		ID:        cmd.ID,
		GroupID:   group.ID,
		NextCheck: group.NextCheck,
		StopError: group.StopError,
	}

	// Prepare a mysql statement to get the latest check data for this command
	stmt, err := r.Server.GetDatabase().Prepare("SELECT * FROM `checks` WHERE `command_id`=? AND `client_id` ORDER BY `timestamp` DESC LIMIT 1")
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}
	defer stmt.Close()

	// Loop through all clients to find which clients to add the command to
	for c := range r.Server.GetHandler().IterClients() {
		// Add the command to groups if needed
		getCheck := false
		for g := range c.IterGroups() {
			if g.GetName() == r.Request.GroupName {
				g.AddCommand(command)
				getCheck = true
			}
		}
		// If command has been added, then initilize a new check for that client
		if getCheck {
			// Get the latest data from database
			check, err := r.Server.GetDatabase().GetCheck(stmt, command.GetID(), c.GetID())
			if err != nil {
				r.Server.GetLogger().Error(err, "Internal error")
				r.Output(APIResponse{Error: true, Message: "Internal error"})
				return
			}

			// Create a new check and add it to the client
			ch, err := NewCheck(check, command)
			if err != nil {
				r.Server.GetLogger().Error(err, "Internal error")
				r.Output(APIResponse{Error: true, Message: "Internal error"})
				return
			}
			c.AddCheck(ch)
		}
	}

	// Output that the command has been successfully created
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

	for c := range r.Server.GetHandler().IterClients() {
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

	for c := range r.Server.GetHandler().IterClients() {
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
