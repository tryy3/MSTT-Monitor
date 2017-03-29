package models

import (
	"sync"
)

// NewCommand skapar ett nytt kommand
func NewCommand(command string, commandID int, nextCheck int, stopError bool) *Command {
	return &Command{
		command:   command,
		commandID: commandID,
		nextCheck: nextCheck,
		stopError: stopError,
	}
}

// Command är en struktur för ett kommand
type Command struct {
	*sync.Mutex
	command   string // Kommandot som kommer att skickas till en klient
	commandID int    // ID från databasen
	nextCheck int    // Interval mellan checks
	stopError bool   // Om checken ska sluta när den stötter på ett error eller inte
}

// GetCommand hämtar kommandot på ett säkert sätt
func (c Command) GetCommand() (cmd string) {
	c.Lock()
	cmd = c.command
	c.Unlock()
	return
}

// GetID Hämtar kommandots ID på ett säkert sätt
func (c Command) GetID() (id int) {
	c.Lock()
	id = c.commandID
	c.Unlock()
	return
}

// GetNextCheck hämtar nextCheck värdet på ett säkert sätt
func (c Command) GetNextCheck() (next int) {
	c.Lock()
	next = c.nextCheck
	c.Unlock()
	return
}

// GetStopError hämtar stopError värdet på ett säkert sätt
func (c Command) GetStopError() (stop bool) {
	c.Lock()
	stop = c.stopError
	c.Unlock()
	return
}

// SetCommand sätter kommandot till ett nytt värde på ett säkert sätt
func (c *Command) SetCommand(command string) {
	c.Lock()
	c.command = command
	c.Unlock()
	return
}

// SetNextCheck sätter ett nytt värde för nextCheck på ett säkert sätt
func (c *Command) SetNextCheck(next int) {
	c.Lock()
	c.nextCheck = next
	c.Unlock()
	return
}

// SetStopError sätter ett nytt värde för stopError på ett säkert sätt
func (c *Command) SetStopError(stop bool) {
	c.Lock()
	c.stopError = stop
	c.Unlock()
	return
}
