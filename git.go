/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	memfs "github.com/go-git/go-billy/v5/memfs"
	memory "github.com/go-git/go-git/v5/storage/memory"

	git "github.com/go-git/go-git/v5"
)

func GitCloneRepository(repository string, auth *http.BasicAuth) (clonedRepository *git.Repository, cloned bool) {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs := memfs.New()

	// Clone repo into memfs
	clonedRepository, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})

	if err != nil {
		fmt.Println("Could not git clone repository", repository, err)
	} else {
		fmt.Println("Repository cloned", repository)
		cloned = true
	}

	return
}
