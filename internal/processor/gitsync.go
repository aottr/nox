package processor

import (
	"fmt"
	"os"

	"github.com/aottr/nox/internal/cache"
	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/crypto"
	"github.com/aottr/nox/internal/git"
	"github.com/aottr/nox/internal/state"
)

func SyncApp(ctx *config.RuntimeContext) error {

	cfg, appName, identities, st := ctx.Config, ctx.App, ctx.Identities, ctx.State

	if appName == nil {
		return fmt.Errorf("app name is required")
	}

	// retrieve app config and repository
	app := cfg.Apps[*appName]
	repoUrl := app.GitConfig.Repo
	if repoUrl == "" {
		repoUrl = cfg.GitConfig.Repo
	}
	branchName := app.GitConfig.Branch
	if branchName == "" {
		branchName = cfg.GitConfig.Branch
	}

	key := cache.RepoKey{Repo: repoUrl, Branch: branchName}
	repo, exists := cache.GlobalCache.Get(key)
	if !exists {
		var err error
		repo, err = cache.GlobalCache.FetchRepo(key)
		if err != nil {
			return fmt.Errorf("failed to fetch repo for app %s: %w", *appName, err)
		}
	}

	// iterate over files and decrypt
	for _, file := range app.Files {
		content, err := git.GetFileContentFromTree(repo, file.Path)
		if err != nil {
			return fmt.Errorf("failed to get file %s: %w", file, err)
		}

		hash := state.HashContent(content)
		cacheKey := state.GenerateKey(*appName, file.Path)

		// skip if file is up to date and force is not set
		if !ctx.Force && !ctx.DryRun {
			if prevHash, ok := st.Data[cacheKey]; ok && prevHash == hash {
				ctx.Logger.Printf("file %s is up to date", file.Path)
				continue
			}
		}

		// decrypt file
		plaintext, err := crypto.DecryptBytes(content, identities)
		if err != nil {
			ctx.Logger.Printf("failed to decrypt file %s: %v", file.Path, err)
			continue
		}

		// skip writing file if dry run is set
		if ctx.DryRun {
			ctx.Logger.Printf("dry run, not writing file %s", file.Output)
			os.Stdout.Write(plaintext)
			continue
		}
		WriteToFile(plaintext, file, &FileProcessorOptions{CreateDir: true})

		ctx.Logger.Printf("decrypted %s for app %s (size: %d bytes)", file, *appName, len(plaintext))

		// update state
		st.Data[cacheKey] = hash
		st.Touch()
	}

	if err := state.Save(st); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}
	return nil
}

func SyncApps(ctx *config.RuntimeContext) error {
	for appName := range ctx.Config.Apps {
		ctx.App = &appName
		ctx.Logger.Printf("Processing app: %s\n", appName)
		if err := SyncApp(ctx); err != nil {
			return err
		}
	}
	return nil
}
