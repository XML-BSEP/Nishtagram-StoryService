package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"story-service/domain"
	"story-service/dto"
	"strings"
	"time"
)


const (
	CreateReportTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Reports (id text, reported_by text, story_id text, story_by text, status text, type text, timestamp timestamp, PRIMARY KEY(status, id));"
	InsertIntoReportTable = "INSERT INTO story_keyspace.Reports (id, reported_by, story_id, story_by, status, type, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	UpdateReportTable = "UPDATE story_keyspace.Reports SET status = ? WHERE id = ? AND status = ?;"
	SelectAllTypes = "SELECT * FROM story_keyspace.ReportType LIMIT 300000000;"
	GetAllRequestsByStatus = "SELECT  id, story_id, timestamp, reported_by, story_by, type, status FROM story_keyspace.Reports WHERE status = ?;"
	DeleteReport = "DELETE FROM story_keyspace.Reports where status = ? and id = ?;"
	GetPendingReportById = "SELECT id, story_id, story_by, reported_by, type, timestamp FROM story_keyspace.Reports " +
		"WHERE status = ? AND id = ?;"
)
type ReportRepository interface {
	ReportStory(report domain.StoryReport, ctx context.Context) error
	GetAllReportTypes(ctx context.Context) ([]string, error)
	ReviewReport(report dto.ReviewReportDTO, ctx context.Context) error
	GetAllPendingReports(ctx context.Context) (*[]dto.ReportDTO, error)
	GetAllApprovedReports(ctx context.Context) (*[]dto.ReportDTO, error)
	GetAllRejectedReports(ctx context.Context) (*[]dto.ReportDTO, error)
}

type reportRepository struct {
	cassandraClient *gocql.Session
}

func (r reportRepository) ReviewReport(report dto.ReviewReportDTO, ctx context.Context) error {
	var reportId, postId, reportedPostBy, reportedBy, reportType string
	var timestamp time.Time

	iter := r.cassandraClient.Query(GetPendingReportById, "CREATED", report.ReportId).Iter()

	if iter == nil {
		return fmt.Errorf("no such element")
	}

	for iter.Scan(&reportId, &postId, &reportedPostBy, &reportedBy, &reportType, &timestamp) {

		if report.DeletePost {
			err := r.cassandraClient.Query(DeleteStory, true, reportedPostBy, postId).Exec()
			if err != nil {
				fmt.Println(err)
			}
		}

		updatedStatus := strings.ToUpper(report.Status)

		err := r.cassandraClient.Query(DeleteReport, "CREATED", reportId).Exec()
		var newUUID = uuid.NewString()

		err = r.cassandraClient.Query(InsertIntoReportTable, newUUID, reportedBy, postId, reportedPostBy, updatedStatus,  reportType, time.Now()).Exec()

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (r reportRepository) GetAllPendingReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	iter := r.cassandraClient.Query(GetAllRequestsByStatus, "CREATED").Iter().Scanner()

	if iter == nil {
		return nil, fmt.Errorf("no pending reports")
	}
	var reports []dto.ReportDTO
	var reportId, postId, reportedBy, reportedPostBy, reportType, status string
	var timestamp time.Time

	for iter.Next() {
		iter.Scan(&reportId, &postId, &timestamp, &reportedBy, &reportedPostBy, &reportType, &status)
		reports = append(reports, dto.NewReportDTO(reportId, postId, timestamp, reportedBy, reportedPostBy, reportType, status))
	}
	return &reports, nil
}

func (r reportRepository) GetAllApprovedReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	iter := r.cassandraClient.Query(GetAllRequestsByStatus, "APPROVED").Iter().Scanner()

	if iter == nil {
		return nil, fmt.Errorf("no pending reports")
	}
	var reports []dto.ReportDTO
	var reportId, postId, reportedBy, reportedPostBy, reportType, status string
	var timestamp time.Time

	for iter.Next() {
		iter.Scan(&reportId, &postId, &timestamp, &reportedBy, &reportedPostBy, &reportType, &status)
		reports = append(reports, dto.NewReportDTO(reportId, postId, timestamp, reportedBy, reportedPostBy, reportType, status))
	}
	return &reports, nil
}

func (r reportRepository) GetAllRejectedReports(ctx context.Context) (*[]dto.ReportDTO, error) {
	iter := r.cassandraClient.Query(GetAllRequestsByStatus, "REJECTED").Iter().Scanner()

	if iter == nil {
		return nil, fmt.Errorf("no pending reports")
	}
	var reports []dto.ReportDTO
	var reportId, postId, reportedBy, reportedPostBy, reportType, status string
	var timestamp time.Time

	for iter.Next() {
		iter.Scan(&reportId, &postId, &timestamp, &reportedBy, &reportedPostBy, &reportType, &status)
		reports = append(reports, dto.NewReportDTO(reportId, postId, timestamp, reportedBy, reportedPostBy, reportType, status))
	}
	return &reports, nil
}

func (r reportRepository) GetAllReportTypes(ctx context.Context) ([]string, error) {
	var retVal []string

	iter := r.cassandraClient.Query(SelectAllTypes).Iter().Scanner()

	var reportType string
	for iter.Next() {

		err := iter.Scan(&reportType)
		if err != nil {
			return nil, err
		}

		retVal = append(retVal, reportType)
	}

	return retVal, nil
}

func (r reportRepository) ReportStory(report domain.StoryReport, ctx context.Context) error {
	newUuid := uuid.NewString()

	err := r.cassandraClient.Query(InsertIntoReportTable, newUuid, report.ReportedBy.Id, report.StoryId,
		report.ReportedStoryBy.Id, "CREATED", strings.ToUpper(report.ReportType.Type), time.Now()).Exec()

	if err != nil {
		return err
	}

	return nil
}

func NewReportRepository(cassandraClient *gocql.Session) ReportRepository {

	r := reportRepository{cassandraClient: cassandraClient}

	err := r.cassandraClient.Query(CreateReportTable).Exec()

	if err != nil {
		return nil
	}

	return r
}
