package models

import (
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
	cl = c.Clients[i]
	c.rw.RUnlock()
	return
}

func (c Clients) GetClients() (cls []*Client) {
	c.rw.RLock()
	cls = c.Clients
	c.rw.RUnlock()
	return
}

func (c Clients) GetClientByID(id int) (cl *Client) {
	for cli := range c.Iter() {
		if cli.GetID() == id {
			cl = cli
			break
		}
	}
	return
}

func (c *Clients) Add(client *Client) {
	c.rw.Lock()
	c.Clients = append(c.Clients, client)
	c.rw.Unlock()
}

func (c *Clients) Remove(i int) (ok bool) {
	ok = false
	if i >= c.Count() || i < 0 {
		return
	}
	c.rw.Lock()
	c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
	ok = true
	c.rw.Unlock()
	return
}

func (c *Clients) RemoveByID(id int) (ok bool) {
	ok = false
	c.rw.Lock()
	for i := c.Count() - 1; i >= 0; i-- {
		cl := c.GetClient(i)
		if cl != nil && cl.GetID() == id {
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
			ok = true
			break
		}
	}
	c.rw.Unlock()
	return
}

func (c Clients) Iter() <-chan *Client {
	ch := make(chan *Client, c.Count())
	go func() {
		c.rw.RLock()
		for _, cl := range c.Clients {
			ch <- cl
		}
		c.rw.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Clients) Count() (count int) {
	c.rw.RLock()
	count = len(c.Clients)
	c.rw.RUnlock()
	return
}
