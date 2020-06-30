package utils

import (
	"crypto/sha256"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

const mnemonic = "bamboo dentist fold word pill patch weekend loop pattern portion shoulder coyote start bone hockey amount cool poverty census food during slight fix use"

func TestCrypto(t *testing.T) {
	derivedPriv, err := keys.StdDeriveKey(mnemonic, "", "44'/23808'/0'/0/0", keys.Secp256k1)

	if err != nil {
		t.Errorf("111, %v", err)
	}

	t.Logf("derivedPriv, %x", derivedPriv)

	privKey, err := keys.StdPrivKeyGen(derivedPriv, keys.Secp256k1)

	if err != nil {
		t.Errorf("222, %v", err)
	}

	t.Logf("privKey, %x", privKey)
	t.Logf("PubKey, %x", privKey.PubKey())

	accAddr := sdk.AccAddress(privKey.PubKey().Address().Bytes())
	t.Logf("accAddr, %x", accAddr.String())

	bytes := []byte{
		0x1,
		0x23,
		0x34,
	}

	hasher := sha256.New()
	hasher.Write(bytes)
	t.Logf("buf, %x", hasher.Sum(nil))

	sig, err := privKey.Sign(bytes)
	if err != nil {
		t.Errorf("777, %v", err)
	}

	t.Logf("sig, %x", sig)
}
