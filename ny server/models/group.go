package models

import (
	"sync"
)

// NewGroup skapar en ny Grupp
func NewGroup(name string, id int) *Group {
	return &Group{
		Name:     name,
		ID:       id,
		Commands: []*Command{},
	}
}

// Command är en struktur för en grupp
type Group struct {
	rw       *sync.RWMutex
	Name     string     `db:"name"`
	ID       int        `db:"id, primarykey, autoincrement"`
	Commands []*Command `db:"-"`
}

// GetName hämtar namnet på gruppen på ett säkert sätt
func (g Group) GetName() (name string) {
	g.rw.RLock()
	name = g.Name
	g.rw.RUnlock()
	return
}

// GetID hämtar idn på gruppen på ett säkert sätt
func (g Group) GetID() (id int) {
	g.rw.RLock()
	id = g.ID
	g.rw.RUnlock()
	return
}

// GetCommands hämtar alla commands i gruppen
func (g Group) GetCommands() (cmds []*Command) {
	g.rw.RLock()
	cmds = g.Commands
	g.rw.RUnlock()
	return
}

// GetCommand hämtar en specifik command i gruppen
func (g Group) GetCommand(i int) (cmd *Command) {
	if i >= g.Count() {
		return
	}
	g.rw.RLock()
	cmd = g.Commands[i]
	g.rw.RUnlock()
	return
}

// GetCommandByID hämtar en specifik command i gruppen med commandID
func (g Group) GetCommandByID(id int) (cmd *Command) {
	for c := range g.Iter() {
		if c.GetID() == id {
			cmd = c
			break
		}
	}
	return
}

func (g Group) HasCommand(id int) (ok bool) {
	ok = false
	c := g.GetCommandByID(id)
	if c != nil {
		ok = true
	}
	return
}

// Length hämtar antal kommandon som tillhör denna grupp
func (g Group) Count() (count int) {
	g.rw.RLock()
	count = len(g.Commands)
	g.rw.RUnlock()
	return
}

// SetName sätter ett nytt värde för namnet på gruppen
func (g *Group) SetName(name string) {
	g.rw.Lock()
	g.Name = name
	g.rw.Unlock()
	return
}

// Add lägger till ett nytt kommand för gruppen
func (g *Group) Add(cmd *Command) {
	g.rw.Lock()
	g.Commands = append(g.Commands, cmd)
	g.rw.Unlock()
	return
}

// Remove tar bort ett specifikt kommand från gruppen
func (g *Group) Remove(i int) (ok bool) {
	ok = false
	if i >= g.Count() || i < 0 {
		return
	}
	g.rw.Lock()
	g.Commands = append(g.Commands[:i], g.Commands[i+1:]...)
	ok = true
	g.rw.Unlock()
	return
}

// RemoveByID tar bort ett specifikt kommand med hjälp av commandID från gruppen
func (g *Group) RemoveByID(id int) (ok bool) {
	ok = false
	g.rw.Lock()
	for i := g.Count() - 1; i >= 0; i-- {
		c := g.GetCommand(i)
		if c != nil && c.GetID() == id {
			g.Commands = append(g.Commands[:i], g.Commands[i+1:]...)
			ok = true
			break
		}
	}
	g.rw.Unlock()
	return
}

func (g Group) Iter() <-chan *Command {
	ch := make(chan *Command, g.Count())
	go func() {
		g.rw.RLock()
		for _, cmd := range g.Commands {
			ch <- cmd
		}
		g.rw.RUnlock()
		close(ch)
	}()
	return ch
}
