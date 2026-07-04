package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"filesweep/internal/domain/actions"
	"filesweep/internal/domain/duplicates"
	scan "filesweep/internal/domain/scan"
	"filesweep/internal/domain/settings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) SaveScan(ctx context.Context, session scan.Session, roots []scan.Root, files []scan.File, groups []duplicates.Group) error {
	paths, _ := json.Marshal(session.SelectedPaths)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO scan_sessions
		(id,status,started_at,completed_at,cancelled_at,selected_paths_json,files_count,total_size_bytes,duplicate_groups_count,reclaimable_bytes,errors_count,created_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		session.ID, session.Status, ts(session.StartedAt), ts(session.CompletedAt), ts(session.CancelledAt), string(paths),
		session.FilesCount, session.TotalSizeBytes, session.DuplicateGroupsCount, session.ReclaimableBytes, session.ErrorsCount, ts(session.CreatedAt))
	if err != nil {
		return err
	}
	for _, r := range roots {
		_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO scan_roots (id,scan_session_id,root_path,root_name,files_count,total_size_bytes) VALUES (?,?,?,?,?,?)`,
			r.ID, r.ScanSessionID, r.RootPath, r.RootName, r.FilesCount, r.TotalSizeBytes)
		if err != nil {
			return err
		}
	}
	stmt, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO scan_files
		(id,scan_session_id,root_id,absolute_path,relative_path,name,extension,mime_type,category,size_bytes,modified_at,created_at_if_available,sha256,is_hidden,is_accessible,is_symlink,is_unstable,scan_error)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, f := range files {
		_, err = stmt.ExecContext(ctx, f.ID, f.ScanSessionID, f.RootID, f.AbsolutePath, f.RelativePath, f.Name, f.Extension, f.MimeType, f.Category,
			f.SizeBytes, ts(f.ModifiedAt), ts(f.CreatedAtIfAvailable), f.SHA256, b(f.IsHidden), b(f.IsAccessible), b(f.IsSymlink), b(f.IsUnstable), f.ScanError)
		if err != nil {
			return err
		}
	}
	for _, g := range groups {
		_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO duplicate_groups (id,scan_session_id,sha256,file_size_bytes,files_count,estimated_reclaimable_bytes,recommended_file_id,category) VALUES (?,?,?,?,?,?,?,?)`,
			g.ID, g.ScanSessionID, g.SHA256, g.FileSizeBytes, g.FilesCount, g.EstimatedReclaimableBytes, g.RecommendedFileID, g.Category)
		if err != nil {
			return err
		}
		for _, m := range g.Members {
			_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO duplicate_group_members (duplicate_group_id,scan_file_id,is_recommended_to_keep,is_selected_for_action) VALUES (?,?,?,0)`,
				g.ID, m.ID, b(m.ID == g.RecommendedFileID))
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Store) UpdateScanStatus(ctx context.Context, id string, status scan.Status) error {
	col := "completed_at"
	if status == scan.StatusCancelled {
		col = "cancelled_at"
	}
	_, err := s.db.ExecContext(ctx, `UPDATE scan_sessions SET status=?, `+col+`=? WHERE id=?`, status, ts(time.Now()), id)
	return err
}

func (s *Store) ListScanSessions(ctx context.Context) ([]scan.Session, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,status,started_at,completed_at,cancelled_at,selected_paths_json,files_count,total_size_bytes,duplicate_groups_count,reclaimable_bytes,errors_count,created_at FROM scan_sessions ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []scan.Session
	for rows.Next() {
		var ss scan.Session
		var selected string
		var started, completed, cancelled, created string
		if err := rows.Scan(&ss.ID, &ss.Status, &started, &completed, &cancelled, &selected, &ss.FilesCount, &ss.TotalSizeBytes, &ss.DuplicateGroupsCount, &ss.ReclaimableBytes, &ss.ErrorsCount, &created); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(selected), &ss.SelectedPaths)
		ss.StartedAt, ss.CompletedAt, ss.CancelledAt, ss.CreatedAt = parseTS(started), parseTS(completed), parseTS(cancelled), parseTS(created)
		out = append(out, ss)
	}
	return out, rows.Err()
}

func (s *Store) GetScanSession(ctx context.Context, id string) (scan.Session, error) {
	list, err := s.ListScanSessions(ctx)
	if err != nil {
		return scan.Session{}, err
	}
	for _, item := range list {
		if item.ID == id {
			return item, nil
		}
	}
	return scan.Session{}, sql.ErrNoRows
}

func (s *Store) ListFiles(ctx context.Context, scanID string, category string, minSize int64, search string, limit, offset int) ([]scan.File, int, error) {
	where := `WHERE scan_session_id=?`
	args := []interface{}{scanID}
	if category != "" {
		where += ` AND category=?`
		args = append(args, category)
	}
	if minSize > 0 {
		where += ` AND size_bytes>=?`
		args = append(args, minSize)
	}
	if search != "" {
		where += ` AND lower(name) LIKE lower(?)`
		args = append(args, "%"+search+"%")
	}
	var total int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM scan_files `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, limit, offset)
	rows, err := s.db.QueryContext(ctx, `SELECT id,scan_session_id,root_id,absolute_path,relative_path,name,extension,mime_type,category,size_bytes,modified_at,sha256,is_hidden,is_accessible,is_symlink,is_unstable,scan_error FROM scan_files `+where+` ORDER BY size_bytes DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	files, err := scanRows(rows)
	return files, total, err
}

func (s *Store) ListDuplicateGroups(ctx context.Context, scanID string, limit, offset int) ([]duplicates.Group, int, error) {
	var total int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM duplicate_groups WHERE scan_session_id=?`, scanID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,scan_session_id,sha256,file_size_bytes,files_count,estimated_reclaimable_bytes,recommended_file_id,category FROM duplicate_groups WHERE scan_session_id=? ORDER BY estimated_reclaimable_bytes DESC LIMIT ? OFFSET ?`, scanID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []duplicates.Group
	for rows.Next() {
		var g duplicates.Group
		if err := rows.Scan(&g.ID, &g.ScanSessionID, &g.SHA256, &g.FileSizeBytes, &g.FilesCount, &g.EstimatedReclaimableBytes, &g.RecommendedFileID, &g.Category); err != nil {
			return nil, 0, err
		}
		out = append(out, g)
	}
	return out, total, rows.Err()
}

func (s *Store) GetDuplicateGroup(ctx context.Context, id string) (duplicates.Group, error) {
	var g duplicates.Group
	err := s.db.QueryRowContext(ctx, `SELECT id,scan_session_id,sha256,file_size_bytes,files_count,estimated_reclaimable_bytes,recommended_file_id,category FROM duplicate_groups WHERE id=?`, id).
		Scan(&g.ID, &g.ScanSessionID, &g.SHA256, &g.FileSizeBytes, &g.FilesCount, &g.EstimatedReclaimableBytes, &g.RecommendedFileID, &g.Category)
	if err != nil {
		return g, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT f.id,f.scan_session_id,f.root_id,f.absolute_path,f.relative_path,f.name,f.extension,f.mime_type,f.category,f.size_bytes,f.modified_at,f.sha256,f.is_hidden,f.is_accessible,f.is_symlink,f.is_unstable,f.scan_error
		FROM scan_files f JOIN duplicate_group_members m ON m.scan_file_id=f.id WHERE m.duplicate_group_id=? ORDER BY f.absolute_path`, id)
	if err != nil {
		return g, err
	}
	defer rows.Close()
	g.Members, err = scanRows(rows)
	return g, err
}

