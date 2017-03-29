package models

import (
	"sync"
)

type Clients struct {
	*sync.RWMutex
	clients []*Client
}

func (c Clients) GetClient(i int) (cl *Client) {
	if i >= c.Count() {
		return
	}
	c.RLock()
	cl = c.clients[i]
	c.RUnlock()
	return
}

func (c Clients) GetClients() (cls []*Client) {
	c.RLock()
	cls = c.clients
	c.RUnlock()
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
	c.Lock()
	c.clients = append(c.clients, client)
	c.Unlock()
}

func (c *Clients) Remove(i int) (ok bool) {
	ok = false
	if i >= c.Count() || i < 0 {
		return
	}
	c.Lock()
	c.clients = append(c.clients[:i], c.clients[i+1:]...)
	ok = true
	c.Unlock()
	return
}

func (c *Clients) RemoveByID(id int) (ok bool) {
	ok = false
	c.Lock()
	for i := c.Count() - 1; i >= 0; i-- {
		cl := c.GetClient(i)
		if cl != nil && cl.GetID() == id {
			c.clients = append(c.clients[:i], c.clients[i+1:]...)
			ok = true
			break
		}
	}
	c.Unlock()
	return
}

func (c Clients) Iter() <-chan *Client {
	ch := make(chan *Client, c.Count())
	go func() {
		c.RLock()
		for _, cl := range c.clients {
			ch <- cl
		}
		c.RUnlock()
		close(ch)
	}()
	return ch
}

func (c Clients) Count() (count int) {
	c.RLock()
	count = len(c.clients)
	c.RUnlock()
	return
}
