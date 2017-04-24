package server

import (
	"database/sql"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewCheckNoFields(command *Command) *Check {
	check := &Check{
		rw:       new(sync.RWMutex),
		Command:  command,
		PastID:   -1,
		Checked:  false,
		Err:      false,
		Finished: true,
	}
	return check
}

func NewCheck(ch checkFields, command *Command) (*Check, error) {
	check := NewCheckNoFields(command)
	check.PastID = ch.ID
	check.Checked = ch.Checked
	check.Err = ch.Error
	check.Finished = ch.Finished
	err := check.SetTimestampFromString(ch.Timestamp)
	return check, err
}

func NewCheckFromDB(db *Database, stmt *sqlx.Stmt, clientID int64, cmd *Command) (*Check, error) {
	var ch *Check
	check, err := db.GetCheck(stmt, clientID, cmd.GetID())
	if err != nil {
		if err == sql.ErrNoRows {
			ch = NewCheckNoFields(cmd)
			return ch, nil
		} else {
			return nil, err
		}
	}

	ch, err = NewCheck(check, cmd)
	return ch, err
}

type Check struct {
	rw            *sync.RWMutex
	Command       *Command
	Group         *Group
	Alerts        []*Alert
	PastID        int64
	NextTimestamp time.Time
	Checked       bool
	Err           bool
	Finished      bool
}

// GetCommand
func (c *Check) GetCommand() *Command {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Command
}

// GetGroup
func (c *Check) GetGroup() *Group {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Group
}

// GetAlerts
func (c *Check) GetAlerts() []*Alert {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Alerts
}

// GetID
func (c *Check) GetID() int64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.PastID
}

// GetTimestamp
func (c *Check) GetTimestamp() time.Time {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.NextTimestamp
}

// GetChecked
func (c *Check) GetChecked() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Checked
}

// GetError
func (c *Check) GetError() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Err
}

// GetFinished
func (c *Check) GetFinished() bool {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Finished
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

func (c *Check) AddAlert(alert *Alert) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Alerts = append(c.Alerts, alert)
}

func (c *Check) RemoveAlertByID(id int64) {
	c.rw.Lock()
	defer c.rw.Unlock()
	for i := len(c.Alerts) - 1; i >= 0; i-- {
		a := c.Alerts[i]
		if a.GetID() == id {
			c.Alerts = append(c.Alerts[:i], c.Alerts[i+1:]...)
			return
		}
	}
}

func (c Check) CountAlerts() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.Alerts)
}

func (c Check) IterAlerts() <-chan *Alert {
	ch := make(chan *Alert, c.CountAlerts())
	go func() {
		c.rw.RLock()
		defer c.rw.RUnlock()
		for _, alert := range c.Alerts {
			ch <- alert
		}
		close(ch)
	}()
	return ch
}
