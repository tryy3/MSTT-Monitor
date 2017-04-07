package server

import (
	"strings"
	"sync"
)

type Clients struct {
	rw      *sync.RWMutex
	Clients []*Client
}

func (c Clients) GetClient(i int) (cl *Client) {
	if i >= c.Count() {
		return
	}
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Clients[i]
}

func (c Clients) GetClients() (cls []*Client) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Clients
}

func (c Clients) GetClientByID(id int64) (cl *Client) {
	for cli := range c.IterClients() {
		if cli.GetID() == id {
			cl = cli
			break
		}
	}
	return
}

func (c *Clients) AddClient(client *Client) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Clients = append(c.Clients, client)
}

func (c *Clients) RemoveClient(i int) (ok bool) {
	if i >= c.Count() || i < 0 {
		return false
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
	return true
}

func (c *Clients) RemoveClientByID(id int64) (ok bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	for i := c.Count() - 1; i >= 0; i-- {
		cl := c.GetClient(i)
		if cl != nil && cl.GetID() == id {
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
			return true
		}
	}
	return false
}

func (c Clients) IterClients() <-chan *Client {
	ch := make(chan *Client, c.Count())
	go func() {
		c.rw.RLock()
		defer c.rw.RUnlock()
		for _, cl := range c.Clients {
			ch <- cl
		}
		close(ch)
	}()
	return ch
}

func (c Clients) Count() (count int) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.Clients)
}

func NewClients(db *Database) (*Clients, error) {
	clients := &Clients{rw: new(sync.RWMutex)}
	cls, err := db.GetClients()
	if err != nil {
		return clients, err
	}
	groups, err := db.GetGroups()
	if err != nil {
		return clients, err
	}
	cmds, err := db.GetCommands()
	if err != nil {
		return clients, err
	}

	c := db.GetRealCommands(cmds)
	g := db.GetRealGroups(groups, c)

	stmt, err := db.Prepare("SELECT * FROM `checks` WHERE `client_id`=? AND `command_id`=? ORDER BY `timestamp` DESC")
	if err != nil {
		return clients, err
	}
	defer stmt.Close()

	for _, client := range cls {
		cl := NewClient(client)

		for _, group := range strings.Split(client.GroupNames, ",") {
			for _, gg := range g {
				if gg.GetName() == group {
					cl.AddGroup(gg)
					for cmd := range gg.IterCommands() {
						check, err := db.GetCheck(stmt, cl.GetID(), cmd.GetID())
						if err != nil {
							return clients, err
						}

						ch, err := NewCheck(check, cmd)
						ch.SetGroup(gg)
						if err != nil {
							return clients, err
						}
						cl.AddCheck(ch)
					}
				}
			}
		}
		clients.AddClient(cl)
	}
	return clients, err
}
