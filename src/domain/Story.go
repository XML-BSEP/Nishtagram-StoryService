package domain

import (
	"time"
)

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
	IsCampaign bool
	CampaignId string
	Link string
}

func NewStory(id string, media string, timestamp time.Time, userId string, closeFriends bool, mentions []Profile, storyType string, locationName string, lat float64, long float64) Story {
	return Story{
		Id: id,
		Media: Media{Path: media},
		Timestamp: timestamp,
		Profile: Profile{Id: userId},
		Mentions: mentions,
		StoryType: StoryType{Type: storyType},
		Location: NewLocation(locationName, lat, long),
		CloseFriends: closeFriends,
		Banned: false,
	}
}


