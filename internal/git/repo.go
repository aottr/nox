package git

import (
	"fmt"
	"io"

	"github.com/aottr/nox/internal/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

type ClonedRepo struct {
	Repo   *git.Repository
	Branch string
	Tree   *object.Tree
	Ref    *plumbing.Reference
	Commit *object.Commit
}

func CloneRepo(c config.GitConfig) (*ClonedRepo, error) {

	auth, err := GetAuth()
	if err != nil {
		return nil, err
	}

	cloneOpts := &git.CloneOptions{
		URL:           c.Repo,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(c.Branch),
		Auth:          auth,
	}

	repo, err := git.Clone(memory.NewStorage(), nil, cloneOpts)
	if err != nil {
		return nil, fmt.Errorf("clone failed: %w", err)
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}
	return &ClonedRepo{
		Repo:   repo,
		Branch: c.Branch,
		Ref:    ref,
		Commit: commit,
		Tree:   tree,
	}, nil
}

func GetFileContentFromTree(tree *object.Tree, path string) ([]byte, error) {
	file, err := tree.File(path)
	if err != nil {
		return nil, fmt.Errorf("file %q not found: %w", path, err)
	}

	reader, err := file.Blob.Reader()
	if err != nil {
		return nil, fmt.Errorf("failed to open reader for %q: %w", path, err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content of %q: %w", path, err)
	}

	return content, nil
}

func (r *ClonedRepo) Refresh() error {
	if err := r.Repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		Prune:      true,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("fetch: %w", err)
	}
	remoteRef := plumbing.NewRemoteReferenceName("origin", r.Branch)
	ref, err := r.Repo.Reference(remoteRef, true)
	if err != nil {
		return err
	}
	r.Ref = ref

	commit, err := r.Repo.CommitObject(r.Ref.Hash())
	if err != nil {
		return err
	}
	r.Commit = commit
	tree, err := commit.Tree()
	if err != nil {
		return err
	}
	r.Tree = tree
	return nil
}

func FileExistsInTree(tree *object.Tree, path string) bool {
	_, err := tree.File(path)
	return err == nil
}
