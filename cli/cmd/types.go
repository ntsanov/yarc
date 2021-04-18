package cmd

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// WrappedTransaction is a wrapper for a transaction that includes all relevant
// data to parse a transaction.
type WrappedTransaction struct {
	RLPBytes     []byte                   `json:"rlp_bytes"`
	IsStaking    bool                     `json:"is_staking"`
	ContractCode hexutil.Bytes            `json:"contract_code"`
	From         *types.AccountIdentifier `json:"from"`
}
