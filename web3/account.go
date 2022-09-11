package web3

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	Accounts = []*Account{}
	Url      = "/cosmos/auth/v1beta1/accounts"
)

type Account struct {
	PrivKey secp256k1.PrivKey
	Address string
	Number  uint64
	Nonce   uint64
}

func NewAccount(ctx Context, privKey string, addr string) *Account {
	privKeyBytes, _ := hex.DecodeString(privKey)
	key := secp256k1.PrivKey{Key: privKeyBytes}

	number, nonce := getAccountNumberAndNonce(ctx, addr)
	account := &Account{
		PrivKey: key,
		Address: addr,
		Number:  number,
		Nonce:   nonce,
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

func QueryAccount(ctx Context, req *authtypes.QueryAccountRequest) (*authtypes.QueryAccountResponse, error) {
	authClient := authtypes.NewQueryClient(ctx.GRPCClient)
	authRes, err := authClient.Account(
		context.Background(),
		req,
	)
	return authRes, err
}

type Response struct {
	Account AccountResponse `json:"account"`
}
type AccountResponse struct {
	Type            string `json:"@type"`
	Address         string `json:"address"`
	PubKey          PubKey `json:"pub_key"`
	AccountNumber   string `json:"account_number"`
	AccountSequence string `json:"sequence"`
}
type PubKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

func getAccountNumberAndNonce(ctx Context, addr string) (uint64, uint64) {
	req := fmt.Sprintf("%v%v/%v", ctx.API, Url, addr)
	rsp, err := http.Get(req)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	responseData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response Response
	json.Unmarshal(responseData, &response)

	account := response
	number, _ := strconv.ParseUint(account.Account.AccountNumber, 10, 64)
	sequence, _ := strconv.ParseUint(account.Account.AccountSequence, 10, 64)

	return number, sequence
}
