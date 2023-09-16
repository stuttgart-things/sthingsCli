/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	repo                  = "https://github.com/stuttgart-things/machineShop.git"
	branchName            = "main"
	commitID              = "1ce75f510bff3b95b0a5998ee22731ec058c3267"
	expectedFileList      = []string{".gitignore", "LICENSE", "README.md"}
	expectedDirectoryList = []string{}
	testCommitData        = []byte("ABC")
	gitToken              = os.Getenv("GITHUB_TOKEN")
	gitUser               = "patrick-hermann-sva"
)

func TestCloneGitRepository(t *testing.T) {

	assert := assert.New(t)

	_, cloned := CloneGitRepository(repo, branchName, "", nil)

	assert.Equal(cloned, true)

}

func TestGetFileListFromGitRepository(t *testing.T) {

	var fileList []string
	var directoryList []string

	repo, cloned := CloneGitRepository(repo, branchName, commitID, nil)

	if cloned {
		fileList, directoryList = GetFileListFromGitRepository("", repo)
		fmt.Println(fileList, directoryList)
	}

	if !reflect.DeepEqual(fileList, expectedFileList) && reflect.DeepEqual(directoryList, expectedDirectoryList) {
		t.Errorf("expected lists differ")
	} else {
		fmt.Println("test successfully")
	}

}

func TestAddCommitFileToGitRepository(t *testing.T) {

	//auth := CreateGitAuth(gitUser, gitToken)
	//AddCommitFileToGitRepository(repo, branchName, auth, testCommitData, "test.txt", "added test file w/ golang")

}
