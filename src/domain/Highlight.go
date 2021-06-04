package domain

import "time"

type Highlight struct {
	Name string
	Timestamp time.Time
	Profile Profile
	Stories []Story
	MainStory Media
}


