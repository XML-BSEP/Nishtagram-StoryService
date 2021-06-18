package repository

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"story-service/domain"
	"strings"
	"time"
)


const (
	CreateReportTable = "CREATE TABLE IF NOT EXISTS story_keyspace.Reports (id text, reported_by text, story_id text, story_by text, status text, type text, timestamp timestamp, PRIMARY KEY(status, id));"
	InsertIntoReportTable = "INSERT INTO story_keyspace.Reports (id, reported_by, story_id, story_by, status, type, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	UpdateReportTable = "UPDATE story_keyspace.Reports SET status = ? WHERE id = ? AND status = ?;"
	SelectAllTypes = "SELECT * FROM story_keyspace.ReportType LIMIT 300000000;"
)
type ReportRepository interface {
	ReportStory(report domain.StoryReport, ctx context.Context) error
	GetAllReportTypes(ctx context.Context) ([]string, error)
}

type reportRepository struct {
	cassandraClient *gocql.Session
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
