/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

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

func ConvertYAMLToJSON(yamlString string) string {

	convertedJSON, err := yaml.YAMLToJSON([]byte(strings.Replace(yamlString, "\t", "\n", 5)))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return string(convertedJSON)
}

func ConvertJSONToYAML(jsonString string) string {

	// REMOVE ALL MAY EXISTING IVALID CHARS
	jsonString = strings.ReplaceAll(jsonString, "\\", "")
	jsonString = strings.TrimRight(jsonString, "\"")
	jsonString = strings.TrimLeft(jsonString, "\"")

	convertedYAML, err := yaml.JSONToYAML([]byte(jsonString))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return string(convertedYAML)
}
