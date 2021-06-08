package dto

type HighlightsPreviewDTO struct {
	UserId string `json:"id" validate:"required"`
	HighlightPhoto string `json:"highlightPhoto" validate:"required"`
	HighlightName string `json:"name" validate:"required"`
}
