/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

var (
	vaultAddr      = os.Getenv("VAULT_ADDR")
	roleId         = os.Getenv("VAULT_ROLE_ID")
	secretId       = os.Getenv("VAULT_SECRET_ID")
	vaultNamespace = os.Getenv("VAULT_NAMESPACE")
)

type Client struct {
	Token         string
	Client        *api.Client
	LeaseDuration int
}

type Token struct {
	Auth struct {
		ClientToken   string `json:"client_token"`
		LeaseDuration int    `json:"lease_duration"`
	} `json:"auth"`
}

type appRoleLogin struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

func GetVaultSecretValue(kvpath, token string) string {

	if !CheckVaultKVExistenceInSecretPath(kvpath, token) {
		fmt.Printf("Vault secret does not exist: %s\n", kvpath)
		return ""
	}

	kv_path := strings.Split(kvpath, ":")

	config := &api.Config{
		Address: vaultAddr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}

	client.SetToken(token)
	client.SetNamespace(vaultNamespace)

	secret, err := client.Logical().Read(kv_path[0])
	if err != nil {
		fmt.Println(err)
	}

	vaultSecretValue := secret.Data["data"].(map[string]interface{})

	return vaultSecretValue[kv_path[1]].(string)
}

func CheckVaultKVExistenceInSecretPath(kvpath, token string) bool {

	kv_path := strings.Split(kvpath, ":")
	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}

	client.SetToken(token)
	client.SetNamespace(vaultNamespace)

	secret, err := client.Logical().Read(kv_path[0])
	if err != nil {
		fmt.Println(err)
	}

	if secret == nil || secret.Data == nil || secret.Data["data"] == nil {
		return false
	}

	vaultSecretValue := secret.Data["data"].(map[string]interface{})

	if vaultSecretValue[kv_path[1]] == nil {
		return false
	}

	return true
}

func VerifyEnvVars(envVars []string) bool {

	set := true

	for _, v := range envVars {

		_, ok := os.LookupEnv(v)

		if !ok {
			fmt.Println(v, "is not present..")
			return false

		} else {
			fmt.Println(v, "is present")
		}
	}

	return set
}

func CreateVaultClient() (*Client, error) {
	vaultClient := Client{}

	client, err := api.NewClient(&api.Config{
		Address: vaultAddr,
	})

	vaultClient.Client = client

	return &vaultClient, err
}

func (vaultClient *Client) GetVaultTokenFromAppRole() (string, error) {
	// step: create the token request

	request := vaultClient.Client.NewRequest("POST", "/v1/auth/approle/login")
	login := appRoleLogin{SecretID: secretId, RoleID: roleId}

	if err := request.SetJSONBody(login); err != nil {
		return "", err
	}

	// step: make the request
	resp, err := vaultClient.Client.RawRequest(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// step: parse and return auth
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

func GetVaultKvSecretEngines(vaultAddr, vaultToken, vaultNamespace string) []string {

	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {

		fmt.Println(err)
		return nil
	}
	client.SetToken(vaultToken)
	client.SetNamespace(vaultNamespace)

	mounts, _ := client.Sys().ListMounts()
	var allMounts []string

	for key := range mounts {
		allMounts = append(allMounts, strings.TrimRight(key, "/"))
	}

	return allMounts
}

func GetVaultKvSecretPaths(vaultAddr, vaultToken, vaultNamespace, kvpath string) []string {

	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {

		fmt.Println(err)
		return nil
	}
	client.SetToken(vaultToken)
	client.SetNamespace(vaultNamespace)

	paths, _ := client.Logical().List(kvpath)

	if paths == nil {
		fmt.Println("secret path " + kvpath + " not found")
		os.Exit(3)
	}

	m := paths.Data["keys"].([]interface{})

	var allpaths []string

	for _, p := range m {
		allpaths = append(allpaths, p.(string))
	}

	return allpaths
}

func ReadVaultSecretEngines(vaultAddr, vaultToken, vaultNamespace, kvpath string) []string {

	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {

		fmt.Println(err)
	}
	client.SetToken(vaultToken)
	client.SetNamespace(vaultNamespace)

	secret, _ := client.Logical().Read(kvpath)

	var allKeys []string

	keys := secret.Data["data"].(map[string]interface{})

	for k := range keys {
		allKeys = append(allKeys, k)
	}

	return allKeys
}

func StoreSecretInSecretEngine(vaultAddr, vaultToken, vaultNamespace, secretEngine, secretName string, secretData map[string]interface{}) {

	config := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(config)
	if err != nil {

		fmt.Println(err)
	}
	client.SetToken(vaultToken)
	client.SetNamespace(vaultNamespace)

	ctx := context.Background()

	// Write a secret
	_, err = client.KVv2(secretEngine).Put(ctx, secretName, secretData)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}

	log.Println("Secret written successfully.")

}

func GetSecretValueFromVaultSecretEngine(vaultConnectionInformation map[string]string) (string, string, string, string) {

	vaultSecretEngine := AskSingleSelectQuestion("Select secret engine:", GetVaultKvSecretEngines(vaultConnectionInformation["addr"], vaultConnectionInformation["token"], vaultConnectionInformation["namespace"]))
	vaultSecretPath := GetVaultKvSecretPaths(vaultConnectionInformation["addr"], vaultConnectionInformation["token"], vaultConnectionInformation["namespace"], vaultSecretEngine+"/metadata")
	vaultSecretName := AskSingleSelectQuestion("Select secret", vaultSecretPath)
	vaultSecretKeyNames := ReadVaultSecretEngines(vaultConnectionInformation["addr"], vaultConnectionInformation["token"], vaultConnectionInformation["namespace"], vaultSecretEngine+"/data/"+vaultSecretName)
	vaultSecretKeyName := AskSingleSelectQuestion("Select secret key:", vaultSecretKeyNames)
	secretValue := GetVaultSecretValue(vaultSecretEngine+"/data/"+vaultSecretName+":"+vaultSecretKeyName, vaultConnectionInformation["token"])

	return secretValue, vaultSecretName, vaultSecretKeyName, vaultSecretEngine
}

func GetVaultConnectionInformation() map[string]string {

	vaultConnectionInformation := make(map[string]string)

	if AskSingleSelectQuestion("Get vault connection information from:", []string{"env", "prompt"}) == "env" {
		VerifyEnvVars([]string{"VAULT_ADDR", "VAULT_TOKEN", "VAULT_NAMESPACE"})

		vaultConnectionInformation["addr"] = os.Getenv("VAULT_ADDR")
		vaultConnectionInformation["token"] = os.Getenv("VAULT_TOKEN")
		vaultConnectionInformation["namespace"] = os.Getenv("VAULT_NAMESPACE")

	} else {

		fmt.Println("to be built in")

	}

	return vaultConnectionInformation

}
