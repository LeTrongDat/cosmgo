package main

import (
	"fmt"
	"os"

	"github.com/LeTrongDat/cosmgo/util"
	"github.com/LeTrongDat/cosmgo/web3"
	"google.golang.org/grpc"
)

const (
	// chain config
	API                 string = "http://10.1.0.44:1317"
	GRPC_Server_Address string = "10.1.0.44:9090"
	Chain_ID            string = "baseblockchain"
	Bech32_Prefix       string = "sosc"

	// contract config
	Contract_Address string = "sosc1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrq43khmt"
	From_Address     string = "sosc154qv5c37vq4nd7m4p6dzvvar6zthch25tpjgmv"
	Private_Key      string = "ea2e25468b256bad1c391a23c24f03e07627a52fcffafbd78c2e07ccd0df2bc2"
	Msg              string = `{"register":{"name":"fred"}}`
)

func main() {
	// err := contract.SendTx("sosc154qv5c37vq4nd7m4p6dzvvar6zthch25tpjgmv", `{"register":{"name":"fred"}}`)
	// fmt.Print(err)

	// set up grpc connection
	grpcConn, err := grpc.Dial(
		GRPC_Server_Address, // Or your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer grpcConn.Close()

	// init client context

	clientCtx := web3.NewContext()
	clientCtx = clientCtx.WithChainID(Chain_ID).
		WithGRPCClient(grpcConn).
		WithAPI(API).
		WithBech32Prefix(Bech32_Prefix)

	// init account
	account := web3.NewAccount(clientCtx, Private_Key, From_Address)

	// init contract
	contract := web3.NewContract()
	msgExecuteContract := web3.NewMsgExecuteContract(
		account.Address,
		Contract_Address,
		[]byte(Msg),
		util.NewFeeAmount("stake", 10000),
		util.NewGasLimit(200000),
	)
	rsp, err := contract.SendTx(clientCtx, msgExecuteContract)
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
	fmt.Println("response", rsp.TxResponse)
}
