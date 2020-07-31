package simapp

import (
	"fmt"
	"sync"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Wallet struct {
	auths   map[string]crypto.PrivKey
	rootKey crypto.PrivKey

	mutex sync.RWMutex
}

func NewWallet() *Wallet {
	res := &Wallet{
		auths: make(map[string]crypto.PrivKey),
	}
	res.NewRootAddress()
	return res
}

func (w *Wallet) newAccAddressNoMutex() types.AccAddress {
	privKey := secp256k1.GenPrivKey()
	address := types.AccAddress(privKey.PubKey().Address())

	addressStr := types.NewAccountIDFromAccAdd(address).String()

	w.auths[addressStr] = privKey

	return address
}

// NewAccAddress new AccAddress and save privkey
func (w *Wallet) NewAccAddress() types.AccAddress {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	address := w.newAccAddressNoMutex()

	return address
}

// NewAccAddressByName new AccAddress and save privkey with name
func (w *Wallet) NewAccAddressByName(name types.Name) types.AccAddress {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	address := w.newAccAddressNoMutex()

	w.auths[name.String()] = w.auths[address.String()]

	return address
}

// NewRootAddress
func (w *Wallet) NewRootAddress() types.AccAddress {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	res := w.newAccAddressNoMutex()
	w.rootKey = w.auths[types.NewAccountIDFromAccAdd(res).String()]
	return res
}

// GetAuth
func (w *Wallet) GetAuth(id types.AccountID) types.AccAddress {
	if acc, ok := id.ToAccAddress(); ok {
		return acc
	}

	if n, ok := id.ToName(); ok {
		res, ok := w.auths[n.String()]
		if !ok {
			panic(fmt.Errorf("no found auth by %s", n.String()))
		}
		return types.AccAddress(res.PubKey().Address())
	}

	panic(fmt.Errorf("id type not support %s", id.String()))
}

// GetRootAuth
func (w *Wallet) GetRootAuth() types.AccAddress {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return types.AccAddress(w.rootKey.PubKey().Address())
}

// GetRootKey
func (w *Wallet) GetRootKey() crypto.PrivKey {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.rootKey
}

// PrivKey get private key
func (w *Wallet) PrivKey(key types.AccAddress) crypto.PrivKey {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	addressStr := types.NewAccountIDFromAccAdd(key).String()

	p, ok := w.auths[addressStr]
	if !ok {
		panic(errors.Errorf("no found private key for %s", key.String()))
	}

	return p
}
