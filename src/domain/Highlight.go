package domain

import "time"

type Highlight struct {
	Id uint
	Name string
	Timestamp time.Time
	Profile Profile
	Stories []Story
}
