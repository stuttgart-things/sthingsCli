/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	repo       = "https://github.com/stuttgart-things/machineShop.git"
	branchName = "main"
)

func TestCloneGitRepository(t *testing.T) {

	assert := assert.New(t)

	_, cloned := CloneGitRepository(repo, branchName, nil)

	assert.Equal(cloned, true)

}

func TestGetFileListFromGitRepository(t *testing.T) {

	assert := assert.New(t)

	_, cloned := CloneGitRepository(repo, branchName, nil)
	assert.Equal(cloned, true)

}
