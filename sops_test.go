/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"

	"github.com/getsops/sops/v3/stores/yaml"

	"testing"
)

var rawSecretManifest = `apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm`

func TestGenerateAgeIdentitdy(t *testing.T) {

	identity := GenerateAgeIdentitdy()
	fmt.Println(identity.Recipient().String())

}

func TestEncryptStore(t *testing.T) {

	identity := GenerateAgeIdentitdy()
	ageKey := identity.Recipient().String()
	store := yaml.Store{}

	EncryptStore(store, ageKey, rawSecretManifest)

}
