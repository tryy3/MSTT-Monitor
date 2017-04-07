package server

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"strconv"

	"github.com/bobziuchkovski/cue"
)

func NewClient(cl *clientFields) *Client {
	return &Client{
		rw: new(sync.RWMutex),
		IP: cl.IP,
		ID: cl.ID,
	}
}

type Client struct {
	rw         *sync.RWMutex
	IP         string
	ID         int64
	Groups     []*Group
	Checks     []*Check
	groupNames []string
}

func (c Client) GetIP() (ip string) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.IP
}

func (c Client) GetID() (id int64) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.ID
}

func (c Client) GetGroup(i int) (group *Group) {
	if i >= c.CountGroups() {
		return
	}
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Groups[i]
}

func (c Client) GetGroups() (groups []*Group) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Groups
}

func (c Client) HasGroupByName(name string) bool {
	for g := range c.IterGroups() {
		if g.GetName() == name {
			return true
		}
	}
	return false
}

func (c Client) GetGroupsByName(name string) (groups []*Group) {
	groups = []*Group{}
	for g := range c.IterGroups() {
		if g.GetName() == name {
			groups = append(groups, g)
		}
	}
	return
}

func (c Client) GetGroupsByCommand(commandID int64) (groups []*Group) {
	groups = []*Group{}
	for g := range c.IterGroups() {
		for cmd := range g.IterCommands() {
			if cmd.GetID() == commandID {
				groups = append(groups, g)
				break
			}
		}
	}
	return
}

func (c Client) GetCheck(i int) (check *Check) {
	if i >= c.CountChecks() {
		return
	}
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Checks[i]
}

func (c Client) GetChecks() (checks []*Check) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.Checks
}

func (c Client) GetCheckByPastID(id int64) (check *Check) {
	for che := range c.IterChecks() {
		if che.GetID() == id {
			return che
		}
	}
	return nil
}

func (c Client) GetChecksByCommandID(id int64) []*Check {
	checks := []*Check{}
	for ch := range c.IterChecks() {
		if ch.GetID() == id {
			checks = append(checks, ch)
		}
	}
	return checks
}

func (c Client) SetIP(ip string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.IP = ip
}

func (c Client) IterGroups() <-chan *Group {
	ch := make(chan *Group, c.CountGroups())
	go func() {
		c.rw.RLock()
		defer c.rw.RUnlock()
		for _, group := range c.Groups {
			ch <- group
		}
		close(ch)
	}()
	return ch
}

func (c Client) IterChecks() <-chan *Check {
	ch := make(chan *Check, c.CountChecks())
	go func() {
		c.rw.RLock()
		defer c.rw.RUnlock()
		for _, check := range c.Checks {
			ch <- check
		}
		close(ch)
	}()
	return ch
}

func (c Client) CountGroups() (count int) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.Groups)
}

func (c Client) CountChecks() (count int) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return len(c.Checks)
}

func (c *Client) AddGroup(group *Group) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.groupNames = append(c.groupNames, group.GetName())
	c.Groups = append(c.Groups, group)
}

func (c *Client) AddCheck(check *Check) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Checks = append(c.Checks, check)
}

func (c *Client) RemoveGroup(i int) (ok bool) {
	if i >= c.CountGroups() || i < 0 {
		return false
	}
	c.rw.Lock()
	group := c.Groups[i]
	c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
	c.rw.Unlock()
	for cmd := range group.IterCommands() {
		c.RemoveCheckByGroupID(cmd.GetGroupID())
	}
	return
}

func (c *Client) RemoveGroupsByName(name string) (ok bool) {
	checks := []int64{}
	c.rw.Lock()
	for i := len(c.Groups) - 1; i >= 0; i-- {
		g := c.Groups[i]
		if g != nil && g.GetName() == name {
			c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
			for cmd := range g.IterCommands() {
				checks = append(checks, cmd.GetGroupID())
			}
		}
	}
	c.rw.Unlock()
	if len(checks) > 0 {
		for _, i := range checks {
			c.RemoveCheckByGroupID(i)
		}
		return true
	}
	return false
}

func (c *Client) RemoveGroupsByCommand(commandID int64) (ok bool) {
	var check int64 = -1
	c.rw.Lock()
	for i := len(c.Groups) - 1; i >= 0; i-- {
		g := c.Groups[i]
		for cmd := range g.IterCommands() {
			if cmd.GetID() == commandID {
				c.Groups = append(c.Groups[:i], c.Groups[i+1:]...)
				check = cmd.GetGroupID()
				break
			}
		}
	}
	c.rw.Unlock()
	if check != -1 {
		c.RemoveCheckByGroupID(check)
		return true
	}
	return false
}

func (c *Client) RemoveCheck(i int) (ok bool) {
	if i >= c.CountChecks() || i < 0 {
		return false
	}
	c.rw.Lock()
	defer c.rw.Unlock()
	c.Checks = append(c.Checks[:i], c.Checks[i+1:]...)
	return true
}

