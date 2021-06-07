package dto

type StoryContent struct {
	IsVideo bool `json:"isVideo" validate:"required"`
	Content string `json:"content" validate:"required"`
}
