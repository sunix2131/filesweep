package scan

import "time"

type Status string

const (
	StatusQueued      Status = "queued"
	StatusDiscovering Status = "discovering"
	StatusAnalyzing   Status = "analyzing"
	StatusHashing     Status = "hashing"
	StatusCompleted   Status = "completed"
	StatusCancelled   Status = "cancelled"
	StatusFailed      Status = "failed"
)

type Session struct {
	ID                   string    `json:"id"`
	Status               Status    `json:"status"`
	StartedAt            time.Time `json:"startedAt"`
	CompletedAt          time.Time `json:"completedAt,omitempty"`
	CancelledAt          time.Time `json:"cancelledAt,omitempty"`
	SelectedPaths        []string  `json:"selectedPaths"`
	FilesCount           int       `json:"filesCount"`
	TotalSizeBytes       int64     `json:"totalSizeBytes"`
	DuplicateGroupsCount int       `json:"duplicateGroupsCount"`
	ReclaimableBytes     int64     `json:"reclaimableBytes"`
	ErrorsCount          int       `json:"errorsCount"`
	CreatedAt            time.Time `json:"createdAt"`
}

type Root struct {
	ID             string `json:"id"`
	ScanSessionID  string `json:"scanSessionId"`
	RootPath       string `json:"rootPath"`
	RootName       string `json:"rootName"`
	FilesCount     int    `json:"filesCount"`
	TotalSizeBytes int64  `json:"totalSizeBytes"`
}
