package def

import "time"

// Conference is a brand like GopherCon
type Conference struct {
	ID     uint
	Name   string
	Events []Event
}

// Event is an instance like GopherCon 2020
type Event struct {
	ID        uint
	Name      string
	Slug      string
	StartDate time.Time
	EndDate   time.Time
	Location  string
}
