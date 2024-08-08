package entity

type Venue struct {
	ID              int
	Name            string
	OpenDays        []int
	OpenAt          string
	ClosedAt        string
	Timezone        string
	SupportedEvents []SupportedEvent
}

type SupportedEvent struct {
	ID              int
	Name            string
	MeetupsCapacity int
}

type GetVenueFilter struct {
	EventID       *int
	MeetupStartTS *string
	MeetupEndTS   *string
}

type VenueEvent struct {
	VenueID         int
	EventID         int
	MeetupsCapacity int
}
