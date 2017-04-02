package models

import (
	"sync"
)

type Client struct {
	rw     *sync.RWMutex
	IP     string   `db:"ip"`
	ID     int      `db:"id, primarykey, autoincrement"`
	Groups []*Group `db:"-"`
	Checks []*Check `db:"-"`
}

func (c Client) GetIP() (ip string) {
	c.rw.RLock()
	ip = c.IP
	c.rw.RUnlock()
	return
}

func (c Client) GetID() (id int) {
	c.rw.RLock()
	id = c.ID
	c.rw.RUnlock()
	return
}

func (c Client) GetGroup(i int) (group *Group) {
	if i >= c.CountGroups() {
		return
	}
	c.rw.RLock()
	group = c.Groups[i]
	c.rw.RUnlock()
	return
}

func (c Client) GetGroups() (groups []*Group) {
	c.rw.RLock()
	groups = c.Groups
	c.rw.RUnlock()
	return
}

func (c Client) GetGroupsByID(id int) (group *Group) {
	for g := range c.IterGroups() {
		if g.GetID() == id {
			group = g
			break
		}
	}
	return
}

func (c Client) GetGroupsByName(name string) (groups []*Group) {
	groups = []*Group{}
	for g := range c.IterGroups() {
		if g.GetName() == name {
			groups = append(groups, g)
		}
	}
	return
}

func (c Client) GetGroupsByCommand(commandID int) (groups []*Group) {
	groups = []*Group{}
	for g := range c.IterGroups() {
		for cmd := range g.Iter() {
			if cmd.GetID() == commandID {
				groups = append(groups, g)
				break
			}
		}
	}
	return
}

func (c Client) GetCheck(i int) (check *Check) {
	if i >= c.CountChecks() {
		return
	}
	c.rw.RLock()
	check = c.Checks[i]
	c.rw.RUnlock()
	return
}

func (c Client) GetChecks() (checks []*Check) {
	c.rw.RLock()
	checks = c.Checks
	c.rw.RUnlock()
	return
}

func (c Client) GetCheckByPastID(id int) (check *Check) {
	for che := range c.IterChecks() {
		if check.GetID() == id {
			check = che
			break
		}
	}
	return
}

// GetChecksReady hämtar alla checks som är redo att kollas
// TODO: Implement this
func (c Client) GetChecksReady() (check []*Check) {
	return nil
}

func (c Client) IterGroups() <-chan *Group {
	ch := make(chan *Group, c.CountGroups())
	go func() {
		c.rw.RLock()
		for _, group := range c.Groups {
			ch <- group
		}
		c.rw.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Client) IterChecks() <-chan *Check {
	ch := make(chan *Check, c.CountChecks())
	go func() {
		c.rw.RLock()
		for _, check := range c.Checks {
			ch <- check
		}
		c.rw.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Client) CountGroups() (count int) {
	c.rw.RLock()
	count = len(c.Groups)
	c.rw.RUnlock()
	return
}

func (c Client) CountChecks() (count int) {
	c.rw.RLock()
	count = len(c.Checks)
	c.rw.RUnlock()
	return
}

func (c *Client) AddGroup(group *Group) {
	c.rw.Lock()
	c.Groups = append(c.Groups, group)
	c.rw.Unlock()
}

func (c *Client) AddCheck(check *Check) {
	c.rw.Lock()
	c.Checks = append(c.Checks, check)
	c.rw.Unlock()
}

func (c *Client) RemoveGroup(i int) (ok bool) {
	ok = false
	if i >= c.CountGroups() || i < 0 {
		return
	}
	c.rw.Lock()
	c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
	ok = true
	c.rw.Unlock()
	return
}

func (c *Client) RemoveGroupByID(id int) (ok bool) {
	ok = false
	c.rw.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil && g.GetID() == id {
			c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
			ok = true
			break
		}
	}
	c.rw.Unlock()
	return
}

func (c *Client) RemoveGroupsByName(name string) (ok bool) {
	ok = false
	c.rw.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil && g.GetName() == name {
			c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
			ok = true
		}
	}
	c.rw.Unlock()
	return
}

func (c *Client) RemoveGroupsByCommand(commandID int) (ok bool) {
	ok = false
	c.rw.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil {
			for cmd := range g.Iter() {
				if cmd.GetID() == commandID {
					c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
					ok = true
					break
				}
			}
		}
	}
	c.rw.Unlock()
	return
}

func (c *Client) RemoveCheck(i int) (ok bool) {
	ok = false
	if i >= c.CountChecks() || i < 0 {
		return
	}
	c.rw.Lock()
	c.Checks = append(c.Checks[:i], c.Checks[i+1:]...)
	ok = true
	c.rw.Unlock()
	return
}

func (c *Client) RemoveCheckByID(id int) (ok bool) {
	ok = false
	c.rw.Lock()
	for i := c.CountChecks() - 1; i >= 0; i-- {
		g := c.GetCheck(i)
		if g != nil && g.GetID() == id {
			c.Checks = append(c.Checks[:i], c.Checks[i+1:]...)
			ok = true
			break
		}
	}
	c.rw.Unlock()
	return
}

//Send
