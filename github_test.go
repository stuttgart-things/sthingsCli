/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/go-github/v68/github"

	"github.com/stretchr/testify/assert"
)

func TestGetReferenceObject_WithExistingBranch(t *testing.T) {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("UNAUTHORIZED: NO TOKEN PRESENT")
	}

	client = github.NewClient(nil).WithAuthToken(token)

	// CALL GETREFERENCEOBJECT
	GetReferenceObject(client, "stuttgart-things", "machineshop", "test", "main")

}

func TestReadFileToVar_WithValidFile(t *testing.T) {
	// CREATE A TEMPORARY FILE
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// WRITE SOME DATA TO THE FILE
	text := []byte("Hello, World!")
	if _, err := tmpfile.Write(text); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// CALL THE READFILETOVAR FUNCTION
	content, err := ReadFileToVar(tmpfile.Name())

	// ASSERT THAT THE RETURNED CONTENT IS THE SAME AS THE DATA WRITTEN TO THE FILE
	assert.Equal(t, text, content)

	// ASSERT THAT THE RETURNED ERROR IS NIL
	assert.Nil(t, err)
}

func TestReadFileToVar_WithNonExistentFile(t *testing.T) {
	// CALL THE READFILETOVAR FUNCTION WITH A NON-EXISTENT FILE PATH
	content, err := ReadFileToVar("non_existent_file.txt")

	// ASSERT THAT THE RETURNED CONTENT IS NIL
	assert.Nil(t, content)

	// ASSERT THAT THE RETURNED ERROR IS NOT NIL
	assert.NotNil(t, err)
}

// func TestMergePullRequest(t *testing.T) {

// 	token := ""
// 	if token == "" {
// 		log.Fatal("UNAUTHORIZED: NO TOKEN PRESENT")
// 	}

// 	client = github.NewClient(nil).WithAuthToken(token)

// 	MergePullRequest(client, "stuttgart-things", "stuttgart-things", "merge", "merge", 241)

// }

func TestGetCommitInformationFromGithubRepo(t *testing.T) {

	commitExists, commitInformation, err := GetCommitInformationFromGithubRepo("stuttgart-things", "kaeffken", "main", "latest")
	fmt.Println(commitExists, commitInformation, err)

}
