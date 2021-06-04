package usecase

import (
	"context"
	"fmt"
	"story-service/dto"
	"story-service/mapper"
	"story-service/repository"
	"time"
)

type StoryUseCase interface {
	AddStory(ctx context.Context, dto dto.StoryDTO) error
	RemoveStory(ctx context.Context, dto dto.RemoveStoryDTO) error
	GetAllStoriesForOneUser(ctx context.Context, userId string) ([]dto.StoryDTO, error)
}

type storyUseCase struct {
	storyRepository repository.StoryRepo
	redisUseCase RedisUseCase
}

func (s storyUseCase) AddStory(ctx context.Context, dto dto.StoryDTO) error {
	var keyToStore, valueToStore string
	expiresAtTime := time.Now().Add(time.Hour*24)
	expiresAt := time.Unix(expiresAtTime.Unix(), 0)
	now := time.Now()

	keyToStore = dto.UserId + "/" + dto.StoryId
	valueToStore = dto.StoryId

	s.redisUseCase.AddKeyValueSet(context.Background(), keyToStore, valueToStore, now.Sub(expiresAt))
	return s.storyRepository.AddStory(context.Background(), mapper.MapDTOToStory(dto))
}

func (s storyUseCase) RemoveStory(ctx context.Context, dto dto.RemoveStoryDTO) error {
	if !s.storyRepository.SeeIfExists(context.Background(), dto.UserId, dto.StoryId) {
		return fmt.Errorf("no cuch story")
	}

	return s.storyRepository.RemoveStory(context.Background(), dto.UserId, dto.StoryId)
}

func (s storyUseCase) GetAllStoriesForOneUser(ctx context.Context, userId string) ([]dto.StoryDTO, error) {
	var userFollowing []string
	var retVal []dto.StoryDTO

	for _, u := range userFollowing {
		searchPattern := u + "*"
		keys, err := s.redisUseCase.ScanKeyByPattern(context.Background(), searchPattern)
		if err != nil {
			continue
		}

		for _, k := range keys {
			storyId, err := s.redisUseCase.GetValueByKey(context.Background(), k)
			if err != nil {
				continue
			}
			story, err := s.storyRepository.GetStoryById(context.Background(), u, storyId)
			if err != nil {
				continue
			}
			retVal = append(retVal, story)
		}

	}

	return retVal, nil

}

func NewStoryUseCase(storyRepository repository.StoryRepo, useCase RedisUseCase) StoryUseCase {
	return &storyUseCase{
		storyRepository: storyRepository,
		redisUseCase: useCase,
	}
}
