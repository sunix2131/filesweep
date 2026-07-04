package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	scan "filesweep/internal/domain/scan"
	"filesweep/internal/domain/settings"
	"filesweep/internal/infrastructure/events"

	"github.com/google/uuid"
)

type ScanResult struct {
	Session scan.Session
	Roots   []scan.Root
	Files   []scan.File
}

type Scanner struct {
	Publisher events.Publisher
}

func (s Scanner) Scan(ctx context.Context, roots []string, cfg settings.Settings) (ScanResult, error) {
	return s.ScanWithID(ctx, uuid.NewString(), roots, cfg)
}

func (s Scanner) ScanWithID(ctx context.Context, scanID string, roots []string, cfg settings.Settings) (ScanResult, error) {
	now := time.Now()
	session := scan.Session{ID: scanID, Status: scan.StatusDiscovering, StartedAt: now, CreatedAt: now, SelectedPaths: roots}
	result := ScanResult{Session: session}
	visited := map[string]bool{}
	for _, rootPath := range roots {
		if err := ctx.Err(); err != nil {
			session.Status = scan.StatusCancelled
			session.CancelledAt = time.Now()
			result.Session = session
			return result, err
		}
		abs, err := filepath.Abs(rootPath)
		if err != nil {
			continue
		}
		root := scan.Root{ID: uuid.NewString(), ScanSessionID: session.ID, RootPath: filepath.Clean(abs), RootName: filepath.Base(abs)}
		err = filepath.WalkDir(abs, func(path string, d os.DirEntry, walkErr error) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			if s.Publisher != nil {
				s.Publisher.ScanProgress(events.ScanProgress{ScanID: session.ID, Phase: "discovering", ProcessedFiles: session.FilesCount, CurrentPath: path})
			}
			if walkErr != nil {
				session.ErrorsCount++
				return nil
			}
			if path != abs && d.IsDir() {
				if !cfg.IncludeHiddenFiles && IsHidden(path) {
					return filepath.SkipDir
				}
				if slices.Contains(cfg.ScanExcludedFolderNames, d.Name()) {
					return filepath.SkipDir
				}
				if cfg.FollowSymlinks {
					clean, err := filepath.EvalSymlinks(path)
					if err == nil {
						if visited[clean] {
							return filepath.SkipDir
						}
						visited[clean] = true
					}
				}
				return nil
			}
			info, err := d.Info()
			if err != nil {
				session.ErrorsCount++
				return nil
			}
			isSymlink := info.Mode()&os.ModeSymlink != 0
			if isSymlink && !cfg.FollowSymlinks {
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if !cfg.IncludeHiddenFiles && isAnyHidden(abs, path) {
				return nil
			}
			rel, _ := filepath.Rel(abs, path)
			category, mimeType := CategoryFor(path)
			f := scan.File{
				ID: uuid.NewString(), ScanSessionID: session.ID, RootID: root.ID,
				AbsolutePath: filepath.Clean(path), RelativePath: rel, Name: filepath.Base(path),
				Extension: strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), "."), MimeType: mimeType,
				Category: category, SizeBytes: info.Size(), ModifiedAt: info.ModTime(),
				IsHidden: isAnyHidden(abs, path), IsAccessible: true, IsSymlink: isSymlink,
			}
			result.Files = append(result.Files, f)
			root.FilesCount++
			root.TotalSizeBytes += f.SizeBytes
			session.FilesCount++
			session.TotalSizeBytes += f.SizeBytes
			return nil
		})
		if err != nil && ctx.Err() != nil {
			session.Status = scan.StatusCancelled
			session.CancelledAt = time.Now()
			result.Session = session
			return result, err
		}
		result.Roots = append(result.Roots, root)
	}
	session.Status = scan.StatusAnalyzing
	result.Session = session
	return result, nil
}

func isAnyHidden(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return IsHidden(path)
	}
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func IsDangerousRoot(path string) bool {
	clean := filepath.Clean(path)
	home, _ := os.UserHomeDir()
	if home != "" && clean == filepath.Clean(home) {
		return true
	}
	volume := filepath.VolumeName(clean)
	withoutVol := strings.TrimPrefix(clean, volume)
	if withoutVol == string(filepath.Separator) || withoutVol == "." {
		return true
	}
	base := strings.ToLower(filepath.Base(clean))
	if base == "windows" || base == "program files" || base == "program files (x86)" {
		return true
	}
	return slices.Contains([]string{"/", "/usr", "/etc", "/var", "/System", "/Library"}, clean)
}