func (c *Client) RemoveCheckByID(id int64) (ok bool) {
	c.rw.Lock()
	defer c.rw.Unlock()
	for i := len(c.Checks) - 1; i >= 0; i-- {
		ch := c.Checks[i]
		if ch.GetID() == id {
			c.Checks = append(c.Checks[:i], c.Checks[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Client) RemoveCheckByGroupID(id int64) (ok bool) {
	ok = false
	c.rw.Lock()
	defer c.rw.Unlock()
	for i := len(c.Checks) - 1; i >= 0; i-- {
		ch := c.Checks[i]
		if ch.GetCommand().GetGroupID() == id {
			c.Checks = append(c.Checks[:i], c.Checks[i+1:]...)
			ok = true
		}
	}
	return ok
}

func (c Client) SendMessage(message string) (string, error) {
	conn, err := net.Dial("tcp", c.GetIP()+":3333")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(message)); err != nil {
		return "", err
	}

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		return "", err
	}

	return string(reply[:n]), nil
}

var re = regexp.MustCompile("-port=\"?([\\d,-]+)\"?")

func (c *Client) SendCheck(s *Server, check *Check) string {
	check.SetChecked(true)
	command := check.GetCommand()
	if check.GetID() != -1 {
		err := s.GetDatabase().UpdatePastCheck(check.GetID())
		if err != nil {
			s.GetLogger().Error(err, "Error updating last check")
			return ""
		}
	}

	var resp string
	var err error
	if strings.HasPrefix(command.GetCommand(), "ping") {
		ports := "3333"
		p := re.FindStringSubmatch(command.GetCommand())
		if len(p) >= 2 {
			ports = p[1]
		}

		r, err := c.Ping(ports)
		b, _ := json.Marshal(r)
		e := ""
		if err != nil {
			e = err.Error()
		}
		resp = fmt.Sprintf(`{"error":"%s","ports":%s}`, e, b)
	} else {
		resp, err = c.SendMessage(command.GetCommand())
	}

	if err != nil || !strings.Contains(resp, `"error":""`) {
		check.SetError(true)
		if err != nil {
			resp = err.Error()
		}
	} else {
		check.SetError(false)
	}

	c.SaveCheck(s, check, resp)
	return resp
}

func (c *Client) SaveCheck(s *Server, check *Check, resp string) {
	defer check.SetChecked(false)

	command := check.GetCommand()

	r, err := s.database.InsertCheck(command.GetID(), c.GetID(), resp, check.GetError(), true)
	if err != nil {
		s.GetLogger().Error(err, "Error inserting new check")
		check.SetError(true)
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		s.GetLogger().Error(err, "Error getting last ID")
		check.SetID(-1)
	} else {
		check.SetID(id)
	}

	time, err := s.database.GetLastCheckTime(check.GetID())
	if err != nil {
		s.GetLogger().Error(err, "Error getting last timestamp")
		check.SetError(true)
		return
	}
	check.SetTimestampFromString(time)

	if !check.GetError() && command.GetStopError() {
		c.ResetCheck(check.GetGroup().GetName())
	}
}

func (c *Client) ResetCheck(name string) {
	for ch := range c.IterChecks() {
		if ch.GetGroup().GetName() == name {
			ch.SetError(false)
			ch.SetChecked(false)
		}
	}
}

func (c *Client) Check(s *Server) {
	for check := range c.IterChecks() {
		if check.GetChecked() {
			continue
		}

		if check.GetError() && !check.GetCommand().GetStopError() {
			continue
		}

		t := check.GetTimestamp()
		if t.IsZero() || time.Now().Before(t) {
			continue
		}

		s.GetLogger().WithFields(cue.Fields{
			"CommandID": check.GetCommand().GetID(),
			"ClientID":  c.GetID(),
		}).Info("Starting a check for client")
		go c.SendCheck(s, check)
	}
}

const (
	minTCPPort = 0
	maxTCPPort = 65535
)

type PingResult struct {
	Port   uint16
	Result bool
}

func (c Client) Ping(port string) ([]PingResult, error) {
	pings := []PingResult{}
	ports := [][]uint16{}

	p := strings.Split(strings.Replace(port, " ", "", -1), ",")
	for _, port := range p {
		pSplit := strings.Split(port, "-")
		minPort, err := strconv.ParseUint(pSplit[0], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("The value: %s can't be converted to an integer.", pSplit[0])
		}
		if minPort < minTCPPort || minPort > maxTCPPort {
			return nil, fmt.Errorf("The value: %s is smaller or larger then the maximum TCP port range.", pSplit[0])
		}
		maxPort := minPort

		if len(pSplit) == 2 {
			m, err := strconv.ParseUint(pSplit[1], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("The value: %s can't be converted to an integer.", pSplit[1])
			}
			if m < minTCPPort || m > maxTCPPort {
				return nil, fmt.Errorf("The value: %s is smaller or larger then the maximum TCP port range.", pSplit[1])
			}
			maxPort = m
		}

		if minPort > maxPort {
			return nil, fmt.Errorf("Min value \"%d\" is larger then the max value \"%d\"", minPort, maxPort)
		}
		ports = append(ports, []uint16{uint16(minPort), uint16(maxPort)})
	}

	var err error = nil
	for _, port := range ports {
		for i := port[0]; i <= port[1]; i++ {
			conn, e := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.GetIP(), i), time.Duration(1)*time.Second)
			if e != nil {
				pings = append(pings, PingResult{Port: i, Result: false})
				err = fmt.Errorf("One or more servers failed")
			} else {
				conn.Close()
				pings = append(pings, PingResult{Port: i, Result: true})
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	return pings, err
}
