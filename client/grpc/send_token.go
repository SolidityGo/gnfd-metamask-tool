package client

import (
	"fmt"
	"github.com/bnb-chain/gnfd-go-sdk/types"
	"github.com/bnb-chain/gnfd-go-sdk/util"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/grpc"
)

func (c *GreenfieldClient) SendToken(req types.SendTokenRequest, txOpt *types.TxOption, opts ...grpc.CallOption) (*types.TxBroadcastResponse, error) {
	if err := util.ValidateToken(req.Token); err != nil {
		return nil, err
	}
	if err := util.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromHexUnsafe(req.ToAddress)
	if err != nil {
		return nil, err
	}
	km, err := c.GetKeyManager()
	if err != nil {
		return nil, err
	}
	transferMsg := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(req.Token, req.Amount)))
	res, err := c.BroadcastTx([]sdk.Msg{transferMsg}, txOpt, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *GreenfieldClient) GenerateTx(signature string, req types.SendTokenRequest) ([]byte, error) {
	if err := util.ValidateToken(req.Token); err != nil {
		return nil, err
	}
	if err := util.ValidateAmount(req.Amount); err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromHexUnsafe(req.ToAddress)
	if err != nil {
		return nil, err
	}
	km, err := c.GetKeyManager()
	if err != nil {
		return nil, err
	}
	transferMsg := banktypes.NewMsgSend(km.GetAddr(), to, sdk.NewCoins(sdk.NewInt64Coin(req.Token, req.Amount)))

	txConfig := authtx.NewTxConfig(types.Cdc(), authtx.DefaultSignModes)
	txBuilder := txConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(transferMsg)
	if err != nil {
		return nil, err
	}
	txBuilder.SetGasLimit(types.DefaultGasLimit)

	address := km.GetAddr().String()
	account, err := c.Account(address)
	if err != nil {
		return nil, err
	}
	accountNum := account.GetAccountNumber()
	accountSeq := account.GetSequence()

	sig := signing.SignatureV2{
		PubKey: km.GetPrivKey().PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_EIP_712,
			Signature: nil,
		},
		Sequence: accountSeq,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	sig = signing.SignatureV2{}

	signerData := xauthsigning.SignerData{
		ChainID:       types.ChainId,
		AccountNumber: accountNum,
		Sequence:      accountSeq,
	}

	sig, err = clitx.SignWithPrivKey(signing.SignMode_SIGN_MODE_EIP_712,
		signerData,
		txBuilder,
		km.GetPrivKey(),
		txConfig,
		accountSeq,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("------------------")
	fmt.Println("generate sig from go-sdk", hexutil.Encode(sig.Data.(*signing.SingleSignatureData).Signature))
	fmt.Println("------------------")

	sig.Data.(*signing.SingleSignatureData).Signature, err = hexutil.Decode(signature)
	if err != nil {
		return nil, err
	}

	fmt.Println("modified sig on go-sdk", hexutil.Encode(sig.Data.(*signing.SingleSignatureData).Signature))

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, err
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil

	/*	txRes, err := c.TxClient.BroadcastTx(
			context.Background(),
			&tx.BroadcastTxRequest{
				Mode:    mode,
				TxBytes: txBytes,
			})

		if err != nil {
			return nil, err
		}
		txResponse := txRes.TxResponse

		return &types.TxBroadcastResponse{
			Ok:     txResponse.Code == 0,
			Log:    txRes.TxResponse.RawLog,
			TxHash: txResponse.TxHash,
			Code:   txResponse.Code,
			Data:   txResponse.Data,
		}, nil
	*/
}
