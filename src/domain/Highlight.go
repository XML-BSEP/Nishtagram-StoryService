package domain

import "time"

type Highlight struct {
	Id string
	Name string
	Timestamp time.Time
	Profile Profile
	Stories []Story
}
