package server

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tryy3/MSTT-Monitor/server/providers"
	"github.com/tryy3/MSTT-Monitor/server/services"
)

func NewAlert(aof alertOptionFields) *Alert {
	alert := &Alert{
		rw:            new(sync.RWMutex),
		ID:            aof.ID,
		ClientID:      aof.ClientID,
		Delay:         aof.Delay,
		PreviousAlert: time.Time{},
	}

	switch aof.Alert {
	case "cpu":
		avg, err := strconv.ParseFloat(aof.Value, 64)
		if err != nil {
			return nil
		}
		alert.Alert = providers.NewAlertProviderCPU(aof.Count, avg)
	}

	for _, s := range strings.Split(aof.Service, ",") {
		switch s {
		case "sms":
			alert.Services = append(alert.Services, services.NewServiceSMS())
		case "email":
		}
	}
	return alert
}

type Alert struct {
	rw            *sync.RWMutex
	ID            int64
	ClientID      int64
	Delay         int64
	Alert         providers.AlertProvider
	PreviousAlert time.Time
	Services      []services.Service
}

func (a Alert) GetID() int64 {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.ID
}

func (a Alert) GetClientID() int64 {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.ClientID
}

func (a Alert) GetDelay() int64 {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.Delay
}

func (a Alert) GetAlert() providers.AlertProvider {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.Alert
}

func (a Alert) GetPreviousAlert() time.Time {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.PreviousAlert
}

func (a Alert) GetServices() []services.Service {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return a.Services
}

func (a *Alert) SetID(id int64) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.ID = id
}

func (a *Alert) SetClientID(clientID int64) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.ClientID = clientID
}

func (a *Alert) SetDelay(delay int64) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.Delay = delay
}

func (a *Alert) SetAlert(alert providers.AlertProvider) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.Alert = alert
}

func (a *Alert) SetPreviousAlert(previous time.Time) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.PreviousAlert = previous
}

func (a *Alert) SetServices(s []services.Service) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.Services = s
}

func (a *Alert) Update(alert alertOptionFields) bool {
	if alert.Delay != a.GetDelay() {
		a.SetDelay(alert.Delay)
	}

	switch alert.Alert {
	case "cpu":
		avg, err := strconv.ParseFloat(alert.Value, 64)
		if err != nil {
			return false
		}
		if strings.ToLower(a.GetAlert().Name()) != alert.Alert {
			a.SetAlert(providers.NewAlertProviderCPU(alert.Count, avg))
		} else {
			if a.GetAlert().(*providers.AlertProviderCPU).Total != alert.Count {
				a.GetAlert().(*providers.AlertProviderCPU).Total = alert.Count
			}
			if a.GetAlert().(*providers.AlertProviderCPU).Avg != avg {
				a.GetAlert().(*providers.AlertProviderCPU).Avg = avg
			}
		}
	}

	a.rw.Lock()
	defer a.rw.Unlock()
	a.Services = []services.Service{}
	for _, s := range strings.Split(alert.Service, ",") {
		switch s {
		case "sms":
			a.Services = append(a.Services, services.NewServiceSMS())
		case "email":
		}
	}
	return true
}

func (a *Alert) Check(resp string, db *Database) {
	a.rw.RLock()
	al := a.Alert
	a.rw.RUnlock()
	if al.Check(resp) {
		if time.Now().After(a.GetPreviousAlert()) {
			for s := range a.IterServices() {
				s.Send(a.Alert.Name(), a.Alert.Message())
			}
		}

		field, err := db.InsertAlert(a.GetID(), a.GetClientID(), al.Value())
		if err != nil {
			a.rw.Lock()
			a.PreviousAlert = time.Now().Add(time.Duration(a.Delay) * time.Second)
			a.rw.Unlock()
		} else {
			err := a.SetTimestampFromString(field.Timestamp)
			if err != nil {
				a.rw.Lock()
				a.PreviousAlert = time.Now().Add(time.Duration(a.Delay) * time.Second)
				a.rw.Unlock()
			}
		}
	}
}

func (a *Alert) SetTimestampFromString(t string) error {
	a.rw.Lock()
	defer a.rw.Unlock()
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	if err != nil {
		return err
	}
	timestamp = timestamp.Add(time.Duration(a.Delay) * time.Second)
	a.PreviousAlert = timestamp
	return nil
}

func (a Alert) CountServices() int {
	a.rw.RLock()
	defer a.rw.RUnlock()
	return len(a.Services)
}

func (a Alert) IterServices() <-chan services.Service {
	ch := make(chan services.Service, a.CountServices())
	go func() {
		a.rw.RLock()
		defer a.rw.RUnlock()
		for _, service := range a.Services {
			ch <- service
		}
		close(ch)
	}()
	return ch
}
