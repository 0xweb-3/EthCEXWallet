package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/0xweb-3/EthCEXWallet/wallet/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// CreateAddressByPrivateKey 通过私钥生成地址信息
func CreateAddressByPrivateKey(ctx context.Context, privateKey *ecdsa.PrivateKey) (*types.EthAddress, error) {
	ethAddress := &types.EthAddress{}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	//转换为十六进制字符串,删除“0x”
	ethAddress.PrivateKey = hexutil.Encode(privateKeyBytes)[2:] // 一般不返回私钥

	// 从私钥转换为公钥
	publicKey := privateKey.Public()

	// 将其转换为十六进制，也去除0x
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	ethAddress.PublicKey = hexutil.Encode(publicKeyBytes)[4:]

	// 生成地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	ethAddress.Address = address

	return ethAddress, nil
}

// CreateAddress 直接生成新的钱包
func CreateAddress(ctx context.Context) (*types.EthAddress, error) {
	//生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return CreateAddressByPrivateKey(ctx, privateKey)
}

// GetAddressByPublicKey 将公钥转换为地址信息
func GetAddressByPublicKey(ctx context.Context, publicKey string) (string, error) {
	// 解码传入的公钥字符串为字节数组
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", errors.New("invalid public key format")
	}

	// 如果公钥长度为 64 字节，自动补上未压缩公钥的前缀 `04`
	if len(pubKeyBytes) == 64 {
		pubKeyBytes = append([]byte{0x04}, pubKeyBytes...)
	}

	var pubKey *ecdsa.PublicKey

	// 检查公钥的前缀以确定是压缩还是未压缩格式
	switch len(pubKeyBytes) {
	case 33:
		// 压缩公钥：33 字节，前缀应为 02 或 03
		pubKey, err = crypto.DecompressPubkey(pubKeyBytes)
		if err != nil {
			return "", errors.New("failed to decompress public key")
		}
	case 65:
		// 未压缩公钥：65 字节，前缀应为 04
		pubKey, err = crypto.UnmarshalPubkey(pubKeyBytes)
		if err != nil {
			return "", errors.New("invalid uncompressed public key")
		}
	default:
		return "", errors.New("invalid public key length")
	}

	// 使用公钥生成以太坊地址
	address := crypto.PubkeyToAddress(*pubKey)

	// 返回地址的十六进制表示
	return address.Hex(), nil
}
