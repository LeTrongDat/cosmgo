package web3

import (
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Context implements a typical context created in SDK modules for transaction
// handling and queries.
type Context struct {
	FromAddress    sdk.AccAddress
	GRPCClient     *grpc.ClientConn
	ChainID        string
	EncodingConfig EncodingConfig
	TxBuilder      client.TxBuilder
	BroadcastMode  string
	SignModeStr    string
	Bech32Prefix   string
}

func NewContext() Context {
	return Context{
		EncodingConfig: MakeEncodingConfig(),
	}
}

// WithGRPCClient returns a copy of the context with an updated GRPC client
// instance.
func (ctx Context) WithGRPCClient(grpcClient *grpc.ClientConn) Context {
	ctx.GRPCClient = grpcClient
	return ctx
}

// WithChainID returns a copy of the context with an updated chain ID.
func (ctx Context) WithChainID(chainID string) Context {
	ctx.ChainID = chainID
	return ctx
}

// WithFromAddress returns a copy of the context with an updated from account
// address.
func (ctx Context) WithFromAddress(addr sdk.AccAddress) Context {
	ctx.FromAddress = addr
	return ctx
}

// WithBroadcastMode returns a copy of the context with an updated broadcast
// mode.
func (ctx Context) WithBroadcastMode(mode string) Context {
	ctx.BroadcastMode = mode
	return ctx
}

// WithSignModeStr returns a copy of the context with an updated SignMode
// value.
func (ctx Context) WithSignModeStr(signModeStr string) Context {
	ctx.SignModeStr = signModeStr
	return ctx
}

// WithEncodingConfig returns the context with an updated EncodingConfig
func (ctx Context) WithEncodingConfig(encodingConfig EncodingConfig) Context {
	ctx.EncodingConfig = encodingConfig
	return ctx
}

// WithBech32Prefix returns the context with an updated Bech32Prefix
func (ctx Context) WithBech32Prefix(bech32Prefix string) Context {
	ctx.Bech32Prefix = bech32Prefix
	return ctx
}

func (ctx Context) WithTxBuilder(txBuilder client.TxBuilder) Context {
	ctx.TxBuilder = txBuilder
	return ctx
}
