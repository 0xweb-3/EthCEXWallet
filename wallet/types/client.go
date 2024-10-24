package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionList struct {
	To   string `json:"to"`
	Hash string `json:"hash"`
}

type RpcBlock struct {
	Hash         common.Hash       `json:"hash"`
	Transactions []TransactionList `json:"transactions"`
	BaseFee      string            `json:"baseFeePerGas"`
}

type Logs struct {
	Logs        []types.Log   `json:"logs"`
	BlockHeader *types.Header `json:"block_header"`
}
