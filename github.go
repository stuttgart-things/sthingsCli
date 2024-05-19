/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

var client *github.Client

// var ctx = context.Background()

// GetGitTree GENERATES THE TREE TO COMMIT BASED ON THE GIVEN FILES AND THE COMMIT
// OF THE REF YOU GOT IN GETREF.
func GetGitTree(ref *github.Reference, sourceFiles, sourceOwner, sourceRepo string) (tree *github.Tree, err error) {

	// CREATE A TREE WITH WHAT TO COMMIT.
	entries := []github.TreeEntry{}

	for _, file := range strings.Split(sourceFiles, ",") {

		// CUT STRING INTO SLICES (SOURCE AND TARGET)
		soureTarget := strings.Split(file, ":")
		filePathLocal := soureTarget[0]
		filePathBranch := soureTarget[1]

		// GET FILE CONTENT
		fileContent, _ := ReadFileToVar(filePathLocal)

		// ADD ENTRIES TO GIT TREE
		entries = append(entries, github.TreeEntry{Path: github.String(filePathBranch), Type: github.String("blob"), Content: github.String(string(fileContent)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

func ReadFileToVar(file string) (content []byte, err error) {

	_, content, err = GetFileContent(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func GetFileContent(fileArg string) (targetName string, b []byte, err error) {

	var localFile string
	files := strings.Split(fileArg, ":")

	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
}
