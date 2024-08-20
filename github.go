/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
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

// PushCommit CREATES THE COMMIT IN THE GIVEN REFERENCE USING THE GIVEN TREE
func PushCommit(client *github.Client, ref *github.Reference, tree *github.Tree, sourceOwner, sourceRepo, authorName, authorEmail, commitMessage string) (err error) {

	// GET THE PARENT COMMIT TO ATTACH THE COMMIT TO.
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}

	// THIS IS NOT ALWAYS POPULATED, BUT IS NEEDED.
	parent.Commit.SHA = parent.SHA

	// CREATE THE COMMIT USING THE TREE.
	date := time.Now()
	author := &github.CommitAuthor{Date: &github.Timestamp{Time: date}, Name: sthingsBase.ConvertStringToPointer(authorName), Email: sthingsBase.ConvertStringToPointer(authorEmail)}
	commit := &github.Commit{Author: author, Message: sthingsBase.ConvertStringToPointer(commitMessage), Tree: tree, Parents: []*github.Commit{parent.Commit}}
	opts := github.CreateCommitOptions{}

	newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit, &opts)
	if err != nil {
		return err
	}

	// ATTACH THE COMMIT TO THE MASTER BRANCH.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	return err
}

// CreatePullRequest CREATES A PULL REQUEST. BASED ON: HTTPS://GODOC.ORG/GITHUB.COM/GOOGLE/GO-GITHUB/GITHUB#EXAMPLE-PULLREQUESTSSERVICE-CREATE
func CreatePullRequest(client *github.Client, prSubject, prRepoOwner, sourceOwner, commitBranch, prRepo, sourceRepo, repoBranch, baseBranch, prDescription string, prLabels []string) (err error, prId string) {

	if prRepoOwner != "" && prRepoOwner != sourceOwner {
		commitBranch = fmt.Sprintf("%s:%s", sourceOwner, commitBranch)
	} else {
		prRepoOwner = sourceOwner
	}

	if prRepo == "" {
		prRepo = sourceRepo
	}

	newPR := &github.NewPullRequest{
		Title:               sthingsBase.ConvertStringToPointer(prSubject),
		Head:                sthingsBase.ConvertStringToPointer(commitBranch),
		HeadRepo:            sthingsBase.ConvertStringToPointer(repoBranch),
		Base:                sthingsBase.ConvertStringToPointer(baseBranch),
		Body:                sthingsBase.ConvertStringToPointer(prDescription),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, prRepoOwner, prRepo, newPR)
	if err != nil {
		return err, "NONE"
	}

	fmt.Printf("PR CREATED: %s\n", pr.GetHTMLURL())
	// for gettimg all fileds fmt.Println(pr)

	if len(prLabels) != 0 {

		newLabels, _, err := client.Issues.AddLabelsToIssue(ctx, prRepoOwner, prRepo, int(*pr.Number), prLabels)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ADDED LABELS TO PR", newLabels)
	}

	prId = strconv.Itoa(int(*pr.Number))

	return nil, prId
}

// CreateRepository CREATES A GITHUB REPOSITORY
func CreateRepository(client *github.Client, name, description, repoOwner string, privateRepo, autoInit bool) (err error, repoName string) {

	r := &github.Repository{
		Name:        sthingsBase.ConvertStringToPointer(name),
		Private:     sthingsBase.ConvertBoolToPointer(privateRepo),
		Description: sthingsBase.ConvertStringToPointer(description),
		AutoInit:    sthingsBase.ConvertBoolToPointer(autoInit),
	}

	repo, _, err := client.Repositories.Create(ctx, repoOwner, r)
	if err != nil {
		return err, ""
	} else {
		fmt.Printf("SUCCESSFULLY CREATED NEW REPO: %v\n", repo.GetName())
	}

	return err, repo.GetName()
}

// CreateRepository CREATES A GITHUB REPOSITORY
func MergePullRequest(client *github.Client, repository, repoOwner, message, mergeMethod string, pullRequestID int) {

	options := &github.PullRequestOptions{
		CommitTitle:        repository,
		MergeMethod:        mergeMethod,
		DontDefaultIfBlank: false,
	}

	mergeResult, resp, _ := client.PullRequests.Merge(ctx, repoOwner, repository, pullRequestID, message, options)
	fmt.Println(mergeResult, resp)

}

// DeleteBranch DELETES A BRANCH
func DeleteBranch(client *github.Client, repository, repoOwner, branchName string) {

	response, error := client.Git.DeleteRef(ctx, repoOwner, repository, branchName)
	fmt.Println(response, error)

}

func GetCommitInformationFromGithubRepo(userName, repoName, branchName, option string) (getCommits bool, allCommits []map[string]interface{}, err error) {

	client := github.NewClient(nil)

	opt := &github.CommitsListOptions{
		SHA: branchName,
	}

	commits, response, err := client.Repositories.ListCommits(context.Background(), userName, repoName, opt)

	if err != nil && response.StatusCode != 200 {
		fmt.Println(err)
		return
	} else {
		getCommits = true
	}

	if option == "latest" {

		commitInformation := make(map[string]interface{})
		commitInformation["REVISION"] = sthingsBase.GetStringPointerValue(commits[0].SHA)
		commitInformation["AUTHOR"] = sthingsBase.GetStringPointerValue(commits[0].Author.Login)

		allCommits = append(allCommits, commitInformation)

	} else {

		for _, commit := range commits {
			commitInformation := make(map[string]interface{})

			commitInformation["REVISION"] = sthingsBase.GetStringPointerValue(commit.SHA)
			commitInformation["AUTHOR"] = sthingsBase.GetStringPointerValue(commit.Author.Login)

			allCommits = append(allCommits, commitInformation)

		}
	}

	return
}
