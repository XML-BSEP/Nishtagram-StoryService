package usecase

import (
	"context"
	"fmt"
	"story-service/domain"
	"story-service/repository"
)

type ReportUseCase interface {
	ReportStory(report domain.StoryReport, ctx context.Context) error
	GetAllReportType(ctx context.Context) ([]string, error)
}

type reportUseCase struct {
	reportRepository repository.ReportRepository
	storyRepository repository.StoryRepo
}

func (r reportUseCase) GetAllReportType(ctx context.Context) ([]string, error) {
	return r.reportRepository.GetAllReportTypes(ctx)
}

func (r reportUseCase) ReportStory(report domain.StoryReport, ctx context.Context) error {
	if !r.storyRepository.SeeIfExists(ctx, report.ReportedStoryBy.Id, report.StoryId) {
		return fmt.Errorf("no such story")
	}

	err := r.reportRepository.ReportStory(report, ctx)

	if err != nil {
		return err
	}

	return nil
}

func NewReportUseCase(reportRepository repository.ReportRepository, storyRepository repository.StoryRepo) ReportUseCase {
	return &reportUseCase{reportRepository: reportRepository, storyRepository: storyRepository}
}
