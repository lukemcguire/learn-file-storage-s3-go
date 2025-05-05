package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	randID := make([]byte, 32)
	rand.Read(randID)
	urlSafeID := base64.RawURLEncoding.EncodeToString(randID)
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", urlSafeID, ext)
}

func getAssetKey(mediaType string) string {
	randID := make([]byte, 32)
	rand.Read(randID)
	fileKey := hex.EncodeToString(randID)
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", fileKey, ext)
}

func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func verifyMP4(file *os.File) (bool, error) {
	buffer := make([]byte, 12)
	n, err := io.ReadFull(file, buffer)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return false, fmt.Errorf("read %d bytes, but reached end of file", n)
		} else {
			return false, fmt.Errorf("error reading from file: %w", err)
		}
	}
	if string(buffer[4:8]) != "ftyp" {
		return false, errors.New("filetype is not mp4")
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return false, errors.New("unable to reset file read pointer")
	}

	return true, nil
}
