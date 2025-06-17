package processor

import (
	"fmt"
	// "os"
	// "path/filepath"
	// "strings"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/gitrepo"
)

func ValidateConfig(cfg *config.Config) error {
	if cfg.AgeKeyPath == "" {
		return fmt.Errorf("age key path is required")
	}

	if cfg.StatePath == "" {
		fmt.Printf("state path is not set, defaulting to default.\n")
	}

	for appName, app := range cfg.Apps {
		fmt.Printf("✅ Validating app %s\n", appName)

		repoURL := app.Repo
		if repoURL == "" {
			repoURL = cfg.DefaultRepo
		}

		repo, err := gitrepo.CloneRepoInMemory(gitrepo.GitFetchOptions{
			RepoURL: repoURL,
			Branch:  app.Branch,
		})
		if err != nil {
			return fmt.Errorf("failed to clone for app %s: %w", appName, err)
		}

		for _, file := range app.Files {
			if !gitrepo.FileExistsInTree(repo.Tree, file.Path) {
				return fmt.Errorf("❌ file %s missing in app %s", file, appName)
			}
			fmt.Printf("✔️ Found file %s in repo\n", file)
		}
	}
	fmt.Println("✨ All checks passed!")
	return nil
}
