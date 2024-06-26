/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"
	"io"
	"os"

	billy "github.com/go-git/go-billy/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	memfs "github.com/go-git/go-billy/v5/memfs"
	memory "github.com/go-git/go-git/v5/storage/memory"

	git "github.com/go-git/go-git/v5"
)

type FilesToAdd struct {
	Filename    string
	Filecontent []byte
}

func ReadFileContentFromGitRepo(repo billy.Filesystem, filePath string) string {

	// OPEN FILE
	file, err := repo.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	// READ FILE CONTENT
	fileContent, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}

	// RETURN STRING
	return string(fileContent)

}

func CloneGitRepository(repository, branchName, commitId string, auth *http.BasicAuth) (fs billy.Filesystem, cloned bool) {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs = memfs.New()

	// Clone repo into memfs
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           repository,
		Auth:          auth,
		RemoteName:    "origin",
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(branchName),
	})

	if err != nil {
		fmt.Println("Could not git clone repository", repository, err)
	} else {
		fmt.Println("Repository cloned")
		cloned = true
	}

	// CHECK OUT SPECIFIC COMMIT BY ID (IF COMMIT IS GIVEN)
	if commitId != "" {

		ref, err := r.Head()
		fmt.Println(ref, err)

		w, err := r.Worktree()
		fmt.Println(err)

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commitId),
		})

		fmt.Println(err)

		ref, _ = r.Head()

		fmt.Println(ref.Hash())
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

func AddCommitFileToGitRepository(repository, branchName string, auth *http.BasicAuth, add []FilesToAdd, remove []string, commitMsg string) (pushed bool) {

	// INIT MEMORY STORAGE AND FS
	storer := memory.NewStorage()
	fs := memfs.New()

	// CLONE REPO INTO MEMFS
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:           repository,
		Auth:          auth,
		RemoteName:    "origin",
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(branchName),
	})

	// CLONE REPO INTO MEMFS
	if err != nil {
		fmt.Println("Could not git clone repository")
	}
	fmt.Println("Repository cloned")

	// GET GIT DEFAULT WORKTREE
	w, err := r.Worktree()
	if err != nil {
		fmt.Println("Could not get git worktree")
	}

	fmt.Println(w)

	// CREATE NEW FILES TO WORKTREE
	for _, file := range add {

		newFile, err := fs.Create(file.Filename)
		if err != nil {
			fmt.Println("COULD NOT CREATE FILE", file.Filename)
		}
		newFile.Write(file.Filecontent)
		newFile.Close()

		// RUN GIT STATUS BEFORE ADDING THE FILE TO THE WORKTREE
		fmt.Println(w.Status())

		// GIT ADD $FILEPATH
		w.Add(file.Filename)
	}

	// REMOVE FILES
	for _, file := range remove {

		err2 := fs.Remove(file)
		if err2 != nil {
			fmt.Println("COULD NOT REMOVE FILE")
		}

		w.Add(file)
	}
	// RUN GIT STATUS AFTER THE FILE HAS BEEN ADDED ADDING TO THE WORKTREE
	fmt.Println(w.Status())

	// GIT COMMIT -M $MESSAGE
	w.Commit(commitMsg, &git.CommitOptions{})

	// PUSH THE CODE TO THE REMOTE
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		fmt.Println("COULD NOT GIT PUSH: %w", err)
	}
	fmt.Println("REMOTE UPDATED.")

	pushed = true

	return
}
