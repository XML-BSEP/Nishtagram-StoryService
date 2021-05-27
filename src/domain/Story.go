package domain

import "time"

type Story struct {
	Id uint
	Media Media
	Timestamp time.Time
	StoryType StoryType
	CloseFriends bool
	Banned bool
}
