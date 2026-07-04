package application

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"filesweep/internal/domain/actions"
	"filesweep/internal/domain/duplicates"
	scan "filesweep/internal/domain/scan"
	"filesweep/internal/domain/settings"
	"filesweep/internal/infrastructure/apppaths"
	"filesweep/internal/infrastructure/database/repositories"
	"filesweep/internal/infrastructure/events"
	"filesweep/internal/infrastructure/filesystem"
	"filesweep/internal/infrastructure/platform"

	"github.com/google/uuid"
)

type Services struct {
	Store     *repositories.Store
	Paths     apppaths.Paths
	Logger    *slog.Logger
	Publisher events.Publisher
	Platform  platform.Services
	scanner   filesystem.Scanner
	mu        sync.Mutex
	cancels   map[string]context.CancelFunc
}

func NewServices(store *repositories.Store, paths apppaths.Paths, logger *slog.Logger, publisher events.Publisher) *Services {
	s := &Services{Store: store, Paths: paths, Logger: logger, Publisher: publisher, cancels: map[string]context.CancelFunc{}}
	s.scanner = filesystem.Scanner{Publisher: publisher}
	return s
}

func (s *Services) StartScan(ctx context.Context, paths []string) (scan.Session, error) {
	cfg, err := s.Store.GetSettings(ctx)
	if err != nil {
		return scan.Session{}, err
	}
	scanID := uuid.NewString()
	session := scan.Session{ID: scanID, Status: scan.StatusQueued, StartedAt: time.Now(), CreatedAt: time.Now(), SelectedPaths: paths}
	scanCtx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancels[scanID] = cancel
	s.mu.Unlock()
	if err := s.Store.SaveScan(ctx, session, nil, nil, nil); err != nil {
		return scan.Session{}, err
	}
	go s.runScan(scanCtx, scanID, paths, cfg)
	return session, nil
}

func (s *Services) runScan(scanCtx context.Context, scanID string, paths []string, cfg settings.Settings) {
	defer func() { s.mu.Lock(); delete(s.cancels, scanID); s.mu.Unlock() }()
	result, err := s.scanner.ScanWithID(scanCtx, scanID, paths, cfg)
	ctx := context.Background()
	if err != nil && scanCtx.Err() != nil {
		result.Session.Status = scan.StatusCancelled
		result.Session.CancelledAt = time.Now()
		_ = s.Store.SaveScan(ctx, result.Session, result.Roots, result.Files, nil)
		if s.Publisher != nil {
			s.Publisher.ScanCancelled(scanID)
		}
		return
	}
	if err != nil {
		result.Session.Status = scan.StatusFailed
		_ = s.Store.SaveScan(ctx, result.Session, result.Roots, result.Files, nil)
		if s.Publisher != nil {
			s.Publisher.ScanFailed(scanID, err.Error())
		}
		return
	}
	candidates := filesystem.DuplicateCandidates(result.Files)
	if s.Publisher != nil {
		s.Publisher.ScanProgress(events.ScanProgress{ScanID: scanID, Phase: "hashing", TotalFiles: len(candidates)})
	}
	candidates = filesystem.ApplyHashesWithProgress(candidates, filesystem.HashWorkers(cfg.MaxHashWorkers), filesystem.HashFile, func(processed int, total int, currentPath string) {
		if s.Publisher == nil {
			return
		}
		percent := 0.0
		if total > 0 {
			percent = float64(processed) / float64(total) * 100
		}
		s.Publisher.ScanProgress(events.ScanProgress{
			ScanID:         scanID,
			Phase:          "hashing",
			ProcessedFiles: processed,
			TotalFiles:     total,
			CurrentPath:    currentPath,
			Percent:        percent,
		})
	})
	if s.Publisher != nil {
		s.Publisher.ScanProgress(events.ScanProgress{ScanID: scanID, Phase: "saving", ProcessedFiles: len(result.Files), TotalFiles: len(result.Files), Percent: 95})
	}
	byID := map[string]scan.File{}
	for _, f := range result.Files {
		byID[f.ID] = f
	}
	for _, f := range candidates {
		byID[f.ID] = f
	}
	result.Files = result.Files[:0]
	for _, f := range byID {
		result.Files = append(result.Files, f)
	}
	groups := filesystem.BuildDuplicateGroups(scanID, result.Files)
	result.Session.Status = scan.StatusCompleted
	result.Session.CompletedAt = time.Now()
	result.Session.DuplicateGroupsCount = len(groups)
	for _, g := range groups {
		result.Session.ReclaimableBytes += g.EstimatedReclaimableBytes
	}
	if err := s.Store.SaveScan(ctx, result.Session, result.Roots, result.Files, groups); err != nil {
		if s.Publisher != nil {
			s.Publisher.ScanFailed(scanID, err.Error())
		}
		return
	}
	if s.Publisher != nil {
		s.Publisher.ScanProgress(events.ScanProgress{ScanID: scanID, Phase: "completed", ProcessedFiles: result.Session.FilesCount, TotalFiles: result.Session.FilesCount, Percent: 100})
		s.Publisher.ScanCompleted(scanID)
	}
}

