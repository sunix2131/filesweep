package actions

import "time"

type Type string
type Status string

const (
	TypeMoveToFolder Type = "move_to_folder"
	TypeMoveToTrash  Type = "move_to_trash"
	TypeUndoMove     Type = "undo_move"
	TypeExportReport Type = "export_report"

	StatusPlanned   Status = "planned"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusSkipped   Status = "skipped"
)

type Action struct {
	ID            string    `json:"id"`
	ActionType    Type      `json:"actionType"`
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	StartedAt     time.Time `json:"startedAt,omitempty"`
	CompletedAt   time.Time `json:"completedAt,omitempty"`
	UndoAvailable bool      `json:"undoAvailable"`
	UndoExpiresAt time.Time `json:"undoExpiresAt,omitempty"`
	Summary       string    `json:"summary"`
	ErrorMessage  string    `json:"errorMessage"`
	Items         []Item    `json:"items,omitempty"`
}

type Item struct {
	ID              string `json:"id"`
	ActionID        string `json:"actionId"`
	SourcePath      string `json:"sourcePath"`
	TargetPath      string `json:"targetPath"`
	SourceSizeBytes int64  `json:"sourceSizeBytes"`
	SourceSHA256    string `json:"sourceSha256"`
	Status          Status `json:"status"`
	ErrorMessage    string `json:"errorMessage"`
}
