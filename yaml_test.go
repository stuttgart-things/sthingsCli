/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"
	"testing"
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

var yamlExample = []byte(`
binary:
  - argocd:
      url: https://github.com/argoproj/argo-cd/releases/download/v{{ or .decksmanVersion "2.5.10" }}/argocd-linux-amd64
      bin: argocd-linux-amd64
  - decksman:
      url: https://artifacts.tiab.labda.sva.de/decksman/decksman-{{ or .decksmanVersion "0.8.436" }}.zip
      bin: decksman
script:
  - argocd:
      script: |
        echo hello
`)

func TestRReadInlineYamlToObject(t *testing.T) {
	config := ReadInlineYamlToObject(yamlExample, config).(Profile)
	fmt.Println(config)
}
