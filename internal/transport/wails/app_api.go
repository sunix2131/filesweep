package wails

import (
	"context"
	"errors"
	"fmt"

	"filesweep/internal/bootstrap"
	"filesweep/internal/domain/actions"
	"filesweep/internal/domain/duplicates"
	scan "filesweep/internal/domain/scan"
	"filesweep/internal/domain/settings"
	"filesweep/internal/infrastructure/events"
	"filesweep/internal/infrastructure/filesystem"
	"filesweep/internal/transport/wails/dto"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppAPI struct {
	ctx       context.Context
	container *bootstrap.Container
}

func NewAppAPI(container *bootstrap.Container) *AppAPI { return &AppAPI{container: container} }

func (a *AppAPI) Startup(ctx context.Context) {
	a.ctx = ctx
	a.container.Services.Publisher = a
}

func (a *AppAPI) SelectFolders() (dto.FolderSelectionResult, error) {
	if a.ctx == nil {
		return dto.FolderSelectionResult{}, errors.New("application is not ready")
	}
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Select folder"})
	if err != nil {
		return dto.FolderSelectionResult{}, err
	}
	res := dto.FolderSelectionResult{}
	if path != "" {
		res.Paths = []string{path}
		if filesystem.IsDangerousRoot(path) {
			res.HasDangerousPath = true
			res.DangerousPathNote = "Selected folder may contain system or home-directory files. Confirm before scanning."
		}
	}
	return res, nil
}

func (a *AppAPI) StartScan(req dto.StartScanRequest) (scan.Session, error) {
	for _, p := range req.Paths {
		if filesystem.IsDangerousRoot(p) && !req.ConfirmDangerousFolders {
			return scan.Session{}, fmt.Errorf("dangerous folder requires explicit confirmation: %s", p)
		}
	}
	return a.container.Services.StartScan(context.Background(), req.Paths)
}

func (a *AppAPI) CancelScan(scanID string) { a.container.Services.CancelScan(scanID) }
func (a *AppAPI) GetScanSession(scanID string) (scan.Session, error) {
	return a.container.Services.GetScanSession(context.Background(), scanID)
}
func (a *AppAPI) ListScanSessions() ([]scan.Session, error) {
	return a.container.Services.ListScanSessions(context.Background())
}
func (a *AppAPI) GetScanSummary(scanID string) (scan.Session, error) { return a.GetScanSession(scanID) }

func (a *AppAPI) ListDuplicateGroups(q dto.PageQuery) (dto.PaginatedDuplicateGroupsDto, error) {
	limit := normalizeLimit(q.Limit)
	items, total, err := a.container.Services.ListDuplicateGroups(context.Background(), q.ScanID, limit, q.Offset)
	return dto.PaginatedDuplicateGroupsDto{Items: items, Total: total}, err
}

func (a *AppAPI) GetDuplicateGroup(groupID string) (duplicates.Group, error) {
	return a.container.Services.GetDuplicateGroup(context.Background(), groupID)
}

func (a *AppAPI) ListLargeFiles(q dto.PageQuery) (dto.PaginatedFilesDto, error) {
	items, total, err := a.container.Services.ListFiles(context.Background(), q.ScanID, q.Category, q.MinSizeBytes, q.Search, normalizeLimit(q.Limit), q.Offset)
	return dto.PaginatedFilesDto{Items: items, Total: total}, err
}

func (a *AppAPI) ListCategories(scanID string) ([]map[string]interface{}, error) {
	return a.container.Services.Categories(context.Background(), scanID)
}

func (a *AppAPI) GetCategoryFiles(q dto.PageQuery) (dto.PaginatedFilesDto, error) {
	return a.ListLargeFiles(q)
}
func (a *AppAPI) BuildActionPlan(req dto.BuildActionPlanRequest) (actions.Action, error) {
	return actions.Action{ID: "pending", ActionType: req.ActionType, Status: actions.StatusPlanned, Items: req.Items, Summary: fmt.Sprintf("%d files planned", len(req.Items))}, nil
}
func (a *AppAPI) ExecuteActionPlan(req dto.BuildActionPlanRequest) (actions.Action, error) {
	return a.container.Services.ExecuteAction(context.Background(), req.ActionType, req.Items)
}
func (a *AppAPI) UndoAction(actionID string) (actions.Action, error) {
	return a.container.Services.UndoAction(context.Background(), actionID)
}
func (a *AppAPI) ListActionHistory(q dto.PageQuery) ([]actions.Action, error) {
	items, _, err := a.container.Services.ListActions(context.Background(), normalizeLimit(q.Limit), q.Offset)
	return items, err
}
func (a *AppAPI) RevealInFileManager(path string) error {
	return a.container.Services.Platform.RevealInFileManager(path)
}
func (a *AppAPI) OpenFile(path string) error { return a.container.Services.Platform.OpenFile(path) }
func (a *AppAPI) ExportScanReport(scanID string, format string) (dto.ExportResultDto, error) {
	if format != "csv" {
		return dto.ExportResultDto{}, fmt.Errorf("unsupported export format: %s", format)
	}
	path, err := a.container.Services.ExportCSV(context.Background(), scanID)
	return dto.ExportResultDto{Path: path}, err
}
func (a *AppAPI) GetSettings() (settings.Settings, error) {
	return a.container.Services.GetSettings(context.Background())
}
func (a *AppAPI) UpdateSettings(cfg settings.Settings) (settings.Settings, error) {
	return a.container.Services.SaveSettings(context.Background(), cfg)
}
func (a *AppAPI) ClearThumbnailCache() error { return a.container.Services.ClearThumbnailCache() }
func (a *AppAPI) GetImagePreview(path string) (string, error) {
	return a.container.Services.ImagePreview(path)
}

func (a *AppAPI) ScanProgress(event events.ScanProgress) {
	runtime.EventsEmit(a.ctx, "scan:progress", event)
}
func (a *AppAPI) ScanCompleted(scanID string) { runtime.EventsEmit(a.ctx, "scan:completed", scanID) }
func (a *AppAPI) ScanCancelled(scanID string) { runtime.EventsEmit(a.ctx, "scan:cancelled", scanID) }
func (a *AppAPI) ScanFailed(scanID string, message string) {
	runtime.EventsEmit(a.ctx, "scan:failed", map[string]string{"scanId": scanID, "message": message})
}
func (a *AppAPI) ActionProgress(actionID string, processed int, total int) {
	runtime.EventsEmit(a.ctx, "action:progress", map[string]interface{}{"actionId": actionID, "processed": processed, "total": total})
}
func (a *AppAPI) ActionCompleted(actionID string) {
	runtime.EventsEmit(a.ctx, "action:completed", actionID)
}
func (a *AppAPI) ActionFailed(actionID string, message string) {
	runtime.EventsEmit(a.ctx, "action:failed", map[string]string{"actionId": actionID, "message": message})
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 100
	}
	if limit > 1000 {
		return 1000
	}
	return limit
}
