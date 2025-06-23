package gitrepo

import (
	"fmt"
	"io"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gitMem "github.com/go-git/go-git/v5/storage/memory"
)

type GitFetchOptions struct {
	RepoURL string
	Branch  string
	File    string
	Token   *string // optional
}

type ClonedRepo struct {
	Repo   *git.Repository
	Tree   *object.Tree
	Ref    *plumbing.Reference
	Commit *object.Commit
}

func CloneRepoInMemory(opts GitFetchOptions) (*ClonedRepo, error) {

	cloneOpts := &git.CloneOptions{
		URL:           opts.RepoURL,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(opts.Branch),
	}

	// check for token authentication
	token := opts.Token
	if token == nil {
		if envToken, exists := os.LookupEnv("GIT_TOKEN"); exists {
			token = &envToken
		}
	}
	if token != nil {
		cloneOpts.Auth = &gitHttp.BasicAuth{
			Username: "nox",
			Password: *token,
		}
	}

	repo, err := git.Clone(gitMem.NewStorage(), nil, cloneOpts)
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

// GetFileContentFromRepo fetches the content of a single file from a Git repository branch.
// It supports optional token-based authentication and performs the clone in-memory.
func GetFileContentFromRepo(opts GitFetchOptions) ([]byte, error) {
	if opts.RepoURL == "" || opts.Branch == "" || opts.File == "" {
		return nil, fmt.Errorf("missing required fields: repo, branch, or file")
	}

	cloneOpts := &git.CloneOptions{
		URL:           opts.RepoURL,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(opts.Branch),
	}

	if opts.Token != nil {
		cloneOpts.Auth = &gitHttp.BasicAuth{
			Username: "nox", // Username is required by go-git but can be any non-empty string
			Password: *opts.Token,
		}
	}

	repo, err := git.Clone(gitMem.NewStorage(), nil, cloneOpts)
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

	file, err := tree.File(opts.File)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	reader, err := file.Blob.Reader()
	if err != nil {
		return nil, fmt.Errorf("failed to open file reader: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}
