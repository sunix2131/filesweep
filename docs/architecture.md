# FileSweep Architecture

FileSweep is a local-only Wails v2 desktop application. The frontend calls a typed Go API; Go owns scanning, SQLite persistence, safe file actions, history, settings and platform integration.

Dependency direction:

`frontend -> transport/wails -> application -> domain`

Infrastructure packages provide SQLite, filesystem and platform adapters. Domain types do not import Wails, React, SQLite or OS-specific packages.

Data is stored under the current user's system config/cache directories:

- `filesweep.db`
- `settings`
- `logs/`
- `thumbnails/`
- `exports/`
- `backups/`

The app never performs permanent delete. Trash operations use the system trash and fail if it is unavailable.
