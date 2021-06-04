package dto

import "story-service/domain"

type OneHighlightDTO struct {
	UserId string `json:"user_id" validate:"required"`
	HighlightName string `json:"highlight_name" validate:"required"`
	MainPicture domain.Media `json:"main_picture" validate:"required"`
	Stories []StoryDTO `json:"stories" validate:"required"`
}