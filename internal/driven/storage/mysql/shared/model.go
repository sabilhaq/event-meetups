package shared

import (
	"strconv"
	"strings"
	"time"

	"github.com/Haraj-backend/hex-monscape/internal/core/entity"
)

type MonsterRow struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Health        int    `db:"health"`
	MaxHealth     int    `db:"max_health"`
	Attack        int    `db:"attack"`
	Defense       int    `db:"defense"`
	Speed         int    `db:"speed"`
	AvatarURL     string `db:"avatar_url"`
	IsPartnerable int    `db:"is_partnerable"`
}

func (r *MonsterRow) ToMonster() *entity.Monster {
	return &entity.Monster{
		ID:   r.ID,
		Name: r.Name,
		BattleStats: entity.BattleStats{
			Health:    r.Health,
			MaxHealth: r.MaxHealth,
			Attack:    r.Attack,
			Defense:   r.Defense,
			Speed:     r.Speed,
		},
		AvatarURL: r.AvatarURL,
	}
}

type MonsterRows []MonsterRow

func (r MonsterRows) ToMonsters() []entity.Monster {
	var monsters []entity.Monster
	for _, row := range r {
		monsters = append(monsters, *row.ToMonster())
	}
	return monsters
}

func ToMonsterRow(monster *entity.Monster) *MonsterRow {
	return &MonsterRow{
		ID:        monster.ID,
		Name:      monster.Name,
		Health:    monster.BattleStats.Health,
		MaxHealth: monster.BattleStats.MaxHealth,
		Attack:    monster.BattleStats.Attack,
		Defense:   monster.BattleStats.Defense,
		Speed:     monster.BattleStats.Speed,
		AvatarURL: monster.AvatarURL,
	}
}

type EventRow struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (r *EventRow) ToEvent() *entity.Event {
	return &entity.Event{
		ID:   r.ID,
		Name: r.Name,
	}
}

type EventRows []EventRow

func (r EventRows) ToEvents() []entity.Event {
	var monsters []entity.Event
	for _, row := range r {
		monsters = append(monsters, *row.ToEvent())
	}
	return monsters
}

type SupportedEventRows []SupportedEventRow

type SupportedEventRow struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	MeetupsCapacity int    `db:"meetups_capacity"`
}

func (r SupportedEventRows) ToSupportedEvents() []entity.SupportedEvent {
	var supportedEvents []entity.SupportedEvent
	for _, row := range r {
		supportedEvents = append(supportedEvents, *row.ToSupportedEvent())
	}
	return supportedEvents
}

// Helper function to add event to existing venue
func (r *SupportedEventRow) ToSupportedEvent() *entity.SupportedEvent {
	return &entity.SupportedEvent{
		ID:              r.ID,
		Name:            r.Name,
		MeetupsCapacity: r.MeetupsCapacity,
	}
}

type VenueEventRows []VenueEventRow

type VenueEventRow struct {
	VenueID         int    `db:"venue_id"`
	VenueName       string `db:"venue_name"`
	VenueOpenDays   string `db:"venue_open_days"`
	VenueOpenAt     string `db:"venue_open_at"`
	VenueClosedAt   string `db:"venue_closed_at"`
	VenueTimezone   string `db:"venue_timezone"`
	EventID         int    `db:"event_id"`
	EventName       string `db:"event_name"`
	MeetupsCapacity int    `db:"meetups_capacity"`
}

func (r *VenueEventRow) ToVenue() *entity.Venue {
	var openDaysArr []int
	for _, day := range strings.Split(r.VenueOpenDays, ",") {
		dayInt, _ := strconv.Atoi(day)
		openDaysArr = append(openDaysArr, dayInt)
	}

	return &entity.Venue{
		ID:       r.VenueID,
		Name:     r.VenueName,
		OpenDays: openDaysArr,
		OpenAt:   r.VenueOpenAt,
		ClosedAt: r.VenueClosedAt,
		Timezone: r.VenueTimezone,
		SupportedEvents: []entity.SupportedEvent{
			{
				ID:              r.EventID,
				Name:            r.EventName,
				MeetupsCapacity: r.MeetupsCapacity,
			},
		},
	}
}

func (r VenueEventRows) ToVenues() []entity.Venue {
	// Map to hold venues and their events
	venuesMap := make(map[int]*entity.Venue)

	for _, row := range r {
		// Check if the venue already exists in the map
		if venue, exists := venuesMap[row.VenueID]; exists {
			// Append the event to the existing venue's supported events
			venue.SupportedEvents = append(venue.SupportedEvents, entity.SupportedEvent{
				ID:              row.EventID,
				Name:            row.EventName,
				MeetupsCapacity: row.MeetupsCapacity,
			})
		} else {
			// Create a new venue entry
			venuesMap[row.VenueID] = row.ToVenue()
		}
	}

	var venues []entity.Venue
	for _, venue := range venuesMap {
		venues = append(venues, *venue)
	}
	return venues
}

type MeetupUserRow struct {
	MeetupID int   `db:"meetup_id"`
	UserID   int   `db:"user_id"`
	JoinedAt int64 `db:"joined_at"`
}

func NewMeetupUserRow(meetup *entity.MeetupUser) *MeetupUserRow {
	now := time.Now().Unix()
	return &MeetupUserRow{
		MeetupID: meetup.MeetupID,
		UserID:   meetup.UserID,
		JoinedAt: now,
	}
}
