/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"

	"filippo.io/age"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	keysource "github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/stores/yaml"
)

var sopsVersion = "3.8.1"
var unencryptedSuffix = "_unencrypted"

func EncryptStore(ageKey, rawData string) (encryptedData string) {

	store := yaml.Store{}

	branches, err := store.LoadPlainFile([]byte(rawData))
	if err != nil {
		panic(err)
	}
	fmt.Println(branches)

	masterKey, err := keysource.MasterKeyFromRecipient(ageKey)
	if err != nil {
		panic(err)
	}
	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			KeyGroups: []sops.KeyGroup{
				[]keys.MasterKey{masterKey},
			},
			UnencryptedSuffix: unencryptedSuffix,
			Version:           sopsVersion,
		},
	}

	dataKey, errs := tree.GenerateDataKeyWithKeyServices(
		[]keyservice.KeyServiceClient{keyservice.NewLocalClient()},
	)
	if errs != nil {
		panic(errs)
	}
	common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  aes.NewCipher(),
	})

	result, err := store.EmitEncryptedFile(tree)
	if err != nil {
		panic(err)
	}

	encryptedData = string(result)

	return encryptedData
}

func GenerateAgeIdentitdy() (identity *age.X25519Identity) {

	identity, err := age.GenerateX25519Identity()
	if err != nil {
		panic(err)
	}

	return
}
