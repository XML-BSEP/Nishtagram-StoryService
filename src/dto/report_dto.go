package dto

import (
	"story-service/domain"
	"time"
)

type ReportDTO struct {
	Id string `json:"id" validate:"required"`
	PostId string `json:"postId" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	ReportBy domain.Profile `json:"reportedBy" validate:"required"`
	ReportType string `json:"reportType" validate:"required"`
	ReportedPostBy domain.Profile `json:"reportedPostBy" validate:"required"`
	ReportStatus string `json:"reportStatus" validate:"required"`
}

func NewReportDTO(id string, postId string, timestamp time.Time, reportedBy string,
	 reportedPostBy string, reportType string, reportStatus string) ReportDTO {
	return ReportDTO{
		Id: id,
		ReportBy: domain.Profile{Id: reportedBy},
		Timestamp: timestamp,
		PostId: postId,
		ReportedPostBy: domain.Profile{Id: reportedPostBy},
		ReportType: reportType,
		ReportStatus:reportStatus,
	}
}
