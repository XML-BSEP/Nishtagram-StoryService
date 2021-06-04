package usecase

import (
	"context"
	"fmt"
	"story-service/domain"
	"story-service/dto"
	"story-service/repository"
)

type HighlightUseCase interface {
	AddStoryToHighlight(ctx context.Context, dto dto.HighlightDTO) error
	RemoveStoryFrom(ctx context.Context, dto dto.HighlightDTO) error
	GetHighlights(ctx context.Context, userId string) ([]dto.HighlightsPreviewDTO, error)
	GetHighlightByName(ctx context.Context, userId string, highlightName string) (dto.OneHighlightDTO, error)
}

type highlightUseCase struct {
	highlightRepository repository.HighlightRepo
	storyRepository repository.StoryRepo
}

func (h highlightUseCase) AddStoryToHighlight(ctx context.Context, dto dto.HighlightDTO) error {
	if !h.storyRepository.SeeIfExists(context.Background(), dto.UserId, dto.StoryId) {
		return fmt.Errorf("no such story")
	}
	stories, _, _ := h.highlightRepository.GetHighlightByName(context.Background(), dto.UserId, dto.HighlightName)
	for _, s := range stories {
		if s == dto.StoryId {
			return nil
		}
	}
	return h.AddStoryToHighlight(context.Background(), dto)
}

func (h highlightUseCase) RemoveStoryFrom(ctx context.Context, dto dto.HighlightDTO) error {
	if !h.storyRepository.SeeIfExists(context.Background(), dto.UserId, dto.StoryId) {
		return fmt.Errorf("no such story")
	}
	stories, _, _ := h.highlightRepository.GetHighlightByName(context.Background(), dto.UserId, dto.HighlightName)
	for i, s := range stories {
		if s == dto.StoryId {
			break
		}

		if i == len(stories) - 1 {
			return fmt.Errorf("no such story")
		}
	}
	return h.RemoveStoryFrom(context.Background(), dto)
}

func (h highlightUseCase) GetHighlights(ctx context.Context, userId string) ([]dto.HighlightsPreviewDTO, error) {
	return h.highlightRepository.GetAllHighlightsByUser(context.Background(), userId)
}

func (h highlightUseCase) GetHighlightByName(ctx context.Context, userId string, highlightName string) (dto.OneHighlightDTO, error) {
	stories, mainStory, _ := h.highlightRepository.GetHighlightByName(context.Background(), userId, highlightName)
	if len(stories) > 0 {
		var storiesToReturn []dto.StoryDTO
		for _, s := range stories {
			story, err := h.storyRepository.GetStoryById(context.Background(), userId, s)
			if err != nil {
				continue
			}
			storiesToReturn = append(storiesToReturn, story)
		}
		return dto.OneHighlightDTO{
			MainPicture: domain.Media{Path: mainStory},
			UserId: userId,
			Stories: storiesToReturn,
		}, nil
	}
	return dto.OneHighlightDTO{}, fmt.Errorf("no such highlight")


}

func NewHighlightUseCase(highlightRepo repository.HighlightRepo, storyRepository repository.StoryRepo) HighlightUseCase {
	return &highlightUseCase{highlightRepository: highlightRepo, storyRepository: storyRepository}
}