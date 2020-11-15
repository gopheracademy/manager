package def

// Conference is a brand like GopherCon
type Conference struct {
	// pk: "true"
	ID     uint
	Name   string
	Events []Event
}

// Event is an instance like GopherCon 2020
type Event struct {
	// pk: "true"
	ID   uint
	Name string
	Slug string
	//	StartDate time.Time
	//	EndDate   time.Time
	Location string
	//	Slots    []EventSlot
}
