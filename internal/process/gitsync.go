package process

import (
	"fmt"
	"log"
	"os"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/crypto"
	"github.com/aottr/nox/internal/gitrepo"
	"github.com/aottr/nox/internal/state"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ProcessApps clones and decrypts configured app secrets efficiently.
func ProcessApps(cfg *config.Config) error {
	type repoKey struct {
		Repo   string
		Branch string
	}

	clones := map[repoKey]*object.Tree{}

	identities, err := crypto.LoadAgeIdentities(cfg.AgeKeyPath)
	if err != nil {
		return fmt.Errorf("Failed to load age identities: %w", err)
	}

	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("Failed to load state: %w", err)
	}

	for appName, app := range cfg.Apps {

		fmt.Printf("Processing app %s\n", appName)

		repoUrl := app.Repo
		if repoUrl == "" {
			repoUrl = cfg.DefaultRepo
		}
		key := repoKey{Repo: repoUrl, Branch: app.Branch}

		clone, ok := clones[key]
		if !ok {
			repo, err := gitrepo.CloneRepoInMemory(gitrepo.GitFetchOptions{
				RepoURL: repoUrl,
				Branch:  app.Branch,
			})
			if err != nil {
				log.Printf("Clone failed for %s/%s: %v", repoUrl, app.Branch, err)
				continue
			}
			clone = repo.Tree
			clones[key] = clone
		}

		for _, file := range app.Files {
			content, err := gitrepo.GetFileContentFromTree(clone, file.Path)
			if err != nil {
				log.Printf("Failed to get file %s: %v", file, err)
				continue
			}

			hash := state.HashContent(content)
			cacheKey := state.GenerateKey(appName, file.Path)

			if prevHash, ok := st.Data[cacheKey]; ok && prevHash == hash {
				log.Printf("File %s is up to date", file.Path)
				continue
			}

			plaintext, err := crypto.DecryptBytes(content, identities)
			if err != nil {
				log.Printf("Failed to decrypt file %s: %v", file.Path, err)
				continue
			}

			outPath := file.Output
			// if outPath == "" {
			// 	// Default output filename if none specified, e.g. replace .age with .env
			// 	outPath = filepath.Base(fileCfg.Path)
			// 	if filepath.Ext(outPath) == ".age" {
			// 		outPath = outPath[:len(outPath)-4] + ".env"
			// 	}
			// }

			// if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			// 	log.Printf("Failed to create directories for %s: %v", outPath, err)
			// 	continue
			// }

			if err := os.WriteFile(outPath, plaintext, 0600); err != nil {
				log.Printf("Failed to write decrypted file to %s: %v", outPath, err)
				continue
			}

			log.Printf("decrypted %s for app %s (size: %d bytes)", file, appName, len(plaintext))

			st.Data[cacheKey] = hash
			st.Touch()

			log.Printf("Decrypted file %s", file)
		}

	}

	if err := state.Save(st); err != nil {
		return fmt.Errorf("Failed to save state: %w", err)
	}

	return nil
}
