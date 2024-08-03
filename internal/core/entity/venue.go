package entity

type Venue struct {
	ID              string
	Name            string
	OpenDays        []int
	OpenAt          string
	ClosedAt        string
	TimeZone        string
	SupportedEvents []SupportedEvent
}

type SupportedEvent struct {
	ID            string
	Name          string
	EventCapacity int
}
