package server

import (
	"sync"
)

// NewGroup skapar en ny Grupp
func NewGroup(name string) *Group {
	return &Group{
		rw:       new(sync.RWMutex),
		Name:     name,
		Commands: []*Command{},
	}
}

// Command är en struktur för en grupp
type Group struct {
	rw       *sync.RWMutex
	Name     string
	Commands []*Command
}

// GetName hämtar namnet på gruppen på ett säkert sätt
func (g Group) GetName() (name string) {
	g.rw.RLock()
	defer g.rw.RUnlock()
	return g.Name
}

// GetCommands hämtar alla commands i gruppen
func (g Group) GetCommands() (cmds []*Command) {
	g.rw.RLock()
	defer g.rw.RUnlock()
	return g.Commands
}

// GetCommand hämtar en specifik command i gruppen
func (g Group) GetCommand(i int) (cmd *Command) {
	if i >= g.Count() {
		return
	}
	g.rw.RLock()
	defer g.rw.RUnlock()
	return g.Commands[i]
}

// GetCommandByID hämtar en specifik command i gruppen med commandID
func (g Group) GetCommandByID(id int64) (cmd *Command) {
	for c := range g.IterCommands() {
		if c.GetID() == id {
			cmd = c
			break
		}
	}
	return
}

func (g Group) HasCommand(id int64) (ok bool) {
	c := g.GetCommandByID(id)
	if c != nil {
		return true
	}
	return false
}

// Length hämtar antal kommandon som tillhör denna grupp
func (g Group) Count() (count int) {
	g.rw.RLock()
	defer g.rw.RUnlock()
	return len(g.Commands)
}

// SetName sätter ett nytt värde för namnet på gruppen
func (g *Group) SetName(name string) {
	g.rw.Lock()
	defer g.rw.Unlock()
	g.Name = name
}

// Add lägger till ett nytt kommand för gruppen
func (g *Group) AddCommand(cmd *Command) {
	g.rw.Lock()
	defer g.rw.Unlock()
	g.Commands = append(g.Commands, cmd)
}

// Remove tar bort ett specifikt kommand från gruppen
func (g *Group) RemoveCommand(i int) (ok bool) {
	if i >= g.Count() || i < 0 {
		return false
	}
	g.rw.Lock()
	defer g.rw.Unlock()
	g.Commands = append(g.Commands[:i], g.Commands[i+1:]...)
	return true
}

// RemoveByID tar bort ett specifikt kommand med hjälp av commandID från gruppen
func (g *Group) RemoveCommandByID(id int64) (ok bool) {
	g.rw.Lock()
	defer g.rw.Unlock()
	for i := g.Count() - 1; i >= 0; i-- {
		c := g.GetCommand(i)
		if c != nil && c.GetID() == id {
			g.Commands = append(g.Commands[:i], g.Commands[i+1:]...)
			return true
		}
	}
	return false
}

func (g Group) IterCommands() <-chan *Command {
	ch := make(chan *Command, g.Count())
	go func() {
		g.rw.RLock()
		defer g.rw.RUnlock()
		for _, cmd := range g.Commands {
			ch <- cmd
		}
		close(ch)
	}()
	return ch
}
