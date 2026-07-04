package events

type ScanProgress struct {
	ScanID         string  `json:"scanId"`
	Phase          string  `json:"phase"`
	ProcessedFiles int     `json:"processedFiles"`
	TotalFiles     int     `json:"totalFiles"`
	CurrentPath    string  `json:"currentPath"`
	Percent        float64 `json:"percent"`
}

type Publisher interface {
	ScanProgress(event ScanProgress)
	ScanCompleted(scanID string)
	ScanCancelled(scanID string)
	ScanFailed(scanID string, message string)
	ActionProgress(actionID string, processed int, total int)
	ActionCompleted(actionID string)
	ActionFailed(actionID string, message string)
}
