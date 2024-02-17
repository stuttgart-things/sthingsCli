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

func TestGetCommitInformationFromGithubRepo(t *testing.T) {

	bla, he, err := GetCommitInformationFromGithubRepo("stuttgart-things", "kaeffken", "feature/issue-1/test", "latest")
	fmt.Println(bla, he, err)

}

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

func TestReadFileContentFromGitRepo(t *testing.T) {

	gitRepository := "https://github.com/stuttgart-things/kaeffken.git"
	gitBranch := "main"
	gitCommitID := "09de9ff7b5c76aff8bb32f68cfb0bbe49cd5a7a8"

	assert := assert.New(t)
	expectedReadMe := "# kaeffken\ngitops cluster management cli \n"

	repo, _ := CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
	readMe := ReadFileContentFromGitRepo(repo, "README.md")
	fmt.Println(readMe)
	fmt.Println(expectedReadMe)

	assert.Equal(readMe, expectedReadMe)
	fmt.Println("TEST SUCCESSFULLY")

}
