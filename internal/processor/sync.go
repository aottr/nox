package processor

import (
	"fmt"
	"os"

	"github.com/aottr/nox/internal/cache"
	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/crypto"
	"github.com/aottr/nox/internal/git"
	"github.com/aottr/nox/internal/logging"
	"github.com/aottr/nox/internal/state"
)

func SyncApp(ctx *config.RuntimeContext) error {

	log := logging.Get()
	var err error
	cfg, appName, identities, st := ctx.Config, ctx.App, ctx.Identities, ctx.State

	if appName == "" {
		return fmt.Errorf("app name is required")
	}

	// retrieve app config and repository
	app := cfg.Apps[appName]
	gitConf := app.GitConfig
	if !gitConf.IsValid() {
		gitConf = cfg.GitConfig
	}

	key := cache.RepoKey{Repo: gitConf.Repo, Branch: gitConf.Branch}
	repo, err := cache.GlobalCache.GetOrFetch(key)
	if err != nil {
		return fmt.Errorf("failed to fetch repo for app %s: %w", appName, err)
	}

	// iterate over files and decrypt
	for _, file := range app.Files {
		content, err := git.GetFileContentFromTree(repo.Tree, file.Path)
		if err != nil {
			return fmt.Errorf("failed to get file %s: %w", file, err)
		}

		hash := state.HashContent(content)
		cacheKey := state.GenerateKey(appName, file.Path)

		// skip if file is up to date and force is not set
		if !ctx.Force && !ctx.DryRun {
			if prevHash, ok := st.Data[cacheKey]; ok && prevHash == hash {
				log.Debug(fmt.Sprintf("file %s is up to date", file.Path))
				continue
			}
		}

		// decrypt file
		plaintext, err := crypto.DecryptBytes(content, identities)
		if err != nil {
			log.Warn("failed to decrypt file %s: %v", file.Path, err)
			continue
		}

		// skip writing file if dry run is set
		if ctx.DryRun {
			log.Debug(fmt.Sprintf("dry run, not writing file %s", file.Output))
			os.Stdout.Write(plaintext)
			continue
		}
		if err := WriteToFile(plaintext, file); err != nil {
			log.Error("failed to write file %s: %v", file.Output, err)
			continue
		}

		log.Debug(fmt.Sprintf("decrypted %s for app %s (size: %d bytes)", file, appName, len(plaintext)))

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
		ctx.App = appName
		logging.Get().Debug(fmt.Sprintf("Processing app: %s", appName))
		if err := SyncApp(ctx); err != nil {
			return err
		}
	}
	return nil
}
