/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func ReadYamlToObject(pathToYAML, fileExtension string, yamlObject interface{}) interface{} {

	// READ YAML
	viper.SetConfigName(strings.ReplaceAll(filepath.Base(pathToYAML), fileExtension, ""))
	viper.AddConfigPath(filepath.Dir(pathToYAML))

	if err := viper.ReadInConfig(); err != nil {
		return yamlObject
	}

	if err := viper.Unmarshal(&yamlObject); err != nil {
		fmt.Println(err)
		return yamlObject
	}

	return yamlObject

}

func ReadInlineYamlToObject(inlineYAML []byte, yamlObject interface{}) interface{} {

	viper.SetConfigType("yaml")

	// READ YAML
	viper.ReadConfig(bytes.NewBuffer(inlineYAML))

	if err := viper.ReadConfig(bytes.NewBuffer(inlineYAML)); err != nil {
		return yamlObject
	}

	if err := viper.Unmarshal(&yamlObject); err != nil {
		fmt.Println(err)
		return yamlObject
	}

	return yamlObject

}
