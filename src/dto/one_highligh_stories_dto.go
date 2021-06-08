package dto

import "story-service/domain"

type OneHighlightDTO struct {
	MainPicture domain.Media `json:"main_picture" validate:"required"`
	Stories []StoryDTO `json:"stories" validate:"required"`
	StoryId []string `json:"storyIds" validate:"required"`
	UserId string `json:"userId" validate:"required"`
	HighlightName string `json:"name" validate:"required"`
	HighlightPhoto string `json:"highlightPhoto" validate:"required"`
}