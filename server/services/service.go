package services

import "time"

type Service interface {
	Send(string, string)
}

func NewServiceSMS() Service {
	return &ServiceSMS{
		Recipients: []Recipient{
			Recipient{Msisdn: ""},
		},
		GW: newGatewayAPI("", time.Second*10),
	}
}
