package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func ThumbnailPath(cacheDir, sourcePath string) string {
	sum := sha256.Sum256([]byte(sourcePath))
	return filepath.Join(cacheDir, hex.EncodeToString(sum[:])+".jpg")
}

func CreateThumbnail(cacheDir, sourcePath string, maxSourceBytes int64) (string, error) {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", err
	}
	if info.Size() > maxSourceBytes {
		return "", fmt.Errorf("image exceeds preview size limit")
	}
	target := ThumbnailPath(cacheDir, sourcePath)
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}
	if strings.EqualFold(filepath.Ext(sourcePath), ".webp") {
		return "", fmt.Errorf("webp preview decoder is unavailable")
	}
	in, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer in.Close()
	img, _, err := image.Decode(in)
	if err != nil {
		return "", err
	}
	thumb := resizeNearest(img, 320)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}
	out, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if err := jpeg.Encode(out, thumb, &jpeg.Options{Quality: 84}); err != nil {
		return "", err
	}
	return target, nil
}

func resizeNearest(src image.Image, maxSide int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxSide && h <= maxSide {
		return src
	}
	if w >= h {
		h = h * maxSide / w
		w = maxSide
	} else {
		w = w * maxSide / h
		h = maxSide
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX := b.Min.X + x*b.Dx()/w
			srcY := b.Min.Y + y*b.Dy()/h
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}
