package filesystem

import (
	"mime"
	"path/filepath"
	"strings"

	scan "filesweep/internal/domain/scan"
)

var extCategories = map[string]string{
	"jpg": scan.CategoryImages, "jpeg": scan.CategoryImages, "png": scan.CategoryImages, "gif": scan.CategoryImages, "webp": scan.CategoryImages, "bmp": scan.CategoryImages, "tiff": scan.CategoryImages, "heic": scan.CategoryImages,
	"mp4": scan.CategoryVideos, "mov": scan.CategoryVideos, "mkv": scan.CategoryVideos, "avi": scan.CategoryVideos, "webm": scan.CategoryVideos, "m4v": scan.CategoryVideos,
	"mp3": scan.CategoryAudio, "wav": scan.CategoryAudio, "flac": scan.CategoryAudio, "aac": scan.CategoryAudio, "m4a": scan.CategoryAudio, "ogg": scan.CategoryAudio,
	"pdf": scan.CategoryDocuments, "doc": scan.CategoryDocuments, "docx": scan.CategoryDocuments, "rtf": scan.CategoryDocuments, "odt": scan.CategoryDocuments, "txt": scan.CategoryDocuments, "md": scan.CategoryDocuments,
	"xls": scan.CategorySpreadsheets, "xlsx": scan.CategorySpreadsheets, "csv": scan.CategorySpreadsheets, "ods": scan.CategorySpreadsheets,
	"ppt": scan.CategoryPresentations, "pptx": scan.CategoryPresentations, "odp": scan.CategoryPresentations,
	"zip": scan.CategoryArchives, "rar": scan.CategoryArchives, "7z": scan.CategoryArchives, "tar": scan.CategoryArchives, "gz": scan.CategoryArchives, "bz2": scan.CategoryArchives, "xz": scan.CategoryArchives,
	"exe": scan.CategoryInstallers, "msi": scan.CategoryInstallers, "dmg": scan.CategoryInstallers, "pkg": scan.CategoryInstallers, "appimage": scan.CategoryInstallers, "deb": scan.CategoryInstallers, "rpm": scan.CategoryInstallers, "apk": scan.CategoryInstallers,
	"go": scan.CategoryCode, "py": scan.CategoryCode, "java": scan.CategoryCode, "js": scan.CategoryCode, "ts": scan.CategoryCode, "jsx": scan.CategoryCode, "tsx": scan.CategoryCode, "html": scan.CategoryCode, "css": scan.CategoryCode, "scss": scan.CategoryCode, "json": scan.CategoryCode, "yaml": scan.CategoryCode, "yml": scan.CategoryCode, "xml": scan.CategoryCode, "sql": scan.CategoryCode, "sh": scan.CategoryCode,
	"ttf": scan.CategoryFonts, "otf": scan.CategoryFonts, "woff": scan.CategoryFonts, "woff2": scan.CategoryFonts,
	"iso": scan.CategoryDiskImages, "img": scan.CategoryDiskImages, "vhd": scan.CategoryDiskImages, "vhdx": scan.CategoryDiskImages,
}

func CategoryFor(path string) (string, string) {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	mimeType := mime.TypeByExtension("." + ext)
	if strings.HasPrefix(mimeType, "image/") {
		return scan.CategoryImages, mimeType
	}
	if strings.HasPrefix(mimeType, "video/") {
		return scan.CategoryVideos, mimeType
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return scan.CategoryAudio, mimeType
	}
	if c, ok := extCategories[ext]; ok {
		return c, mimeType
	}
	return scan.CategoryOther, mimeType
}

func IsHidden(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".") && base != "." && base != ".."
}
