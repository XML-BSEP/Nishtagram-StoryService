package dto

import (
	"story-service/domain"
	"time"
)

type StoryDTO struct {
	StoryId string `json:"story_id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
	Mentions []string `json:"mentions" validate: "required"`
	MediaPath domain.Media `json:"media_path" validate:"required"`
	Type string `json:"type" validate:"required"`
	Location domain.Location `json:"location" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	CloseFriends bool `json:"close_friends" validate:"required"`

}

func NewStoryDTO(storyId string, userId string, mentions []string, media domain.Media, mediaType string, location domain.Location, timestamp time.Time, closeFriends bool) StoryDTO {
	return StoryDTO{
		StoryId: storyId,
		CloseFriends: closeFriends,
		UserId: userId,
		Mentions: mentions,
		MediaPath: media,
		Type: mediaType,
		Location: location,
		Timestamp: timestamp,
	}
}
