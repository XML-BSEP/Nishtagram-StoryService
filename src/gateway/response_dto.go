package gateway

import "story-service/domain"

type FollowingResponseDTO struct {
	Data []domain.Profile `json:"data"`
}
