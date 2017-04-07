package server

import (
	"sync"
	"time"
)

// Clients är en struct för alla klienter, så att man kan modifiera
// klienter väldigt enkelt över flera gorutines.
type Clients struct {
	sync.Mutex           // Mutex så att man kan modifiera utan att det blir race conditions
	clients    []*Client // clients innehåller alla klienter
}

// Get hämtar en client från listan
// och svarar med den rätta klienten, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Clients) Get(id int) *Client {
	if len(c.clients)-1 < id {
		return nil
	}
	c.Lock()
	cl := c.clients[id]
	c.Unlock()
	return cl
}

// GetClientByID letar efter en klient med den korrekta IDn
// och svarar med den rätta klienten, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Clients) GetClientByID(id int) *Client {
	c.Lock()
	for _, cl := range c.clients {
		if cl.clientID == id {
			c.Unlock()
			return cl
		}
	}
	c.Unlock()
	return nil
}

// Length svarar med antal klienter.
func (c *Clients) Length() int {
	return len(c.clients)
}

// RemoveByClientID tar bort en klient och svarar med en boolean
// om klienten hittades eller inte, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Clients) RemoveByClientID(id int) bool {
	c.Lock()
	for i := len(c.clients) - 1; i >= 0; i-- {
		cl := c.clients[i]
		if cl.clientID == id {
			c.clients = append(c.clients[:i], c.clients[i+1:]...)
			c.Unlock()
			return true
		}
	}
	c.Unlock()
	return false
}

// Add lägger till en ny klient, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Clients) Add(cl *Client) {
	c.Lock()
	c.clients = append(c.clients, cl)
	c.Unlock()
}

// Client är en strukt för alla klienter
type Client struct {
	sync.Mutex          // Mutex så att man kan modifiera utan att det blir race conditions
	ip         string   // IPn till klienten
	clientID   int      // IDn för klienten
	groups     []string // Grupper som klienten hör till
	checks     []*Check // Alla checks som hör till klienten.
}

// Get hämtar en check från listan
// och svarar med den rätta klienten, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Client) Get(id int) *Check {
	if len(c.checks)-1 < id {
		return nil
	}
	c.Lock()
	ch := c.checks[id]
	c.Unlock()
	return ch
}

// GetCheckByCommandID letar efter en check med den korrekta commandID
// och svarar med den rätta klienten, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Client) GetCheckByCommandID(commandID int) *Check {
	c.Lock()
	for _, ch := range c.checks {
		if ch.commandID == commandID {
			c.Unlock()
			return ch
		}
	}
	c.Unlock()
	return nil
}

// Length svarar med antal klienter.
func (c *Client) Length() int {
	return len(c.checks)
}

// RemoveByCommandID tar bort en klient och svarar med en boolean
// om klienten hittades eller inte, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Client) RemoveByCommandID(commandID int) bool {
	c.Lock()
	for i := len(c.checks) - 1; i >= 0; i-- {
		ch := c.checks[i]
		if ch.commandID == commandID {
			c.checks = append(c.checks[:i], c.checks[i+1:]...)
			c.Unlock()
			return true
		}
	}
	c.Unlock()
	return false
}

// Add lägger till en ny klient, den använder
// sig också av mutexes för att förhindra race conditions.
func (c *Client) Add(ch *Check) {
	c.Lock()
	c.checks = append(c.checks, ch)
	c.Unlock()
}

// Check har alltid senaste check och aldrig någon historik.
type Check struct {
	sync.Mutex              // Mutex så att man kan modifiera utan att det blir race conditions
	command       string    // Command att skicka
	gruppNamn     string    // Grupp namnet
	commandID     int       // ID för command, används bara för mysql
	nextCheck     int64     // Hur ofta denna check ska kollas
	pastID        int64     // Senaste IDn eller om en Check pågår, så är denna ID den tidigare checken
	nextTimestamp time.Time // När nästa timestamp ska ske
	checked       bool      // Om denna check har kollats
	err           bool      // Om förra checken failade
	failErr       bool      // Om det blir ett error, fortsätt kolla.
	done          bool      // Om denna check är färdig eller inte
}
