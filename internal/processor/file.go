package processor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aottr/nox/internal/config"
)

type FileProcessorOptions struct {
	CreateDir bool
}

func WriteToFile(data []byte, file config.FileConfig, opts *FileProcessorOptions) error {
	path := file.Output
	if path == "" {
		// Default output filename if none specified, e.g. replace .age with .env
		path = filepath.Base(file.Path)
		fmt.Println(path)
		if filepath.Ext(path) == ".age" {
			path = path[:len(path)-4] + ".env"
		}
	}
	if opts.CreateDir {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create directories for %s: %w", path, err)
		}
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write decrypted file to %s: %w", path, err)
	}
	return nil
}
