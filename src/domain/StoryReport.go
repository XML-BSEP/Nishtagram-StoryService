package domain

import "time"

type StoryReport struct {
	Id uint
	StoryId uint
	Timestamp time.Time
	ReportedBy Profile
	ReportType ReportType
	ReportStatus ReportStatus
}
