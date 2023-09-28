/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"
	"testing"
)

var (
	downloadURL    = "https://dl.k8s.io/release/v1.27.4/bin/linux/amd64/kubectl"
	downloadTarget = "/tmp"
)

func TestDownloadFileWithProgressBar(t *testing.T) {
	DownloadFileWithProgressBar(downloadURL, downloadTarget)
}

func TestCheckUrlAvailability(t *testing.T) {

	// assert := assert.New(t)

	urlUp := CheckUrlAvailability("www.google.de")
	fmt.Println(urlUp)
	// urlDown := CheckUrlAvailability("www.stuttgart-things.com")

	// assert.Equal(urlUp, true)
	// assert.Equal(urlDown, false)

}
