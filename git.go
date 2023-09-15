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

func AddCommitFileToGitRepository(repository string, auth *http.BasicAuth, fileContent []byte, filePath, commitMsg string) error {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs := memfs.New()

	// Clone repo into memfs
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})
	if err != nil {
		return fmt.Errorf("Could not git clone repository %s: %w", repository, err)
	}
	fmt.Println("Repository cloned")

	// Get git default worktree
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("Could not get git worktree: %w", err)
	}

	fmt.Println(w)

	// Create new file
	newFile, err := fs.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create new file: %w", err)
	}
	newFile.Write(fileContent)
	newFile.Close()

	// Run git status before adding the file to the worktree
	fmt.Println(w.Status())

	// git add $filePath
	w.Add(filePath)

	// Run git status after the file has been added adding to the worktree
	fmt.Println(w.Status())

	// git commit -m $message
	w.Commit(commitMsg, &git.CommitOptions{})

	//Push the code to the remote
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		return fmt.Errorf("Could not git push: %w", err)
	}
	fmt.Println("Remote updated.", filePath)

	return nil
}
