# FileSweep

FileSweep scans local folders and shows exact duplicates, large files and storage use by category. File operations are collected into a plan and require confirmation in a native dialog; scanning never changes the selected folders. Scanning a filesystem root or home directory requires a separate warning confirmation.

The application is built with Go and Wails. Scan results, settings and action history are stored in a local SQLite database.

## Current behavior

- recursive scanning with exclusions, progress and cancellation;
- exact duplicate detection using file size followed by SHA-256;
- large-file and category views;
- a reviewable action plan for move and trash operations;
- undo for completed moves;
- CSV export of scan results;
- Russian and English interface;
- light, dark and system themes.

Duplicate candidates are grouped by size before hashing. A file whose size or modification time changes during hashing is marked unstable and excluded from duplicate groups.

Overlapping selected folders count each path once. Cancellation is checked during both discovery and hashing, between read buffers. Symbolic links are skipped; FileSweep does not follow them outside the selected folders.

## File operation checks

FileSweep does not contain a permanent-delete fallback. A trash operation fails if the operating system trash service is unavailable.

Before a move, the application compares the current file with the recorded size and, when available, SHA-256. Existing destination files are not overwritten. FileSweep first attempts a no-overwrite hard link followed by removal of the source. When linking the source is unavailable, the file is copied to a unique temporary file, verified and linked into the destination before the source is removed.

If removing the source fails, the new destination entry is rolled back. The original file remains in place.

A destination name conflict fails that item; choose another target in the plan. No hidden auto-renaming is performed, so move history records the path that was actually used. Batch operations are not transactional: if one item fails, earlier successful moves remain in place. Undo is available only for a fully successful move batch and will fail if a source or destination has changed.

## Development

Requirements:

- Go 1.26.8 or newer;
- Node.js 22 or newer;
- Wails CLI 2.12;
- platform packages required by Wails.

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
npm ci --prefix frontend
wails dev
```

## Checks

The frontend must be built first because `main.go` embeds `frontend/dist`.

```bash
npm ci --prefix frontend
npm run lint --prefix frontend
npm run typecheck --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend

test -z "$(gofmt -l . | grep -v '^frontend/')"
go vet ./internal/... ./tests .
go test -race ./internal/... ./tests .
```

GitHub Actions also runs `staticcheck`, `govulncheck` and a Wails build on Linux. Tagged releases are built separately on macOS, Windows and Linux runners.

## Layout

```text
internal/
  application/      scan and action workflows
  domain/           scan, duplicate, action and settings types
  infrastructure/   SQLite, filesystem and platform adapters
  transport/wails/  API exposed to the frontend

frontend/src/
  app/              application shell
  pages/            scan result screens
  components/       shared interface components
  i18n/             Russian and English strings
```

## Local data

File paths, hashes and scan results are not sent to an application server. Thumbnails, exports, settings and logs remain in the operating system's application-data directory.

The code is licensed under the MIT License. See [LICENSE](LICENSE).
