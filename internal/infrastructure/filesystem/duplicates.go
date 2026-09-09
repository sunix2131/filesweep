package filesystem

import (
	"context"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"filesweep/internal/domain/duplicates"
	scan "filesweep/internal/domain/scan"

	"github.com/google/uuid"
)

func HashWorkers(setting int) int {
	n := runtime.NumCPU() - 1
	if n < 2 {
		n = 2
	}
	if n > 8 {
		n = 8
	}
	if setting > 0 && setting < n {
		n = setting
	}
	return n
}

func DuplicateCandidates(files []scan.File) []scan.File {
	bySize := map[int64]int{}
	for _, f := range files {
		if f.IsAccessible && !f.IsUnstable && f.SizeBytes > 0 {
			bySize[f.SizeBytes]++
		}
	}
	var out []scan.File
	for _, f := range files {
		if bySize[f.SizeBytes] > 1 {
			out = append(out, f)
		}
	}
	return out
}

func BuildDuplicateGroups(scanID string, files []scan.File) []duplicates.Group {
	byHash := map[string][]scan.File{}
	for _, f := range files {
		if f.SHA256 != "" && !f.IsUnstable {
			byHash[f.SHA256] = append(byHash[f.SHA256], f)
		}
	}
	var out []duplicates.Group
	for hash, members := range byHash {
		if len(members) < 2 {
			continue
		}
		sort.Slice(members, func(i, j int) bool { return members[i].AbsolutePath < members[j].AbsolutePath })
		rec := RecommendedFile(members)
		out = append(out, duplicates.Group{
			ID: uuid.NewString(), ScanSessionID: scanID, SHA256: hash,
			FileSizeBytes: members[0].SizeBytes, FilesCount: len(members),
			EstimatedReclaimableBytes: members[0].SizeBytes * int64(len(members)-1),
			RecommendedFileID:         rec.ID, Category: members[0].Category, Members: members,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EstimatedReclaimableBytes > out[j].EstimatedReclaimableBytes })
	return out
}

var copyPattern = regexp.MustCompile(`(?i)(\([0-9]+\)|copy|копия|duplicate|final-final|new)`)

func RecommendedFile(files []scan.File) scan.File {
	best := files[0]
	for _, f := range files[1:] {
		if scoreFile(f) < scoreFile(best) {
			best = f
		}
	}
	return best
}

func scoreFile(f scan.File) int {
	path := strings.ToLower(f.AbsolutePath)
	score := 0
	if strings.Contains(path, string(filepath.Separator)+"downloads"+string(filepath.Separator)) {
		score += 1000
	}
	if strings.Contains(path, "tmp") || strings.Contains(path, "temp") {
		score += 500
	}
	if copyPattern.MatchString(f.Name) {
		score += 150
	}
	score += len(f.Name) / 2
	score += len(f.RelativePath) / 4
	score += int(f.ModifiedAt.Unix() / 86400 / 100000)
	return score
}

func ApplyHashes(files []scan.File, workers int, hash func(string) (HashResult, error)) []scan.File {
	return ApplyHashesWithProgress(files, workers, hash, nil)
}

func ApplyHashesWithProgress(files []scan.File, workers int, hash func(string) (HashResult, error), progress func(processed int, total int, currentPath string)) []scan.File {
	return ApplyHashesContext(context.Background(), files, workers, hash, progress)
}

func ApplyHashesContext(ctx context.Context, files []scan.File, workers int, hash func(string) (HashResult, error), progress func(processed int, total int, currentPath string)) []scan.File {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	var processed atomic.Int64
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				if ctx.Err() != nil {
					return
				}
				if progress != nil {
					progress(int(processed.Load()), len(files), files[idx].AbsolutePath)
				}
				res, err := hash(files[idx].AbsolutePath)
				if err != nil {
					files[idx].IsAccessible = false
					files[idx].ScanError = err.Error()
				} else {
					files[idx].SHA256 = res.SHA256
					files[idx].IsUnstable = res.Unstable
				}
				next := int(processed.Add(1))
				if progress != nil {
					progress(next, len(files), files[idx].AbsolutePath)
				}
			}
		}()
	}
	for idx := range files {
		select {
		case jobs <- idx:
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return files
		}
	}
	close(jobs)
	wg.Wait()
	return files
}
