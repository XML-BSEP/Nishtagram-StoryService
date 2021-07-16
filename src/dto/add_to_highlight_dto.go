package dto

type HighlightDTO struct {
	StoryId []string `json:"storyIds" validate:"required"`
	UserId string `json:"userId" validate:"required"`
	HighlightName string `json:"name" validate:"required"`
	HighlightPhoto string `json:"highlightPhoto" validate:"required"`
	Stories []StoryDTO `json:"stories" validate:"required"`
	Id string `json:"id"`
}