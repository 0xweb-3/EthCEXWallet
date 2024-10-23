package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// 构建ERC-20的交易数据
func BuildErc20Data(toAddress common.Address, amount *big.Int) ([]byte, error) {
	var data []byte

	// 获取 transfer 函数的签名哈希
	//transferFnSignature := []byte("transfer(address,uint256)")
	//hash := sha3.NewLegacyKeccak256()
	//hash.Write(transferFnSignature)
	//methodID := hash.Sum(nil)[:4]
	methodID := crypto.Keccak256Hash([]byte("transfer(address,uint256)")).Bytes()[:4]

	// 对 toAddress 和 amount 进行 ABI 编码
	dataAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// 组合交易数据
	data = append(data, methodID...)
	data = append(data, dataAddress...)
	data = append(data, dataAmount...)
	return data, nil
}

// BuildErc721Data 构建 ERC-721 合约的 safeTransferFrom 交易数据
func BuildErc721Data(fromAddress, toAddress common.Address, tokenID *big.Int) ([]byte, error) {
	var data []byte

	// 获取 safeTransferFrom(address,address,uint256) 函数的签名哈希
	methodID := crypto.Keccak256Hash([]byte("safeTransferFrom(address,address,uint256)")).Bytes()[:4]

	// 对 fromAddress, toAddress 和 tokenID 进行 ABI 编码
	dataFromAddress := common.LeftPadBytes(fromAddress.Bytes(), 32)
	dataToAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	dataTokenID := common.LeftPadBytes(tokenID.Bytes(), 32)

	// 组合交易数据
	data = append(data, methodID...)
	data = append(data, dataFromAddress...)
	data = append(data, dataToAddress...)
	data = append(data, dataTokenID...)

	return data, nil
}

// 使用EIP1559的方式实现交易
func OfflineSignTx(feeTx *types.DynamicFeeTx, privateKeyStr string, chainId *big.Int) (string, error) {
	// 将私钥字符串转换为 ECDSA 私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	// 构造交易
	tx := types.NewTx(feeTx)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	if err != nil {
		return "", err
	}

	// 将签名后的交易编码为字节流
	signedTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", err
	}

	// 返回签名后的交易，转换为十六进制字符串
	return common.Bytes2Hex(signedTxBytes), nil
}
