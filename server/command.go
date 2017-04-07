package server

import (
	"sync"
)

// NewCommand skapar ett nytt kommand
func NewCommand(command string, commandID int64, nextCheck int, stopError bool) *Command {
	return &Command{
		rw:        new(sync.RWMutex),
		Command:   command,
		ID:        commandID,
		NextCheck: nextCheck,
		StopError: stopError,
	}
}

// Command är en struktur för ett kommand
type Command struct {
	rw        *sync.RWMutex
	Command   string // Kommandot som kommer att skickas till en klient
	ID        int64  // ID från databasen
	GroupID   int64  // ID från gruppen
	NextCheck int    // Interval mellan checks
	StopError bool   // Om checken ska sluta när den stötter på ett error eller inte
}

// GetCommand hämtar kommandot på ett säkert sätt
func (c Command) GetCommand() (cmd string) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Command
}

// GetID Hämtar kommandots ID på ett säkert sätt
func (c Command) GetID() (id int64) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.ID
}

// GetGroupID hämtar kommandots Grupp ID på ett säkert sätt
func (c Command) GetGroupID() int64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.GroupID
}

// GetNextCheck hämtar nextCheck värdet på ett säkert sätt
func (c Command) GetNextCheck() (next int) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.NextCheck
}

// GetStopError hämtar stopError värdet på ett säkert sätt
func (c Command) GetStopError() (stop bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.StopError
}

// SetGroupID sätter kommandots grupp id til lett nytt värde på ett säkert sätt
func (c *Command) SetGroupID(id int64) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.GroupID = id
}

// SetCommand sätter kommandot till ett nytt värde på ett säkert sätt
func (c *Command) SetCommand(command string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Command = command
}

// SetNextCheck sätter ett nytt värde för nextCheck på ett säkert sätt
func (c *Command) SetNextCheck(next int) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.NextCheck = next
}

// SetStopError sätter ett nytt värde för stopError på ett säkert sätt
func (c *Command) SetStopError(stop bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.StopError = stop
}

func (c *Command) Clone() *Command {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return &Command{rw: new(sync.RWMutex), ID: c.ID, Command: c.Command, NextCheck: c.NextCheck, StopError: c.StopError}
}
