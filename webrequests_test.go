/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import "testing"

var (
	downloadURL    = "https://dl.k8s.io/release/v1.27.4/bin/linux/amd64/kubectl"
	downloadTarget = "/tmp"
)

func TestDownloadFileWithProgressBar(t *testing.T) {
	DownloadFileWithProgressBar(downloadURL, downloadTarget)
}
