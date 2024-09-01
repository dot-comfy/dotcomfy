package services

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func CloneTo(url string, path string) error {
	// @REF [Basic go-git example](https://github.com/go-git/go-git/blob/master/_examples/clone/main.go)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               url, // Guaranteed at least one because cobra
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	fmt.Println(head.Hash())

	return nil
}
