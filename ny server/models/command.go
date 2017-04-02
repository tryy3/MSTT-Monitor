package models

import (
	"sync"
)

// NewCommand skapar ett nytt kommand
func NewCommand(command string, commandID int, nextCheck int, stopError bool) *Command {
	return &Command{
		Command:   command,
		CommandID: commandID,
		NextCheck: nextCheck,
		StopError: stopError,
	}
}

// Command är en struktur för ett kommand
type Command struct {
	rw        *sync.RWMutex
	Command   string `db:"command"`                       // Kommandot som kommer att skickas till en klient
	CommandID int    `db:"id, primarykey, autoincrement"` // ID från databasen
	NextCheck int    `db:"next_check"`                    // Interval mellan checks
	StopError bool   `db:"stop_error"`                    // Om checken ska sluta när den stötter på ett error eller inte
}

// GetCommand hämtar kommandot på ett säkert sätt
func (c Command) GetCommand() (cmd string) {
	c.rw.RLock()
	cmd = c.Command
	c.rw.RUnlock()
	return
}

// GetID Hämtar kommandots ID på ett säkert sätt
func (c Command) GetID() (id int) {
	c.rw.RLock()
	id = c.CommandID
	c.rw.RUnlock()
	return
}

// GetNextCheck hämtar nextCheck värdet på ett säkert sätt
func (c Command) GetNextCheck() (next int) {
	c.rw.RLock()
	next = c.NextCheck
	c.rw.RUnlock()
	return
}

// GetStopError hämtar stopError värdet på ett säkert sätt
func (c Command) GetStopError() (stop bool) {
	c.rw.RLock()
	stop = c.StopError
	c.rw.RUnlock()
	return
}

// SetCommand sätter kommandot till ett nytt värde på ett säkert sätt
func (c *Command) SetCommand(command string) {
	c.rw.Lock()
	c.Command = command
	c.rw.Unlock()
	return
}

// SetNextCheck sätter ett nytt värde för nextCheck på ett säkert sätt
func (c *Command) SetNextCheck(next int) {
	c.rw.Lock()
	c.NextCheck = next
	c.rw.Unlock()
	return
}

// SetStopError sätter ett nytt värde för stopError på ett säkert sätt
func (c *Command) SetStopError(stop bool) {
	c.rw.Lock()
	c.StopError = stop
	c.rw.Unlock()
	return
}
