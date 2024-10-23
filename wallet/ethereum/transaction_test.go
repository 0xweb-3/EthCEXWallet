package ethereum

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

/*
*
privateKey 17a01d2d0862c190dd3d286f5233039938c0522da31fd7d580569cdc07e642f4
publicKey a2e624a9cba7cb4b6b2814c8535b11a130f4818027f125c2ecdf23da303491bb822726e209563a1046b20d7dccfebd3857d4b5e9e3d079dcfbf4c9737fc06d18
address 0xEB80a127b2b763C631D8ADCeBb0976b190C8C227
*/
func TestOfflineSignTx(t *testing.T) {
	client, err := ethclient.Dial("https://sepolia.drpc.org") // Sepolia 测试网
	if err != nil {
		t.Error(err)
	}

	privateKeyStr := "17a01d2d0862c190dd3d286f5233039938c0522da31fd7d580569cdc07e642f4"
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		t.Error(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Log("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	amount := new(big.Int)
	amount.SetString("1000000000000000", 10)

	// 获得帐户nonce,回你应该使用的下一个nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		t.Error(err)
	}
	t.Log("nonce", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("chainID", chainID)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("gasPrice", gasPrice)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if err != nil {
		t.Error(err)
	}
	t.Log("gasLimit", gasLimit)

	toAddress := common.HexToAddress("0x8ff44C9b5Eab5E5CE8d1d642184b70e9b9587F74")
	t.Log(toAddress)

	// 这是给矿工的“小费”，通常用来加快交易的处理。
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log("gasTipCap", gasTipCap)

	// 这是你愿意为每单位 gas 支付的最大费用。
	gasFeeCap := new(big.Int).Add(gasPrice, gasTipCap)
	t.Log("gasFeeCap", gasFeeCap)

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     amount,
		Data:      nil,
	}

	txHex, err := OfflineSignTx(tx, privateKeyStr, chainID)
	if err != nil {
		t.Error(err)
	}
	t.Log("Signed Transaction Hex: ", txHex)

	// 解码签名后的交易十六进制
	signedTx := new(types.Transaction)
	signedTxBytes := common.Hex2Bytes(txHex)
	if err != nil {
		t.Error(err)
	}
	err = signedTx.UnmarshalBinary(signedTxBytes)
	if err != nil {
		t.Error(err)
	}

	// 发送签名好的交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Transaction sent: %s", signedTx.Hash().Hex())
}