func (s *Store) CategorySummary(ctx context.Context, scanID string) ([]map[string]interface{}, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT category, COUNT(*), COALESCE(SUM(size_bytes),0) FROM scan_files WHERE scan_session_id=? GROUP BY category ORDER BY SUM(size_bytes) DESC`, scanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var c string
		var n int
		var size int64
		if err := rows.Scan(&c, &n, &size); err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{"category": c, "filesCount": n, "totalSizeBytes": size})
	}
	return out, rows.Err()
}

func (s *Store) SaveAction(ctx context.Context, a actions.Action) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO file_actions (id,action_type,status,created_at,started_at,completed_at,undo_available,undo_expires_at,summary,error_message) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		a.ID, a.ActionType, a.Status, ts(a.CreatedAt), ts(a.StartedAt), ts(a.CompletedAt), b(a.UndoAvailable), ts(a.UndoExpiresAt), a.Summary, a.ErrorMessage)
	if err != nil {
		return err
	}
	for _, item := range a.Items {
		_, err = tx.ExecContext(ctx, `INSERT OR REPLACE INTO file_action_items (id,action_id,source_path,target_path,source_size_bytes,source_sha256,status,error_message) VALUES (?,?,?,?,?,?,?,?)`,
			item.ID, a.ID, item.SourcePath, item.TargetPath, item.SourceSizeBytes, item.SourceSHA256, item.Status, item.ErrorMessage)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListActions(ctx context.Context, limit, offset int) ([]actions.Action, int, error) {
	var total int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM file_actions`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,action_type,status,created_at,started_at,completed_at,undo_available,undo_expires_at,summary,error_message FROM file_actions ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []actions.Action
	for rows.Next() {
		var a actions.Action
		var created, started, completed, undo string
		var undoAvailable int
		if err := rows.Scan(&a.ID, &a.ActionType, &a.Status, &created, &started, &completed, &undoAvailable, &undo, &a.Summary, &a.ErrorMessage); err != nil {
			return nil, 0, err
		}
		a.CreatedAt, a.StartedAt, a.CompletedAt, a.UndoExpiresAt = parseTS(created), parseTS(started), parseTS(completed), parseTS(undo)
		a.UndoAvailable = undoAvailable == 1
		out = append(out, a)
	}
	return out, total, rows.Err()
}

