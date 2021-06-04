package domain

import (
	"story-service/dto"
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

func NewStoryFromDTO(dto dto.StoryDTO) Story {
	var mentions []Profile
	for _, s := range dto.Mentions {
		mentions = append(mentions, Profile{Id: s})
	}
	return Story{
		Id: dto.StoryId,
		Profile: Profile{Id: dto.UserId},
		Timestamp: dto.Timestamp,
		Location: dto.Location,
		Mentions: mentions,
		CloseFriends: dto.CloseFriends,
		Banned: false,
		StoryType: StoryType{Type: dto.Type},
		Media: dto.MediaPath,
	}
}
