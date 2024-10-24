package node

import (
	"context"
	"github.com/0xweb-3/EthCEXWallet/global"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"

	WalletTypes "github.com/0xweb-3/EthCEXWallet/wallet/types"
)

// todo 以下代码可以使用ethclient.go 直接简化处理

type EthClient interface {
	//BlockHeaderByNumber 通过块儿id获取块儿头信息
	BlockHeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	//BlockByNumber 通过块儿id获取块儿信息
	BlockByNumber(context.Context, *big.Int) (*WalletTypes.RpcBlock, error)
	//SafeBlockHeaderByNumber 最新Safe块儿的头信息
	SafeBlockHeaderByNumber(context.Context) (*types.Header, error)
	FinalizedBlockHeaderByNumber(context.Context) (*types.Header, error)
	BlockHeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	BlockHeadersByRange(ctx context.Context, start *big.Int, end *big.Int, chainId uint) ([]types.Header, error)

	TxByHash(context.Context, common.Hash) (*types.Transaction, error)
	TxReceiptByHash(context.Context, common.Hash) (*types.Receipt, error)
	// StorageHash 获取指定账户地址在特定区块的 存储根哈希值用于验证账户存储数据的完整性
	StorageHash(context.Context, common.Address, *big.Int) (common.Hash, error)
	// GetAddressNonce 获取地址的nonce
	GetAddressNonce(ctx context.Context, address common.Address) (hexutil.Uint64, error)
	// SendRawTransaction 发送交易到链上
	SendRawTransaction(ctx context.Context, rawTx string) error

	// 合约事件的监听
	FilterLogs(filterQuery ethereum.FilterQuery, chainID *big.Int) (WalletTypes.Logs, error)

	// gasPrice获取
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	// 获取当前网络上建议的 优先费
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type RPC interface {
	Close()
	CallContext(ctx context.Context, result any, method string, args ...any) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type clnt struct {
	rpc RPC
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	return rpc.BlockNumber(number.Int64()).String()
}

func (c *clnt) BlockHeaderByNumber(ctx context.Context, blockNUmber *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", toBlockNumArg(blockNUmber), false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return header, nil
}

func (c *clnt) BlockByNumber(ctx context.Context, blockNUmber *big.Int) (*WalletTypes.RpcBlock, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var block *WalletTypes.RpcBlock
	err := c.rpc.CallContext(ctx, &block, "eth_getBlockByNumber", toBlockNumArg(blockNUmber), true)
	if err != nil {
		return nil, err
	} else if block == nil {
		return nil, ethereum.NotFound
	}

	return block, nil
}

func (c *clnt) SafeBlockHeaderByNumber(ctx context.Context) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", "safe", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return header, nil
}

func (c *clnt) FinalizedBlockHeaderByNumber(ctx context.Context) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctx, &header, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return header, nil
}

func (c *clnt) BlockHeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var head *types.Header
	err := c.rpc.CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

func (c *clnt) BlockHeadersByRange(ctx context.Context, start *big.Int, end *big.Int, uint2 uint) ([]types.Header, error) {
	//TODO implement me
	panic("implement me")
}

func (c *clnt) TxByHash(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()
	var tx *types.Transaction
	err := c.rpc.CallContext(ctx, &tx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, ethereum.NotFound
	}
	return tx, nil
}

func (c *clnt) TxReceiptByHash(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	var r *types.Receipt
	err := c.rpc.CallContext(ctx, &r, "eth_getTransactionReceipt", hash)
	if err == nil && r == nil {
		return nil, ethereum.NotFound
	}
	return r, err
}

func (c *clnt) StorageHash(ctx context.Context, address common.Address, blockNumber *big.Int) (common.Hash, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	proof := struct{ StorageHash common.Hash }{}
	err := c.rpc.CallContext(ctx, &proof, "eth_getProof", address, nil, toBlockNumArg(blockNumber))
	if err != nil {
		return common.Hash{}, err
	}

	return proof.StorageHash, nil
}

func (c *clnt) GetAddressNonce(ctx context.Context, address common.Address) (hexutil.Uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var result hexutil.Uint64
	err := c.rpc.CallContext(ctx, &result, "eth_getTransactionCount", address, "pending")
	return result, err
}

func (c *clnt) SendRawTransaction(ctx context.Context, rawTx string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	return c.rpc.CallContext(ctx, nil, "eth_sendRawTransaction", rawTx)
}

func (c *clnt) FilterLogs(filterQuery ethereum.FilterQuery, chainID *big.Int) (WalletTypes.Logs, error) {
	//TODO implement me
	panic("implement me")
}

func (c *clnt) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var hex hexutil.Big
	if err := c.rpc.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *clnt) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(global.ServerConfig.MaxRequestTime))
	defer cancel()

	var hex hexutil.Big
	if err := c.rpc.CallContext(ctx, &hex, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func DailEthClient(ctx context.Context, rpcUrl string) (EthClient, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	client, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	return &clnt{
		rpc: NewRPC(client), // 使用初始化的 rpc.Client 创建 RPC 实例
	}, nil
}

type rpcClent struct {
	rpc *rpc.Client
}

func (r rpcClent) Close() {
	r.rpc.Close()
}

func (r rpcClent) CallContext(ctx context.Context, result any, method string, args ...any) error {
	err := r.rpc.CallContext(ctx, result, method, args...)
	return err
}

func (r rpcClent) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	err := r.rpc.BatchCallContext(ctx, b)
	return err
}

func NewRPC(client *rpc.Client) RPC {
	return &rpcClent{
		rpc: client,
	}
}
