package models

import (
	"sync"
)

type Client struct {
	*sync.RWMutex
	ip     string
	id     int
	groups []*Group
	checks []*Check
}

func (c Client) GetIP() (ip string) {
	c.RLock()
	ip = c.ip
	c.RUnlock()
	return
}

func (c Client) GetID() (id int) {
	c.RLock()
	id = c.id
	c.RUnlock()
	return
}

func (c Client) GetGroup(i int) (group *Group) {
	if i >= c.CountGroups() {
		return
	}
	c.RLock()
	group = c.groups[i]
	c.RUnlock()
	return
}

func (c Client) GetGroups() (groups []*Group) {
	c.RLock()
	groups = c.groups
	c.RUnlock()
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
	c.RLock()
	check = c.checks[i]
	c.RUnlock()
	return
}

func (c Client) GetChecks() (checks []*Check) {
	c.RLock()
	checks = c.checks
	c.RUnlock()
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
		c.RLock()
		for _, group := range c.groups {
			ch <- group
		}
		c.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Client) IterChecks() <-chan *Check {
	ch := make(chan *Check, c.CountChecks())
	go func() {
		c.RLock()
		for _, check := range c.checks {
			ch <- check
		}
		c.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Client) CountGroups() (count int) {
	c.RLock()
	count = len(c.groups)
	c.RUnlock()
	return
}

func (c Client) CountChecks() (count int) {
	c.RLock()
	count = len(c.checks)
	c.RUnlock()
	return
}

func (c *Client) AddGroup(group *Group) {
	c.Lock()
	c.groups = append(c.groups, group)
	c.Unlock()
}

func (c *Client) AddCheck(check *Check) {
	c.Lock()
	c.checks = append(c.checks, check)
	c.Unlock()
}

func (c *Client) RemoveGroup(i int) (ok bool) {
	ok = false
	if i >= c.CountGroups() || i < 0 {
		return
	}
	c.Lock()
	c.groups = append(c.groups[:i], c.groups[i+1:]...)
	ok = true
	c.Unlock()
	return
}

func (c *Client) RemoveGroupByID(id int) (ok bool) {
	ok = false
	c.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil && g.GetID() == id {
			c.groups = append(c.groups[:i], c.groups[i+1:]...)
			ok = true
			break
		}
	}
	c.Unlock()
	return
}

func (c *Client) RemoveGroupsByName(name string) (ok bool) {
	ok = false
	c.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil && g.GetName() == name {
			c.groups = append(c.groups[:i], c.groups[i+1:]...)
			ok = true
		}
	}
	c.Unlock()
	return
}

func (c *Client) RemoveGroupsByCommand(commandID int) (ok bool) {
	ok = false
	c.Lock()
	for i := c.CountGroups() - 1; i >= 0; i-- {
		g := c.GetGroup(i)
		if g != nil {
			for cmd := range g.Iter() {
				if cmd.GetID() == commandID {
					c.groups = append(c.groups[:i], c.groups[i+1:]...)
					ok = true
					break
				}
			}
		}
	}
	c.Unlock()
	return
}

func (c *Client) RemoveCheck(i int) (ok bool) {
	ok = false
	if i >= c.CountChecks() || i < 0 {
		return
	}
	c.Lock()
	c.checks = append(c.checks[:i], c.checks[i+1:]...)
	ok = true
	c.Unlock()
	return
}

func (c *Client) RemoveCheckByID(id int) (ok bool) {
	ok = false
	c.Lock()
	for i := c.CountChecks() - 1; i >= 0; i-- {
		g := c.GetCheck(i)
		if g != nil && g.GetID() == id {
			c.checks = append(c.checks[:i], c.checks[i+1:]...)
			ok = true
			break
		}
	}
	c.Unlock()
	return
}

//Send
