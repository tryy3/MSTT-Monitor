package server

import (
	"strings"
	"sync"
)

type Handler struct {
	rw      *sync.RWMutex
	Clients []*Client
	Groups  []*Group
}

func (h Handler) GetClient(i int) (cl *Client) {
	if i >= h.CountClients() {
		return
	}
	h.rw.RLock()
	defer h.rw.RUnlock()
	return h.Clients[i]
}

func (h Handler) GetClients() (cls []*Client) {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return h.Clients
}

func (h Handler) GetClientByID(id int64) (cl *Client) {
	for cli := range h.IterClients() {
		if cli.GetID() == id {
			cl = cli
			break
		}
	}
	return
}

func (h Handler) GetGroupByName(name string) *Group {
	for g := range h.IterGroups() {
		if g.GetName() == name {
			return g
		}
	}
	return nil
}

func (h *Handler) AddClient(client *Client) {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.Clients = append(h.Clients, client)
}

func (h *Handler) AddGroup(group *Group) {
	h.rw.Lock()
	defer h.rw.Unlock()
	h.Groups = append(h.Groups, group)
}

func (h *Handler) AddGroupCheckName(group *Group) {
	if h.GetGroupByName(group.GetName()) != nil {
		return
	}
	h.AddGroup(group)
}

func (h *Handler) RemoveClient(i int) (ok bool) {
	if i >= h.CountClients() || i < 0 {
		return false
	}
	h.rw.Lock()
	defer h.rw.Unlock()
	h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
	return true
}

func (h *Handler) RemoveClientByID(id int64) (ok bool) {
	h.rw.Lock()
	defer h.rw.Unlock()
	for i := h.CountClients() - 1; i >= 0; i-- {
		cl := h.GetClient(i)
		if cl != nil && cl.GetID() == id {
			h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
			return true
		}
	}
	return false
}

func (h *Handler) RemoveGroup(i int) (ok bool) {
	if i >= h.CountGroups() || i < 0 {
		return false
	}
	h.rw.Lock()
	defer h.rw.Unlock()
	h.Groups = append(h.Groups[:i], h.Groups[i+1:]...)
	return true
}

func (h Handler) IterClients() <-chan *Client {
	ch := make(chan *Client, h.CountClients())
	go func() {
		h.rw.RLock()
		defer h.rw.RUnlock()
		for _, cl := range h.Clients {
			ch <- cl
		}
		close(ch)
	}()
	return ch
}

func (h Handler) IterGroups() <-chan *Group {
	ch := make(chan *Group, h.CountGroups())
	go func() {
		h.rw.RLock()
		defer h.rw.RUnlock()
		for _, cl := range h.Groups {
			ch <- cl
		}
		close(ch)
	}()
	return ch
}

func (h Handler) CountClients() (count int) {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return len(h.Clients)
}

func (h Handler) CountGroups() (count int) {
	h.rw.RLock()
	defer h.rw.RUnlock()
	return len(h.Groups)
}

func NewHandler(db *Database) (*Handler, error) {
	handler := &Handler{rw: new(sync.RWMutex)}
	cls, err := db.GetClients()
	if err != nil {
		return handler, err
	}
	groups, err := db.GetGroups()
	if err != nil {
		return handler, err
	}
	cmds, err := db.GetCommands()
	if err != nil {
		return handler, err
	}
	alerts, err := db.GetAlertOptions()
	if err != nil {
		return handler, err
	}

	c := db.GetRealCommands(cmds)
	g := db.GetRealGroups(groups, c)

	for _, group := range g {
		handler.AddGroupCheckName(group)
	}

	checkStmt, err := db.Prepare("SELECT * FROM `checks` WHERE `client_id`=? AND `command_id`=? ORDER BY `timestamp` DESC")
	if err != nil {
		return handler, err
	}
	defer checkStmt.Close()

	alertStmt, err := db.Prepare("SELECT * FROM `alerts` WHERE `alert_id`=? AND `client_id`=? ORDER BY `timestamp` DESC")
	if err != nil {
		return handler, err
	}
	defer alertStmt.Close()

	for _, client := range cls {
		cl := NewClient(client)

		for _, group := range strings.Split(client.GroupNames, ",") {
			for _, gg := range g {
				if gg.GetName() == group {
					cl.AddGroup(gg)
					for cmd := range gg.IterCommands() {
						ch, err := NewCheckFromDB(db, checkStmt, cl.GetID(), cmd)
						if err != nil {
							return handler, err
						}
						ch.SetGroup(gg)
						cl.AddCheck(ch)

						for _, a := range alerts {
							if a.ClientID == cl.GetID() && a.CommandID == cmd.GetID() {
								al := NewAlert(a)
								alert, err := db.GetAlert(alertStmt, a.ID, cl.GetID())
								if err == nil {
									err = al.SetTimestampFromString(alert.Timestamp)
									if err != nil {
										return handler, err
									}
								}
								ch.AddAlert(al)
							}
						}
					}
				}
			}
		}
		handler.AddClient(cl)
	}
	return handler, err
}
