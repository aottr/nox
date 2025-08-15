package cache

import (
	"sync"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/git"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type RepoKey struct {
	Repo   string
	Branch string
}

type RepoCache struct {
	mu    sync.RWMutex
	repos map[RepoKey]*object.Tree
}

var (
	GlobalCache = &RepoCache{
		repos: make(map[RepoKey]*object.Tree),
	}
)

func (c *RepoCache) Get(key RepoKey) (*object.Tree, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tree, exists := c.repos[key]
	return tree, exists
}

func (c *RepoCache) Set(key RepoKey, tree *object.Tree) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.repos[key] = tree
}

func (c *RepoCache) FetchRepo(key RepoKey) (*object.Tree, error) {
	r, err := git.CloneRepo(config.GitConfig{
		Repo:   key.Repo,
		Branch: key.Branch,
	})
	if err != nil {
		return nil, err
	}

	c.Set(key, r.Tree)
	return r.Tree, nil
}

func ClearRepoCache() {
	GlobalCache.mu.Lock()
	defer GlobalCache.mu.Unlock()
	GlobalCache.repos = make(map[RepoKey]*object.Tree)
}
