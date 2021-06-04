package dto

type RemoveStoryDTO struct {
	UserId string `json:"user_id" validate:"required"`
	StoryId string `json:"story_id" validate:"required"`
}
