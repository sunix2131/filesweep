package dto

import (
	"filesweep/internal/domain/actions"
	scan "filesweep/internal/domain/scan"
	"filesweep/internal/domain/settings"
)

type FolderSelectionResult struct {
	Paths             []string `json:"paths"`
	HasDangerousPath  bool     `json:"hasDangerousPath"`
	DangerousPathNote string   `json:"dangerousPathNote"`
}

type StartScanRequest struct {
	Paths                   []string `json:"paths"`
	ConfirmDangerousFolders bool     `json:"confirmDangerousFolders"`
}

type PageQuery struct {
	ScanID       string `json:"scanId"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	Category     string `json:"category"`
	Search       string `json:"search"`
	MinSizeBytes int64  `json:"minSizeBytes"`
}

type PaginatedFilesDto struct {
	Items []scan.File `json:"items"`
	Total int         `json:"total"`
}

type PaginatedDuplicateGroupsDto struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

type BuildActionPlanRequest struct {
	ActionType actions.Type   `json:"actionType"`
	Items      []actions.Item `json:"items"`
}

type ExportResultDto struct {
	Path string `json:"path"`
}

type SettingsDto = settings.Settings
