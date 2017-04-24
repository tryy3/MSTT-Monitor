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

	alertStmt, err := r.Server.GetDatabase().Prepare("SELECT * FROM `checks` WHERE `alert_id`=? AND `client_id`=? ORDER BY `timestamp` DESC")
	if err != nil {
		r.Server.GetLogger().Error(err, "Internal error")
		r.Output(APIResponse{Error: true, Message: "Internal error"})
		return
	}
	defer alertStmt.Close()

	for cl := range r.Server.GetHandler().IterClients() {
		for ch := range cl.IterChecks() {
			for al := range ch.IterAlerts() {
				if al.GetID() == r.Request.ID {
					r.Output(APIResponse{Error: true, Message: "There is already an alert with this ID in cache."})
					return
				}
			}
		}
	}

	for cl := range r.Server.GetHandler().IterClients() {
		if cl.GetID() == alert.ClientID {
			for ch := range cl.IterChecks() {
				if ch.GetCommand().GetID() == alert.CommandID {
					a := NewAlert(alert)

					al, err := r.Server.GetDatabase().GetAlert(alertStmt, a.GetID(), cl.GetID())
					if err == nil {
						err = a.SetTimestampFromString(al.Timestamp)
						if err != nil {
							r.Server.GetLogger().Error(err, "Internal error")
							r.Output(APIResponse{Error: true, Message: "Internal error"})
							return
						}
					}
					r.Output(APIResponse{Error: false, Message: "Added the alert to the client"})
					return
				}
			}
		}
	}
	r.Output(APIResponse{Error: true, Message: "Can't find the correct client or command to add the alert to"})
}

func (AlertHandler) Update(r *HTTPHandler) {
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

	for cl := range r.Server.GetHandler().IterClients() {
		if cl.GetID() == alert.ClientID {
			for ch := range cl.IterChecks() {
				if ch.GetCommand().GetID() == alert.CommandID {
					for a := range ch.IterAlerts() {
						if a.GetID() == alert.ID {
							a.Update(alert)
						}
					}
					r.Output(APIResponse{Error: false, Message: "Updated the alert for the client"})
					return
				}
			}
		}
	}
	r.Output(APIResponse{Error: true, Message: "Can't find the correct client or command to update the alert"})
}

func (AlertHandler) Delete(r *HTTPHandler) {
	if r.Request.ID == -1 {
		r.Output(APIResponse{Error: true, Message: "You need to supply an alert option ID"})
		return
	}

	for cl := range r.Server.GetHandler().IterClients() {
		for ch := range cl.IterChecks() {
			for alert := range ch.IterAlerts() {
				if alert.GetID() == r.Request.ID {
					ch.RemoveAlertByID(r.Request.ID)
					r.Output(APIResponse{Error: false, Message: "Deleted the alert from the client"})
					return
				}
			}
		}
	}
	r.Output(APIResponse{Error: true, Message: "Can't find the correct client or command to delete the alert from"})
}
