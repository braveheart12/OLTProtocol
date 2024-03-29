package keys

import (
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Algorithm int

const (
	UNKNOWN Algorithm = iota
	ED25519
	SECP256K1

	ED25519_PUB_SIZE  int = ed25519.PubKeyEd25519Size
	ED25519_PRIV_SIZE int = 64

	SECP256K1_PUB_SIZE  int = secp256k1.PubKeySecp256k1Size
	SECP256K1_PRIV_SIZE int = 32
)

func (a Algorithm) Name() string {
	switch a {
	case ED25519:
		return "ed25519"
	case SECP256K1:
		return "secp256k1"
	}
	return "Unknown algorithm"
}

func GetAlgorithmFromTmKeyName(name string) Algorithm {

	switch name {
	case "ed25519":
		return ED25519
	case "secp256k1":
		return SECP256K1
	}
	return UNKNOWN
}
