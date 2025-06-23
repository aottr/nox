package gitrepo

import (
	"github.com/go-git/go-git/v5/plumbing/object"
)

func FileExistsInTree(tree *object.Tree, path string) bool {
	_, err := tree.File(path)
	return err == nil
}
