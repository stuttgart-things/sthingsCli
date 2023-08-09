/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"

	billy "github.com/go-git/go-billy/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	memfs "github.com/go-git/go-billy/v5/memfs"
	memory "github.com/go-git/go-git/v5/storage/memory"

	git "github.com/go-git/go-git/v5"
)

func CloneGitRepository(repository, branchName string, auth *http.BasicAuth) (fs billy.Filesystem, cloned bool) {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs = memfs.New()

	// Clone repo into memfs
	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           repository,
		Auth:          auth,
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(branchName),
	})

	if err != nil {
		fmt.Println("Could not git clone repository", repository, err)
	} else {
		fmt.Println("Repository cloned")
		cloned = true
	}

	return
}

func GetFileListFromGitRepository(directory string, fs billy.Filesystem) (fileList, directoryList []string) {

	files, _ := fs.ReadDir(directory)

	for _, file := range files {

		if file.IsDir() {
			directoryList = append(directoryList, file.Name())
		} else {
			fileList = append(fileList, file.Name())
		}
	}

	return

}

func CreateGitAuth(gitUser, gitToken string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: gitUser,
		Password: gitToken,
	}
}
