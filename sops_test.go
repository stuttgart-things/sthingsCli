/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"fmt"

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

	EncryptStore(ageKey, rawSecretManifest)

}
