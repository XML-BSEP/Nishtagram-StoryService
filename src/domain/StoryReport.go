package domain

import "time"

type StoryReport struct {
	Id string
	StoryId string
	Timestamp time.Time
	ReportedBy Profile
	ReportType ReportType
	ReportStatus ReportStatus
	ReportedStoryBy Profile
}
