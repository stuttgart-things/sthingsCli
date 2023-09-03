/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	zipDownloadURL    = "https://github.com/stuttgart-things/sthingsCli/archive/refs/tags/v0.1.28.zip"
	zipFileArchive    = "v0.1.28.zip"
	zipFileUnArchived = "sthingsCli-0.1.28"
)

func TestUnzipArchive(t *testing.T) {

	assert := assert.New(t)

	DownloadFileWithProgressBar(zipDownloadURL, downloadTarget)
	UnzipArchive(downloadTarget+"/"+zipFileArchive, downloadTarget)
	unarchivedDirExist, _ := sthingsBase.VerifyIfFileOrDirExists(downloadTarget+"/"+zipFileUnArchived, "dir")

	assert.Equal(unarchivedDirExist, true)

}
