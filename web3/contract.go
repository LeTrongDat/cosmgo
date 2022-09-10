package web3

import (
	"context"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/LeTrongDat/cosmgo/util"
	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type Contract struct{}

func NewContract() *Contract {
	return &Contract{}
}

type MsgExecuteContract struct {
	wasmtypes.MsgExecuteContract
	FeeAmount sdk.Coins
	GasLimit  uint64
}

func NewMsgExecuteContract(sender string, contract string, msg []byte, feeAmount sdk.Coins, gasLimit uint64) *MsgExecuteContract {
	msgExecuteContract := &MsgExecuteContract{}
	msgExecuteContract = msgExecuteContract.
		WithSender(sender).
		WithContract(contract).
		WithMsg(msg).
		WithFeeAmount(feeAmount).
		WithGasLimit(gasLimit)

	return msgExecuteContract
}

func (mec *MsgExecuteContract) WithSender(sender string) *MsgExecuteContract {
	mec.Sender = sender
	return mec
}

func (mec *MsgExecuteContract) WithContract(contract string) *MsgExecuteContract {
	mec.Contract = contract
	return mec
}

func (mec *MsgExecuteContract) WithMsg(msg []byte) *MsgExecuteContract {
	mec.Msg = msg
	return mec
}

func (mec *MsgExecuteContract) WithFeeAmount(feeAmount sdk.Coins) *MsgExecuteContract {
	mec.FeeAmount = feeAmount
	return mec
}

func (mec *MsgExecuteContract) WithGasLimit(gasLimit uint64) *MsgExecuteContract {
	mec.GasLimit = gasLimit
	return mec
}

func (c *Contract) SendTx(ctx Context, msg *MsgExecuteContract) (*tx.BroadcastTxResponse, error) {
	if err := setMsgs(&ctx, msg); err != nil {
		return nil, err
	}

	if err := setSignaturesRound1(&ctx, msg); err != nil {
		return nil, err
	}

	if err := setSignaturesRound2(&ctx, msg); err != nil {
		return nil, err
	}

	rsp, err := broadcastTx(&ctx, msg)
	return rsp, err
}

func setMsgs(ctx *Context, msg *MsgExecuteContract) error {
	txBuilder := extractTxBuilderFromCtx(ctx)
	(*txBuilder).SetFeeAmount(msg.FeeAmount)
	(*txBuilder).SetGasLimit(msg.GasLimit)

	err := (*txBuilder).SetMsgs(msg)
	if err != nil {
		return err
	}
	return nil
}

func setSignaturesRound1(ctx *Context, msg *MsgExecuteContract) error {
	txBuilder := extractTxBuilderFromCtx(ctx)
	signMode := ctx.EncodingConfig.TxConfig.SignModeHandler().DefaultMode()

	var sigsV2 []signing.SignatureV2
	for _, account := range Accounts {
		sigV2 := signing.SignatureV2{
			PubKey: account.PrivKey.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: nil,
			},
			Sequence: account.Nonce,
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := (*txBuilder).SetSignatures(sigsV2...)
	if err != nil {
		return err
	}
	return nil
}

func setSignaturesRound2(ctx *Context, msg *MsgExecuteContract) error {
	txBuilder := extractTxBuilderFromCtx(ctx)
	signMode := ctx.EncodingConfig.TxConfig.SignModeHandler().DefaultMode()

	sigsV2 := []signing.SignatureV2{}
	for _, account := range Accounts {
		signerData := xauthsigning.SignerData{
			ChainID:       ctx.ChainID,
			AccountNumber: account.Number,
			Sequence:      account.Nonce,
		}
		sigV2, err := util.SignWithPrivKey(
			signMode, signerData,
			*txBuilder, cryptotypes.PrivKey(&account.PrivKey),
			ctx.EncodingConfig.TxConfig,
			account.Nonce)

		if err != nil {
			return err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := (*txBuilder).SetSignatures(sigsV2...)
	if err != nil {
		return err
	}
	return nil
}
func extractTxBuilderFromCtx(ctx *Context) *client.TxBuilder {
	if ctx.TxBuilder == nil {
		ctx.TxBuilder = ctx.EncodingConfig.TxConfig.NewTxBuilder()
	}
	return &ctx.TxBuilder
}
func broadcastTx(ctx *Context, msg *MsgExecuteContract) (*tx.BroadcastTxResponse, error) {
	txBuilder := extractTxBuilderFromCtx(ctx)

	transaction := (*txBuilder).GetTx()
	txBytes, err := ctx.EncodingConfig.TxConfig.TxEncoder()(transaction)
	if err != nil {
		return nil, err
	}

	txClient := tx.NewServiceClient(ctx.GRPCClient)
	grpcRes, err := txClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes, // Proto-binary of the signed transaction, see previous step.
		},
	)
	if err != nil {
		return nil, err
	}
	return grpcRes, nil
}

func NewFeeAmount(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(denom, amount))
}

func NewGasLimit(gasLimit uint64) uint64 {
	return gasLimit
}
