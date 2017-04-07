package server

import (
	"sync"
	"time"
)

func NewCheck(ch checkFields, command *Command) (*Check, error) {
	check := &Check{
		rw:       new(sync.RWMutex),
		Command:  command,
		PastID:   ch.ID,
		Checked:  ch.Checked,
		Err:      ch.Error,
		Finished: ch.Finished,
	}
	err := check.SetTimestampFromString(ch.Timestamp)
	return check, err
}

type Check struct {
	rw            *sync.RWMutex
	Command       *Command
	Group         *Group
	PastID        int64
	NextTimestamp time.Time
	Checked       bool
	Err           bool
	Finished      bool
}

// GetCommand
func (c *Check) GetCommand() (cmd *Command) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	cmd = c.Command
	return
}

// GetGroup
func (c *Check) GetGroup() (g *Group) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	g = c.Group
	return
}

// GetID
func (c *Check) GetID() (id int64) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	id = c.PastID
	return
}

// GetTimestamp
func (c *Check) GetTimestamp() (timestamp time.Time) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	timestamp = c.NextTimestamp
	return
}

// GetChecked
func (c *Check) GetChecked() (checked bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	checked = c.Checked
	return
}

// GetError
func (c *Check) GetError() (err bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	err = c.Err
	return
}

// GetFinished
func (c *Check) GetFinished() (finished bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	finished = c.Finished
	return
}

func (c *Check) SetGroup(g *Group) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Group = g
}

func (c *Check) SetID(id int64) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.PastID = id
}

func (c *Check) SetTimestamp(t time.Time) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.NextTimestamp = t
}

func (c *Check) SetTimestampFromString(t string) error {
	c.rw.Lock()
	defer c.rw.Unlock()
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	if err != nil {
		return err
	}
	timestamp = timestamp.Add(time.Duration(c.Command.GetNextCheck()) * time.Second)
	c.NextTimestamp = timestamp
	return nil
}

func (c *Check) SetChecked(checked bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Checked = checked
}

func (c *Check) SetError(err bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Err = err
}

func (c *Check) SetFinished(finished bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Finished = finished
}
