package cache

import (
	"sync"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/git"
)

type RepoKey struct {
	Repo   string
	Branch string
}

type RepoCache struct {
	mu    sync.RWMutex
	repos map[RepoKey]*git.ClonedRepo
	// sf    singleflight.Group TODO
}

var (
	GlobalCache = &RepoCache{
		repos: make(map[RepoKey]*git.ClonedRepo),
	}
)

func (c *RepoCache) Get(key RepoKey) (*git.ClonedRepo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tree, exists := c.repos[key]
	return tree, exists
}

func (c *RepoCache) Set(key RepoKey, repo *git.ClonedRepo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.repos[key] = repo
}

func (c *RepoCache) FetchRepo(key RepoKey) (*git.ClonedRepo, error) {
	r, err := git.CloneRepo(config.GitConfig{
		Repo:   key.Repo,
		Branch: key.Branch,
	})
	if err != nil {
		return nil, err
	}

	c.Set(key, r)
	return r, nil
}

func (c *RepoCache) GetOrFetch(key RepoKey) (*git.ClonedRepo, error) {
	tree, exists := c.Get(key)
	if exists {
		return tree, nil
	} else {
		return c.FetchRepo(key)
	}
}

func (c *RepoCache) RefreshCache() error {
	c.mu.RLock()
	repos := make([]*git.ClonedRepo, 0, len(c.repos))
	for _, r := range c.repos {
		repos = append(repos, r)
	}
	c.mu.RUnlock()

	for _, repo := range repos {
		if err := repo.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

func (c *RepoCache) Has(key RepoKey) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.repos[key]
	return exists
}

func ClearRepoCache() {
	GlobalCache.mu.Lock()
	defer GlobalCache.mu.Unlock()
	GlobalCache.repos = make(map[RepoKey]*git.ClonedRepo)
}
