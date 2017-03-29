package models

import (
	"sync"
)

// NewGroup skapar en ny Grupp
func NewGroup(name string, id int) *Group {
	return &Group{
		name:     name,
		id:       id,
		commands: []*Command{},
	}
}

// Command är en struktur för en grupp
type Group struct {
	*sync.RWMutex
	name     string
	id       int
	commands []*Command
}

// GetName hämtar namnet på gruppen på ett säkert sätt
func (g Group) GetName() (name string) {
	g.RLock()
	name = g.name
	g.RUnlock()
	return
}

// GetID hämtar idn på gruppen på ett säkert sätt
func (g Group) GetID() (id int) {
	g.RLock()
	id = g.id
	g.RUnlock()
	return
}

// GetCommands hämtar alla commands i gruppen
func (g Group) GetCommands() (cmds []*Command) {
	g.RLock()
	cmds = g.commands
	g.RUnlock()
	return
}

// GetCommand hämtar en specifik command i gruppen
func (g Group) GetCommand(i int) (cmd *Command) {
	if i >= g.Count() {
		return
	}
	g.RLock()
	cmd = g.commands[i]
	g.RUnlock()
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
	g.RLock()
	count = len(g.commands)
	g.RUnlock()
	return
}

// SetName sätter ett nytt värde för namnet på gruppen
func (g *Group) SetName(name string) {
	g.Lock()
	g.name = name
	g.Unlock()
	return
}

// Add lägger till ett nytt kommand för gruppen
func (g *Group) Add(cmd *Command) {
	g.Lock()
	g.commands = append(g.commands, cmd)
	g.Unlock()
	return
}

// Remove tar bort ett specifikt kommand från gruppen
func (g *Group) Remove(i int) (ok bool) {
	ok = false
	if i >= g.Count() || i < 0 {
		return
	}
	g.Lock()
	g.commands = append(g.commands[:i], g.commands[i+1:]...)
	ok = true
	g.Unlock()
	return
}

// RemoveByID tar bort ett specifikt kommand med hjälp av commandID från gruppen
func (g *Group) RemoveByID(id int) (ok bool) {
	ok = false
	g.Lock()
	for i := g.Count() - 1; i >= 0; i-- {
		c := g.GetCommand(i)
		if c != nil && c.GetID() == id {
			g.commands = append(g.commands[:i], g.commands[i+1:]...)
			ok = true
			break
		}
	}
	g.Unlock()
	return
}

func (g Group) Iter() <-chan *Command {
	ch := make(chan *Command, g.Count())
	go func() {
		g.RLock()
		for _, cmd := range g.commands {
			ch <- cmd
		}
		g.RUnlock()
		close(ch)
	}()
	return ch
}
