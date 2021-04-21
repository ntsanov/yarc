package cmd

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fatih/color"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/store"
	hmyTypes "github.com/harmony-one/harmony/core/types"
	"github.com/harmony-one/harmony/numeric"
	stakingTypes "github.com/harmony-one/harmony/staking/types"

	"github.com/spf13/viper"
)

type ErrorResponse struct {
	Err  string `json:"err,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Code int    `json:"code,omitempty"`
}

// HandleError pretty prints the error and exits setting the exit code to something negative
func HandleError(err error, message string, code int) {
	resp := ErrorResponse{
		Err:  err.Error(),
		Msg:  message,
		Code: code,
	}
	out, err := json.Marshal(resp)
	if err != nil {
		color.New(color.BgRed).Println("Something went really wrong")
	}
	color.New(color.BgRed).Println(string(out))
	os.Exit(-1)
}

// PrintResult should be called only on success
// it does not exit so that should be handled by the cmd
func PrintResult(res interface{}) {
	out, err := json.Marshal(res)
	if err != nil {
		HandleError(err, "could not serialize result", 0)
	}
	color.New(color.FgGreen).Println(string(out))
}

// ListNetworks lists available network on the server
func ListNetworks() (networks *types.NetworkListResponse, err error) {
	var (
		nodeURL string = viper.GetString("node")
	)
	ctx := context.Background()
	fetcherOpts := []fetcher.Option{
		fetcher.WithTimeout(time.Duration(viper.GetInt("timeout")) * time.Second),
	}
	f := fetcher.New(
		nodeURL,
		fetcherOpts...,
	)
	networks, fetchErr := f.NetworkList(ctx, nil)
	if fetchErr != nil {
		return nil, fetchErr.Err
	}
	return networks, nil
}

// GetNetwork returns network to be used. Network is chosen by
// the following order.
// 1. Specific network set in config
// 2. Network index set in config
// 3. Default network index(0) is used
func GetNetwork() (network *types.NetworkIdentifier, err error) {
	configNetwork := viper.Get("network")
	// No network is specifically set in config, use idx set in config

	if configNetwork == nil {
		idx := viper.GetInt("network_idx")
		networks, err := ListNetworks()
		if err != nil {
			return nil, err
		}
		if len(networks.NetworkIdentifiers) < idx {
			return nil, ErrSelectedNetworkIndexOutOfRange
		}
		network = networks.NetworkIdentifiers[idx]
	} else {
		network = configNetwork.(*types.NetworkIdentifier)
	}
	return network, nil
}

func NewTransferOperation(from, to, amount, currency string) ([]*types.Operation, error) {
	amt, err := numeric.NewDecFromStr(amount)
	if err != nil {
		return nil, err
	}
	amountInOne := amt.Mul(oneAsDec)
	return []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: opType,
			Account: &types.AccountIdentifier{
				Address: from,
			},
			Amount: &types.Amount{
				Value: amountInOne.Neg().RoundInt().String(),
				Currency: &types.Currency{
					Symbol:   currency,
					Decimals: 18,
				},
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*types.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Type: opType,
			Account: &types.AccountIdentifier{
				Address: to,
			},
			Amount: &types.Amount{
				Value: amountInOne.RoundInt().String(),
				Currency: &types.Currency{
					Symbol:   currency,
					Decimals: 18,
				},
			},
		},
	}, nil
}

func NewFetcher(ctx context.Context, network *types.NetworkIdentifier) (*fetcher.Fetcher, error) {

	var nodeURL string = viper.GetString("node")

	fetcherOpts := []fetcher.Option{
		fetcher.WithTimeout(time.Duration(viper.GetInt("timeout")) * time.Second),
	}

	f := fetcher.New(
		nodeURL,
		fetcherOpts...,
	)

	_, _, fetchErr := f.InitializeAsserter(ctx, network)
	if fetchErr != nil {
		return nil, fetchErr.Err
	}

	return f, nil

}

func DecodeWrappedTransaction(marshaledStr string) (*WrappedTransaction, hmyTypes.PoolTransaction, error) {
	wt := &WrappedTransaction{}
	err := json.Unmarshal([]byte(marshaledStr), wt)
	var tx hmyTypes.PoolTransaction
	stream := rlp.NewStream(bytes.NewBuffer(wt.RLPBytes), 0)
	if wt.IsStaking {
		stakingTx := &stakingTypes.StakingTransaction{}
		if err := stakingTx.DecodeRLP(stream); err != nil {
			return nil, nil, errors.New("rlp encoding error for staking transaction")
		}
		tx = stakingTx
	} else {
		plainTx := &hmyTypes.Transaction{}
		if err := plainTx.DecodeRLP(stream); err != nil {
			return nil, nil, errors.New("rlp encoding error for plain transaction")
		}
		tx = plainTx
	}
	return wt, tx, err
}

// TODO make this a factory and set factory at start from flag/config
func GetKeys(address, passphrase string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	pkStr := viper.GetString("private_key")
	if pkStr == "" {
		ks, acct, err := store.UnlockedKeystore(address, passphrase)
		if err != nil {
			return nil, nil, err
		}
		_, key, err := ks.GetDecryptedKey(*acct, passphrase)
		if err != nil {
			return nil, nil, err
		}
		publicKey := key.PrivateKey.Public()
		return key.PrivateKey, publicKey.(*ecdsa.PublicKey), nil
	}
	privateKeyBytes, err := hex.DecodeString(pkStr)
	if err != nil {
		return nil, nil, err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return nil, nil, common.ErrBadKeyLength
	}
	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	return privateKey.ToECDSA(), publicKey.ToECDSA(), nil

}

// Signature signs transaction and returns Signature ready to be used
func Signature(account *types.AccountIdentifier, passphrase string, tx hmyTypes.PoolTransaction) (*types.Signature, error) {
	// ks, acct, err := store.UnlockedKeystore(account.Address, passphrase)
	// if err != nil {
	// 	return nil, err
	// }
	// _, key, err := ks.GetDecryptedKey(*acct, passphrase)
	// if err != nil {
	// 	return nil, err
	// }
	// publicKey := key.PrivateKey.Public()
	privateKey, publicKey, err := GetKeys(account.Address, passphrase)
	if err != nil {
		return nil, err
	}
	compressedPublicKey := crypto.CompressPubkey(publicKey)

	signature := types.Signature{
		PublicKey: &types.PublicKey{
			Bytes:     compressedPublicKey,
			CurveType: types.Secp256k1,
		},
		SigningPayload: &types.SigningPayload{
			AccountIdentifier: account,
			SignatureType:     types.EcdsaRecovery,
		},
		SignatureType: types.EcdsaRecovery,
	}
	var ChainID = big.NewInt(viper.GetInt64("chain_id"))
	switch orgTx := tx.(type) {
	case *stakingTypes.StakingTransaction:
		signer := stakingTypes.NewEIP155Signer(ChainID)
		signature.SigningPayload.Bytes = signer.Hash(orgTx).Bytes()
	case *hmyTypes.Transaction:
		signer := hmyTypes.NewEIP155Signer(ChainID)
		signature.SigningPayload.Bytes = signer.Hash(orgTx).Bytes()
	default:
		return nil, errors.New("unknown transaction type")
	}
	// signedPayload, err := ks.SignHash(*acct, signature.SigningPayload.Bytes)
	signedPayload, err := crypto.Sign(signature.SigningPayload.Bytes, privateKey)
	if err != nil {
		return nil, err
	}
	signature.Bytes = signedPayload
	return &signature, nil
}

// func GetPublicKey(address, passphrase string) (*ecdsa.PublicKey, error) {
// 	ks, acct, err := store.UnlockedKeystore(address, passphrase)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, key, err := ks.GetDecryptedKey(*acct, passphrase)
// 	if err != nil {
// 		return nil, err
// 	}
// 	publicKey := key.PrivateKey.Public()
// 	return publicKey.(*ecdsa.PublicKey), nil
// }
