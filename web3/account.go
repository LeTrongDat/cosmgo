package web3

import (
	"encoding/hex"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	Accounts = []*Account{}
)

type Account struct {
	PrivKey secp256k1.PrivKey
	Address sdk.AccAddress
	Number  uint64
	Nonce   uint64
}

func NewAccount(privKey string) *Account {
	privKeyBytes, _ := hex.DecodeString(privKey)
	key := secp256k1.PrivKey{Key: privKeyBytes}
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())

	account := &Account{
		PrivKey: key,
		Address: addr,
	}
	Accounts = append(Accounts, account)
	return account
}

func (acc *Account) WithNumber(number uint64) *Account {
	acc.Number = number
	return acc
}

func (acc *Account) WithNonce(nonce uint64) *Account {
	acc.Nonce = nonce
	return acc
}
