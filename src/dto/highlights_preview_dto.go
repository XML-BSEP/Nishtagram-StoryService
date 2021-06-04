package dto

import "story-service/domain"

type HighlightsPreviewDTO struct {
	UserId string `json:"user_id" validate:"required"`
	MainStory domain.Media `json:"main_story" validate:"requrired"`
	HighlightName string `json:"highlight_name" validate:"required"`
}
