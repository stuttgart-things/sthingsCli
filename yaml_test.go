/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var config Profile

type Install struct {
	Url string `mapstructure:"url"`
	Bin string `mapstructure:"bin"`
}

type Script struct {
	Script string `mapstructure:"script"`
}

type Profile struct {
	BinaryProfile []map[string]Install `mapstructure:"binary"`
	ScriptProfile []map[string]Script  `mapstructure:"script"`
}

var expectedUrl = "https://github.com/argoproj/argo-cd/releases/download/v2.5.10/argocd-linux-amd64"
var expectedScript = "echo hello"

var yamlExample = []byte(`
binary:
  - argocd:
      url: https://github.com/argoproj/argo-cd/releases/download/v2.5.10/argocd-linux-amd64
      bin: argocd-linux-amd64
script:
  - argocd:
      script: |
        echo hello`)

func TestRReadInlineYamlToObject(t *testing.T) {

	assert := assert.New(t)

	config := ReadInlineYamlToObject(yamlExample, config).(Profile)

	// CHECK FOR ONE VALUE IN BINARY
	assert.Equal(expectedUrl, config.BinaryProfile[0]["argocd"].Url)

	// CHECK FOR ONE VALUE IN SCRIPT
	assert.Equal(expectedScript, config.ScriptProfile[0]["argocd"].Script)

}
