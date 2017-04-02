package models

import (
	"sync"
	"time"
)

type Check struct {
	rw            *sync.RWMutex
	Command       *Command  `db:"-"`
	PastID        int       `db:"id, primarykey, autoincrement"`
	NextTimestamp time.Time `db:"-"`
	Checked       bool      `db:"checked"`
	Err           bool      `db:"err"`
	Finished      bool      `db:"finished"`
}

// GetCommand
func (c *Check) GetCommand() (cmd *Command) {
	c.rw.RLock()
	cmd = c.Command
	c.rw.RUnlock()
	return
}

// GetID
func (c *Check) GetID() (id int) {
	c.rw.RLock()
	id = c.PastID
	c.rw.RUnlock()
	return
}

// GetTimestamp
func (c *Check) GetTimestamp() (timestamp time.Time) {
	c.rw.RLock()
	timestamp = c.NextTimestamp
	c.rw.RUnlock()
	return
}

// GetChecked
func (c *Check) GetChecked() (checked bool) {
	c.rw.RLock()
	checked = c.Checked
	c.rw.RUnlock()
	return
}

// GetError
func (c *Check) GetError() (err bool) {
	c.rw.RLock()
	err = c.Err
	c.rw.RUnlock()
	return
}

// GetFinished
func (c *Check) GetFinished() (finished bool) {
	c.rw.RLock()
	finished = c.Finished
	c.rw.RUnlock()
	return
}