func (s *Services) CancelScan(scanID string) {
	s.mu.Lock()
	cancel := s.cancels[scanID]
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (s *Services) ListScanSessions(ctx context.Context) ([]scan.Session, error) {
	return s.Store.ListScanSessions(ctx)
}
func (s *Services) GetScanSession(ctx context.Context, id string) (scan.Session, error) {
	return s.Store.GetScanSession(ctx, id)
}
func (s *Services) ListDuplicateGroups(ctx context.Context, scanID string, limit, offset int) ([]duplicates.Group, int, error) {
	return s.Store.ListDuplicateGroups(ctx, scanID, limit, offset)
}
func (s *Services) GetDuplicateGroup(ctx context.Context, id string) (duplicates.Group, error) {
	return s.Store.GetDuplicateGroup(ctx, id)
}
func (s *Services) ListFiles(ctx context.Context, scanID, category string, minSize int64, search string, limit, offset int) ([]scan.File, int, error) {
	return s.Store.ListFiles(ctx, scanID, category, minSize, search, limit, offset)
}
func (s *Services) Categories(ctx context.Context, scanID string) ([]map[string]interface{}, error) {
	return s.Store.CategorySummary(ctx, scanID)
}
func (s *Services) GetSettings(ctx context.Context) (settings.Settings, error) {
	return s.Store.GetSettings(ctx)
}
func (s *Services) SaveSettings(ctx context.Context, cfg settings.Settings) (settings.Settings, error) {
	if cfg.MaxHashWorkers < 1 {
		cfg.MaxHashWorkers = 1
	}
	if cfg.MaxPreviewFileSizeMB < 1 {
		cfg.MaxPreviewFileSizeMB = 20
	}
	return cfg, s.Store.SaveSettings(ctx, cfg)
}

func (s *Services) ExecuteAction(ctx context.Context, actionType actions.Type, items []actions.Item) (actions.Action, error) {
	a := actions.Action{ID: uuid.NewString(), ActionType: actionType, Status: actions.StatusRunning, CreatedAt: time.Now(), StartedAt: time.Now(), Items: items}
	for i := range a.Items {
		if a.Items[i].ID == "" {
			a.Items[i].ID = uuid.NewString()
		}
		a.Items[i].ActionID = a.ID
		switch actionType {
		case actions.TypeMoveToFolder:
			err := filesystem.SafeMove(a.Items[i].SourcePath, a.Items[i].TargetPath, a.Items[i].SourceSizeBytes, a.Items[i].SourceSHA256)
			if err != nil {
				a.Items[i].Status = actions.StatusFailed
				a.Items[i].ErrorMessage = err.Error()
			} else {
				a.Items[i].Status = actions.StatusCompleted
			}
		case actions.TypeMoveToTrash:
			err := s.Platform.MoveToTrash([]string{a.Items[i].SourcePath})
			if err != nil {
				a.Items[i].Status = actions.StatusFailed
				a.Items[i].ErrorMessage = err.Error()
			} else {
				a.Items[i].Status = actions.StatusCompleted
			}
		}
		if s.Publisher != nil {
			s.Publisher.ActionProgress(a.ID, i+1, len(a.Items))
		}
	}
	failed := 0
	for _, item := range a.Items {
		if item.Status == actions.StatusFailed {
			failed++
		}
	}
	a.CompletedAt = time.Now()
	a.Status = actions.StatusCompleted
	if failed > 0 {
		a.Status = actions.StatusFailed
		a.ErrorMessage = fmt.Sprintf("%d items failed", failed)
	}
	a.UndoAvailable = actionType == actions.TypeMoveToFolder && failed == 0
	a.Summary = fmt.Sprintf("%s: %d files", actionType, len(a.Items))
	if err := s.Store.SaveAction(ctx, a); err != nil {
		return a, err
	}
	if s.Publisher != nil {
		if a.Status == actions.StatusCompleted {
			s.Publisher.ActionCompleted(a.ID)
		} else {
			s.Publisher.ActionFailed(a.ID, a.ErrorMessage)
		}
	}
	return a, nil
}

func (s *Services) UndoAction(ctx context.Context, actionID string) (actions.Action, error) {
	orig, err := s.Store.GetAction(ctx, actionID)
	if err != nil {
		return actions.Action{}, err
	}
	if !orig.UndoAvailable || orig.ActionType != actions.TypeMoveToFolder {
		return actions.Action{}, fmt.Errorf("undo is unavailable")
	}
	var items []actions.Item
	for _, item := range orig.Items {
		items = append(items, actions.Item{SourcePath: item.TargetPath, TargetPath: item.SourcePath, SourceSizeBytes: item.SourceSizeBytes, SourceSHA256: item.SourceSHA256})
	}
	return s.ExecuteAction(ctx, actions.TypeMoveToFolder, items)
}

func (s *Services) ExportCSV(ctx context.Context, scanID string) (string, error) {
	files, _, err := s.Store.ListFiles(ctx, scanID, "", 0, "", 1000000, 0)
	if err != nil {
		return "", err
	}
	path := filepath.Join(s.Paths.ExportsDir, "filesweep-"+scanID+".csv")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{"name", "path", "size_bytes", "category", "modified_at", "sha256"})
	for _, file := range files {
		_ = w.Write([]string{file.Name, file.AbsolutePath, fmt.Sprint(file.SizeBytes), file.Category, file.ModifiedAt.Format(time.RFC3339), file.SHA256})
	}
	return path, w.Error()
}

func (s *Services) ListActions(ctx context.Context, limit, offset int) ([]actions.Action, int, error) {
	return s.Store.ListActions(ctx, limit, offset)
}

func (s *Services) ClearThumbnailCache() error {
	entries, err := os.ReadDir(s.Paths.ThumbnailsDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			_ = os.Remove(filepath.Join(s.Paths.ThumbnailsDir, e.Name()))
		}
	}
	return nil
}

func (s *Services) ImagePreview(path string) (string, error) {
	cfg, err := s.Store.GetSettings(context.Background())
	if err != nil {
		return "", err
	}
	return filesystem.CreateThumbnail(s.Paths.ThumbnailsDir, path, int64(cfg.MaxPreviewFileSizeMB)*1024*1024)
}
