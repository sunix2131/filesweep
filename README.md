# FileSweep

FileSweep is a local desktop app for sorting out messy folders safely. It scans selected folders, finds exact duplicates, shows large files, groups storage by category, and lets the user prepare a file action plan before anything is moved or sent to trash.

The app is intentionally local-first: no accounts, no cloud sync, no telemetry, no background cleanup, and no automatic deletion.

## Screenshots

### Home

![FileSweep home screen](docs/screenshots/home.svg)

### Duplicates

![FileSweep duplicates screen](docs/screenshots/duplicates.svg)

### Categories

![FileSweep categories screen](docs/screenshots/categories.svg)

## Why

File managers are good at showing folders, but they do not answer the questions I usually have when a directory gets messy:

- which files are actually identical;
- what is taking the most space;
- what types of files dominate the folder;
- what can be moved safely;
- what will happen before I confirm an action.

FileSweep is built around one rule: the app should show the situation clearly, but the user stays in control.

## Current Features

- [x] Native desktop shell with Wails v2.
- [x] React + TypeScript frontend.
- [x] macOS Finder-inspired layout.
- [x] Native folder picker.
- [x] Recursive folder scanning.
- [x] Scan progress with current file path and scan phase.
- [x] Scan cancellation.
- [x] Local SQLite database.
- [x] Embedded SQL migrations.
- [x] Exact duplicate detection by file size and SHA-256.
- [x] Worker pool for hashing duplicate candidates.
- [x] Detection of files that changed while being hashed.
- [x] File categories by extension and MIME type.
- [x] Large files view with size filter and search.
- [x] Category summary view.
- [x] Action plan before file operations.
- [x] Safe move with size/hash checks.
- [x] Cross-volume copy, verify, then remove source.
- [x] Undo for move operations.
- [x] System trash integration without permanent delete fallback.
- [x] Action history.
- [x] CSV export.
- [x] Local settings.
- [x] Russian and English UI dictionaries.
- [x] Light, dark, and system theme modes.
- [x] Go unit/integration tests.
- [x] Frontend lint, typecheck, and Vitest.
- [x] GitHub Actions for PR checks and release builds.

## Roadmap

- [ ] Duplicate group details: open a group, compare all copies, choose which file to keep, and add selected copies to the action plan.
- [ ] Image previews for duplicate images and large image files.
- [ ] Better file category names and localized system errors.
- [ ] Row virtualization for very large scan results.
- [ ] Smarter filters for large files: folder, category, modified date, and custom size.
- [ ] Safer Windows trash integration through Shell APIs.
- [ ] Release artifacts for macOS, Windows, and Linux.
- [ ] Playwright smoke test for the full scan-to-action-plan flow.

## Stack

- Go
- Wails v2.12.0
- SQLite with `database/sql`
- `modernc.org/sqlite`
- React
- TypeScript
- Vite
- Zustand
- React Router
- Lucide Icons
- Tailwind CSS
- Vitest
- GitHub Actions

## Supported Platforms

Target platforms:

- macOS Intel / Apple Silicon
- Windows 10/11
- Linux x64

The app is currently developed and tested primarily on macOS.

## Development

Install Wails:

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
```

Install frontend dependencies:

```sh
npm install --prefix frontend
```

Run the app:

```sh
wails dev
```

## Checks

Backend:

```sh
gofmt -w .
go vet ./...
go test ./internal/... ./tests .
```

Frontend:

```sh
npm run lint --prefix frontend
npm run typecheck --prefix frontend
npm run test --prefix frontend
```

Build:

```sh
npm run build --prefix frontend
wails build
```

## Project Structure

```text
internal/
  application/      scan, actions, undo, export, settings
  domain/           scan, duplicates, actions, settings models
  infrastructure/   SQLite, filesystem, platform adapters
  transport/wails/  Wails API exposed to the frontend

frontend/
  src/app/          application shell
  src/pages/        screens
  src/components/   shared UI
  src/i18n/         translation dictionaries
```

## Privacy

FileSweep does not upload file lists, hashes, paths, or scan results. The database, settings, logs, thumbnails, and exports are stored locally in the user's application data directories.

The app does not perform irreversible deletion. Trash operations use the operating system trash and fail if system trash is unavailable.

## License

MIT
