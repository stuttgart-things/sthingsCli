/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
)

var client *github.Client
var backgroundContext = context.Background()

// GetGitTree GENERATES THE TREE TO COMMIT BASED ON THE GIVEN FILES AND THE COMMIT
// OF THE REF YOU GOT IN GETREF.
func GetGitTree(client *github.Client, ref *github.Reference, sourceFiles []string, sourceOwner, sourceRepo string) (tree *github.Tree, err error) {

	// CREATE A TREE WITH WHAT TO COMMIT.
	entries := []*github.TreeEntry{}

	for _, file := range sourceFiles {

		// CUT STRING INTO SLICES (SOURCE AND TARGET)
		soureTarget := strings.Split(file, ":")
		filePathLocal := soureTarget[0]
		filePathBranch := soureTarget[1]

		// GET FILE CONTENT
		fileContent, _ := ReadFileToVar(filePathLocal)

		// ADD ENTRIES TO GIT TREE
		entries = append(entries, &github.TreeEntry{Path: github.String(filePathBranch), Type: github.String("blob"), Content: github.String(string(fileContent)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(backgroundContext, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
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

// GetReferenceObject RETURNS THE COMMIT BRANCH REFERENCE OBJECT IF IT EXISTS OR CREATES IT
// FROM THE BASE BRANCH BEFORE RETURNING IT.

func GetReferenceObject(client *github.Client, sourceOwner, sourceRepo, commitBranch, baseBranch string) (ref *github.Reference, err error) {

	if ref, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	if commitBranch == baseBranch {
		return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if baseBranch == "" {
		return nil, errors.New("the `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, sourceOwner, sourceRepo, newRef)
	return ref, err
}
