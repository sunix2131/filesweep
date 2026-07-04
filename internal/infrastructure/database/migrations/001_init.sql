CREATE TABLE IF NOT EXISTS scan_sessions (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    started_at TEXT,
    completed_at TEXT,
    cancelled_at TEXT,
    selected_paths_json TEXT NOT NULL,
    files_count INTEGER NOT NULL DEFAULT 0,
    total_size_bytes INTEGER NOT NULL DEFAULT 0,
    duplicate_groups_count INTEGER NOT NULL DEFAULT 0,
    reclaimable_bytes INTEGER NOT NULL DEFAULT 0,
    errors_count INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS scan_roots (
    id TEXT PRIMARY KEY,
    scan_session_id TEXT NOT NULL,
    root_path TEXT NOT NULL,
    root_name TEXT NOT NULL,
    files_count INTEGER NOT NULL DEFAULT 0,
    total_size_bytes INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (scan_session_id) REFERENCES scan_sessions(id)
);

CREATE TABLE IF NOT EXISTS scan_files (
    id TEXT PRIMARY KEY,
    scan_session_id TEXT NOT NULL,
    root_id TEXT NOT NULL,
    absolute_path TEXT NOT NULL,
    relative_path TEXT NOT NULL,
    name TEXT NOT NULL,
    extension TEXT,
    mime_type TEXT,
    category TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    modified_at TEXT,
    created_at_if_available TEXT,
    sha256 TEXT,
    is_hidden INTEGER NOT NULL DEFAULT 0,
    is_accessible INTEGER NOT NULL DEFAULT 1,
    is_symlink INTEGER NOT NULL DEFAULT 0,
    is_unstable INTEGER NOT NULL DEFAULT 0,
    scan_error TEXT,
    FOREIGN KEY (scan_session_id) REFERENCES scan_sessions(id),
    FOREIGN KEY (root_id) REFERENCES scan_roots(id)
);

CREATE INDEX IF NOT EXISTS idx_scan_files_session_size ON scan_files(scan_session_id, size_bytes DESC);
CREATE INDEX IF NOT EXISTS idx_scan_files_session_sha256 ON scan_files(scan_session_id, sha256);
CREATE INDEX IF NOT EXISTS idx_scan_files_category ON scan_files(scan_session_id, category);

CREATE TABLE IF NOT EXISTS duplicate_groups (
    id TEXT PRIMARY KEY,
    scan_session_id TEXT NOT NULL,
    sha256 TEXT NOT NULL,
    file_size_bytes INTEGER NOT NULL,
    files_count INTEGER NOT NULL,
    estimated_reclaimable_bytes INTEGER NOT NULL,
    recommended_file_id TEXT,
    category TEXT NOT NULL DEFAULT 'other',
    FOREIGN KEY (scan_session_id) REFERENCES scan_sessions(id)
);

CREATE TABLE IF NOT EXISTS duplicate_group_members (
    duplicate_group_id TEXT NOT NULL,
    scan_file_id TEXT NOT NULL,
    is_recommended_to_keep INTEGER NOT NULL DEFAULT 0,
    is_selected_for_action INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (duplicate_group_id, scan_file_id),
    FOREIGN KEY (duplicate_group_id) REFERENCES duplicate_groups(id),
    FOREIGN KEY (scan_file_id) REFERENCES scan_files(id)
);

CREATE TABLE IF NOT EXISTS file_actions (
    id TEXT PRIMARY KEY,
    action_type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TEXT NOT NULL,
    started_at TEXT,
    completed_at TEXT,
    undo_available INTEGER NOT NULL DEFAULT 0,
    undo_expires_at TEXT,
    summary TEXT NOT NULL,
    error_message TEXT
);

CREATE TABLE IF NOT EXISTS file_action_items (
    id TEXT PRIMARY KEY,
    action_id TEXT NOT NULL,
    source_path TEXT NOT NULL,
    target_path TEXT,
    source_size_bytes INTEGER NOT NULL,
    source_sha256 TEXT,
    status TEXT NOT NULL,
    error_message TEXT,
    FOREIGN KEY (action_id) REFERENCES file_actions(id)
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
