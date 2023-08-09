/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"testing"
)

var (
	repo = "https://github.com/stuttgart-things/machineShop.git"
)

func TestCloneGitRepository(t *testing.T) {

	GitCloneRepository(repo, nil)
}
