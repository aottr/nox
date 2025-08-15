package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/constants"
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

// IOWrapper is a generic wrapper for reading/writing files/STDIN/STDOUT
func IOWrapper[T any](input, output string, additional T, process func([]byte, T) ([]byte, error)) error {
	var inputBytes []byte
	var err error

	if input == constants.StandardInput {
		inputBytes, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		inputBytes, err = os.ReadFile(input)
		if err != nil {
			return err
		}
	}

	outbutBytes, err := process(inputBytes, additional)
	if err != nil {
		return err
	}

	if output == constants.StandardOutput {
		if _, err = os.Stdout.Write(outbutBytes); err != nil {
			return err
		}
	} else {
		if err = os.WriteFile(output, outbutBytes, 0600); err != nil {
			return err
		}
	}
	return nil
}
