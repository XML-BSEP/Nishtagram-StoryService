package usecase

import (
	"context"
	"fmt"
	logger "github.com/jelena-vlajkov/logger/logger"
	"story-service/domain"
	"story-service/dto"
	"story-service/gateway"
	"story-service/repository"
)

type ReportUseCase interface {
	ReportStory(report domain.StoryReport, ctx context.Context) error
	GetAllReportType(ctx context.Context) ([]string, error)
	ReviewReport(report dto.ReviewReportDTO, ctx context.Context) error
	GetAllPendingReports(ctx context.Context) (*[]dto.ReportDTO, error)
	GetAllApprovedReports(ctx context.Context) (*[]dto.ReportDTO, error)
	GetAllRejectedReports(ctx context.Context) (*[]dto.ReportDTO, error)
}

type reportUseCase struct {
	reportRepository repository.ReportRepository
	storyRepository repository.StoryRepo
	logger *logger.Logger
}

func (r reportUseCase) ReviewReport(report dto.ReviewReportDTO, ctx context.Context) error {
	return r.reportRepository.ReviewReport(report, context.Background())
}

func (r reportUseCase) GetAllPendingReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	reports, err := r.reportRepository.GetAllPendingReports(context.Background())
	if err != nil {
		return nil, err
	}

	for i, report := range *reports {
		reportedBy, err := gateway.GetUser(ctx, report.ReportBy.Id, r.logger)
		if err == nil {
			report.ReportBy = domain.Profile{Id: report.ReportBy.Id, ProfilePhoto: reportedBy.ProfilePhoto, Username: reportedBy.Username}
		}


		reportedUser, err := gateway.GetUser(ctx, report.ReportedPostBy.Id, r.logger)
		if err == nil {
			report.ReportedPostBy = domain.Profile{Id: report.ReportedPostBy.Id, ProfilePhoto: reportedUser.ProfilePhoto, Username: reportedUser.Username}
		}

		(*reports)[i] = report

	}

	return reports, nil
}

func (r reportUseCase) GetAllApprovedReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	reports, err := r.reportRepository.GetAllApprovedReports(context.Background())
	if err != nil {
		return nil, err
	}

	for i, report := range *reports {
		reportedBy, err := gateway.GetUser(ctx, report.ReportBy.Id, r.logger)
		if err == nil {
			report.ReportBy = domain.Profile{Id: report.ReportBy.Id, ProfilePhoto: reportedBy.ProfilePhoto, Username: reportedBy.Username}
		}


		reportedUser, err := gateway.GetUser(ctx, report.ReportedPostBy.Id, r.logger)
		if err == nil {
			report.ReportedPostBy = domain.Profile{Id: report.ReportedPostBy.Id, ProfilePhoto: reportedUser.ProfilePhoto, Username: reportedUser.Username}
		}

		(*reports)[i] = report

	}

	return reports, nil
}

func (r reportUseCase) GetAllRejectedReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	reports, err := r.reportRepository.GetAllRejectedReports(context.Background())
	if err != nil {
		return nil, err
	}

	for i, report := range *reports {
		reportedBy, err := gateway.GetUser(ctx, report.ReportBy.Id, r.logger)
		if err == nil {
			report.ReportBy = domain.Profile{Id: report.ReportBy.Id, ProfilePhoto: reportedBy.ProfilePhoto, Username: reportedBy.Username}
		}


		reportedUser, err := gateway.GetUser(ctx, report.ReportedPostBy.Id, r.logger)
		if err == nil {
			report.ReportedPostBy = domain.Profile{Id: report.ReportedPostBy.Id, ProfilePhoto: reportedUser.ProfilePhoto, Username: reportedUser.Username}
		}

		(*reports)[i] = report

	}

	return reports, nil
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

func NewReportUseCase(reportRepository repository.ReportRepository, storyRepository repository.StoryRepo, logger *logger.Logger) ReportUseCase {
	return &reportUseCase{reportRepository: reportRepository, storyRepository: storyRepository, logger: logger}
}
