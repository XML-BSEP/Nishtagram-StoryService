package domain

import "time"

type Story struct {
	Id string
	Media Media
	Timestamp time.Time
	Profile Profile
	Mentions []Profile
	StoryType StoryType
	Location Location
	CloseFriends bool
	Banned bool
}
