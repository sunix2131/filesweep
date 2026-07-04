package settings

type Theme string
type Language string

const (
	ThemeSystem Theme    = "system"
	ThemeLight  Theme    = "light"
	ThemeDark   Theme    = "dark"
	LangRU      Language = "ru"
	LangEN      Language = "en"
)

type Settings struct {
	Language                Language `json:"language"`
	Theme                   Theme    `json:"theme"`
	IncludeHiddenFiles      bool     `json:"includeHiddenFiles"`
	FollowSymlinks          bool     `json:"followSymlinks"`
	MaxHashWorkers          int      `json:"maxHashWorkers"`
	MaxPreviewFileSizeMB    int      `json:"maxPreviewFileSizeMb"`
	ScanExcludedFolderNames []string `json:"scanExcludedFolderNames"`
	RecentFolders           []string `json:"recentFolders"`
}

func Default() Settings {
	return Settings{
		Language:                LangRU,
		Theme:                   ThemeSystem,
		MaxHashWorkers:          4,
		MaxPreviewFileSizeMB:    20,
		ScanExcludedFolderNames: []string{".git", "node_modules", "System Volume Information", "$RECYCLE.BIN"},
		RecentFolders:           []string{},
	}
}
