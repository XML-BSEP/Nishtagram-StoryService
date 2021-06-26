package dto

type ReviewReportDTO struct {
	ReportId string `json:"reportId" validate:"required"`
	Status string `json:"status" validate:"required"`
	DeletePost bool `json:"deletePost" validate:"required"`
}