func (s *Store) GetAction(ctx context.Context, id string) (actions.Action, error) {
	actionsList, _, err := s.ListActions(ctx, 1000, 0)
	if err != nil {
		return actions.Action{}, err
	}
	for _, a := range actionsList {
		if a.ID == id {
			rows, err := s.db.QueryContext(ctx, `SELECT id,action_id,source_path,target_path,source_size_bytes,source_sha256,status,error_message FROM file_action_items WHERE action_id=?`, id)
			if err != nil {
				return a, err
			}
			defer rows.Close()
			for rows.Next() {
				var item actions.Item
				if err := rows.Scan(&item.ID, &item.ActionID, &item.SourcePath, &item.TargetPath, &item.SourceSizeBytes, &item.SourceSHA256, &item.Status, &item.ErrorMessage); err != nil {
					return a, err
				}
				a.Items = append(a.Items, item)
			}
			return a, rows.Err()
		}
	}
	return actions.Action{}, sql.ErrNoRows
}

func (s *Store) GetSettings(ctx context.Context) (settings.Settings, error) {
	cfg := settings.Default()
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key='settings'`).Scan(&raw)
	if err == sql.ErrNoRows {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return settings.Default(), err
	}
	return cfg, nil
}

func (s *Store) SaveSettings(ctx context.Context, cfg settings.Settings) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT OR REPLACE INTO settings (key,value) VALUES ('settings',?)`, string(raw))
	return err
}

func scanRows(rows *sql.Rows) ([]scan.File, error) {
	var out []scan.File
	for rows.Next() {
		var f scan.File
		var modified string
		var hidden, accessible, symlink, unstable int
		if err := rows.Scan(&f.ID, &f.ScanSessionID, &f.RootID, &f.AbsolutePath, &f.RelativePath, &f.Name, &f.Extension, &f.MimeType, &f.Category, &f.SizeBytes, &modified, &f.SHA256, &hidden, &accessible, &symlink, &unstable, &f.ScanError); err != nil {
			return nil, err
		}
		f.ModifiedAt = parseTS(modified)
		f.IsHidden = hidden == 1
		f.IsAccessible = accessible == 1
		f.IsSymlink = symlink == 1
		f.IsUnstable = unstable == 1
		out = append(out, f)
	}
	return out, rows.Err()
}

func ts(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func parseTS(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339Nano, v)
	return t
}

func b(v bool) int {
	if v {
		return 1
	}
	return 0
}
