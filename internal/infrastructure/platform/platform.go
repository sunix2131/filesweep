package platform

type TrashService interface {
	MoveToTrash(paths []string) error
}

type FileRevealService interface {
	RevealInFileManager(path string) error
}

type FileOpenService interface {
	OpenFile(path string) error
}

type AppPathsService interface {
	GetAppDataDir() (string, error)
	GetCacheDir() (string, error)
	GetLogsDir() (string, error)
}

type Services struct{}
