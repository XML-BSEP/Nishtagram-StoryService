package dto

type HighlightDTO struct {
	StoryId string `json:"story_id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
	HighlightName string `json:"highlight_name" validate:"required"`
}