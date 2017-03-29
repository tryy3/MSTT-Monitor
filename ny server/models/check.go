package models

import (
	"sync"
	"time"
)

type Check struct {
	*sync.Mutex
	command       *Command
	pastID        int
	nextTimestamp time.Time
	checked       bool
	err           bool
	finished      bool
}

// GetCommand
func (c *Check) GetCommand() (cmd *Command) {
	c.Lock()
	cmd = c.command
	c.Unlock()
	return
}

// GetID
func (c *Check) GetID() (id int) {
	c.Lock()
	id = c.pastID
	c.Unlock()
	return
}

// GetTimestamp
func (c *Check) GetTimestamp() (timestamp time.Time) {
	c.Lock()
	timestamp = c.nextTimestamp
	c.Unlock()
	return
}

// GetChecked
func (c *Check) GetChecked() (checked bool) {
	c.Lock()
	checked = c.checked
	c.Unlock()
	return
}

// GetError
func (c *Check) GetError() (err bool) {
	c.Lock()
	err = c.err
	c.Unlock()
	return
}

// GetFinished
func (c *Check) GetFinished() (finished bool) {
	c.Lock()
	finished = c.finished
	c.Unlock()
	return
}
