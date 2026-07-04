package scan

import "time"

type File struct {
	ID                   string    `json:"id"`
	ScanSessionID        string    `json:"scanSessionId"`
	RootID               string    `json:"rootId"`
	AbsolutePath         string    `json:"absolutePath"`
	RelativePath         string    `json:"relativePath"`
	Name                 string    `json:"name"`
	Extension            string    `json:"extension"`
	MimeType             string    `json:"mimeType"`
	Category             string    `json:"category"`
	SizeBytes            int64     `json:"sizeBytes"`
	ModifiedAt           time.Time `json:"modifiedAt"`
	CreatedAtIfAvailable time.Time `json:"createdAtIfAvailable,omitempty"`
	SHA256               string    `json:"sha256"`
	IsHidden             bool      `json:"isHidden"`
	IsAccessible         bool      `json:"isAccessible"`
	IsSymlink            bool      `json:"isSymlink"`
	IsUnstable           bool      `json:"isUnstable"`
	ScanError            string    `json:"scanError"`
}

const (
	CategoryImages        = "images"
	CategoryVideos        = "videos"
	CategoryAudio         = "audio"
	CategoryDocuments     = "documents"
	CategorySpreadsheets  = "spreadsheets"
	CategoryPresentations = "presentations"
	CategoryArchives      = "archives"
	CategoryInstallers    = "installers"
	CategoryCode          = "code"
	CategoryFonts         = "fonts"
	CategoryDiskImages    = "disk_images"
	CategoryOther         = "other"
)
