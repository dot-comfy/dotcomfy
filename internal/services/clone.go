package services

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	Log "dotcomfy/internal/logger"
)

func Clone(url string, branch, path string) error {
	LOGGER = Log.GetLogger()
	// @REF [Basic go-git example](https://github.com/go-git/go-git/blob/master/_examples/clone/main.go)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               url, // Guaranteed at least one because cobra
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.ReferenceName(branch),
		SingleBranch:      true,
	})
	if err != nil {
		LOGGER.Error(err)
		return err
	}

	head, err := repo.Head()
	if err != nil {
		LOGGER.Error(err)
		return err
	}

	fmt.Println(head.Hash())

	return nil
}
