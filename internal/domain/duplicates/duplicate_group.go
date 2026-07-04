package duplicates

import scan "filesweep/internal/domain/scan"

type Group struct {
	ID                        string      `json:"id"`
	ScanSessionID             string      `json:"scanSessionId"`
	SHA256                    string      `json:"sha256"`
	FileSizeBytes             int64       `json:"fileSizeBytes"`
	FilesCount                int         `json:"filesCount"`
	EstimatedReclaimableBytes int64       `json:"estimatedReclaimableBytes"`
	RecommendedFileID         string      `json:"recommendedFileId"`
	Category                  string      `json:"category"`
	Members                   []scan.File `json:"members,omitempty"`
}
