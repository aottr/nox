package processor

import (
	"fmt"
	// "os"
	// "path/filepath"
	// "strings"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/git"
)

func ValidateConfig(cfg *config.Config) error {
	if cfg.Age.Identity == "" {
		return fmt.Errorf("age key path is required")
	}

	if cfg.StatePath == "" {
		fmt.Printf("state path is not set, defaulting to default.\n")
	}

	for appName, app := range cfg.Apps {
		fmt.Printf("✅ Validating app %s\n", appName)

		gitConf := app.GitConfig
		if !gitConf.IsValid() {
			gitConf = cfg.GitConfig
		}

		repo, err := git.CloneRepo(gitConf)
		if err != nil {
			return fmt.Errorf("failed to clone for app %s: %w", appName, err)
		}

		for _, file := range app.Files {
			if !git.FileExistsInTree(repo.Tree, file.Path) {
				return fmt.Errorf("❌ file %s missing in app %s", file, appName)
			}
			fmt.Printf("✔️ Found file %s in repo\n", file)
		}
	}
	fmt.Println("all checks passed!")
	return nil
}
