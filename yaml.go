/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"

	"sigs.k8s.io/yaml"

	"github.com/spf13/viper"
)

func ReadYamlKeyValuesFromFile(yamlFileContent []byte) (yamlStructure map[string]interface{}) {

	yamlStructure = make(map[string]interface{})
	data := make(map[interface{}]interface{})

	err := yamlv3.Unmarshal(yamlFileContent, &data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range data {
		yamlStructure[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
	}

	return
}

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
