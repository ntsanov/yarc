/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PayloadsResponse struct {
	UnsignedTransaction string                  `json:"unsigned_transaction,omitempty"`
	Payloads            []*types.SigningPayload `json:"payloads,omitempty"`
}

var (
	metaFlag string
	metaResp MetaResponse
)

// payloadsCmd represents the payloads command
var payloadsCmd = &cobra.Command{
	Use:   "payloads",
	Short: "A brief description of your command",
	Long: `Payloads is called with an array of operations and the response from metadata
	Usage:
		payloads --meta <path_to_metadata.json> --from-file <path_to_oprations.json> 
		payloads --meta <path_to_metadata.json> --from <address> --to <address> --amount <amount> --currency <currency> --from-shard <idx> --to-shard <idx> [--type <operation_type>]
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fileMeta, err := os.Open(metaFlag)
		if err != nil {
			HandleError(err, "could not read meta file:"+metaFlag, 0)
		}
		defer fileMeta.Close()
		serializedMetaResp, err := ioutil.ReadAll(fileMeta)
		if err != nil {
			HandleError(err, "could not serialize op"+metaFlag, 0)
		}
		err = json.Unmarshal(serializedMetaResp, &metaResp)
		if err != nil {
			HandleError(err, "could not parse metadata", 0)
		}

		if fromFile != "" {
			pathToOp := fromFile
			file, err := os.Open(pathToOp)
			if err != nil {
				HandleError(err, "could not read operations file:"+pathToOp, 0)
			}
			defer file.Close()

			serializedOp, err := ioutil.ReadAll(file)
			if err != nil {
				HandleError(err, "could not serialize op:"+pathToOp, 0)
			}

			err = json.Unmarshal(serializedOp, &operations)
			if err != nil {
				HandleError(err, "could not parse operations", 0)
			}
		} else {
			for _, flagName := range []string{"from", "to", "amount", "from-shard", "to-shard"} {
				cmd.MarkFlagRequired(flagName)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			operations OperationsInput
			passphrase = viper.GetString("passphrase")
		)

		address := fromAddress
		_, publicKey, err := GetKeys(address, passphrase)
		if err != nil {
			HandleError(err, "could not parse public key", 0)
		}
		compressedPkey := crypto.CompressPubkey(publicKey)
		if fromFile != "" {
			pathToOp := fromFile
			fileOp, err := os.Open(pathToOp)
			if err != nil {
				HandleError(err, "could not read operations file:"+pathToOp, 0)
			}
			defer fileOp.Close()
			serializedOp, err := ioutil.ReadAll(fileOp)
			if err != nil {
				HandleError(err, "could not serialize op"+pathToOp, 0)
			}
			err = json.Unmarshal(serializedOp, &operations)
			if err != nil {
				HandleError(err, "could not parse operations", 0)
			}
		} else {
			op, err := NewTransferOperation(fromAddress, toAddress, amountFlag, currency)
			if err != nil {
				HandleError(err, "could not create operation", 0)
			}
			operations = OperationsInput{
				Operations: op,
			}
		}
		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		keys := []*types.PublicKey{
			{
				Bytes:     compressedPkey,
				CurveType: types.Secp256k1,
			},
		}

		unsignedTx, payloads, fetchErr := f.ConstructionPayloads(
			ctx,
			network,
			operations.Operations,
			metaResp.Metadata,
			keys,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve metadata", 0)
		}

		resp := PayloadsResponse{
			UnsignedTransaction: unsignedTx,
			Payloads:            payloads,
		}

		PrintResult(resp)
	},
}

func init() {
	payloadsCmd.Flags().StringVar(&passphraseFlag, "passphrase", "", "passphrase for sender account")
	viper.BindPFlag("passphrase", payloadsCmd.Flags().Lookup("passphrase"))
	payloadsCmd.Flags().IntVar(&fromShard, "from-shard", 0, "From shard")
	payloadsCmd.Flags().IntVar(&toShard, "to-shard", 0, "To shard")
	payloadsCmd.Flags().StringVar(&amountFlag, "amount", "0", "amount to send")
	payloadsCmd.Flags().StringVar(&opType, "type", NativeTransferOperation, "Opration type")
	payloadsCmd.Flags().StringVar(&fromAddress, "from", "", "sender's address, keystore must exist locally")
	payloadsCmd.Flags().StringVar(&toAddress, "to", "", "the destination address")
	payloadsCmd.Flags().StringVar(&currency, "currency", "ONE", "currency used in the transaction")
	payloadsCmd.Flags().StringVar(&metaFlag, "meta", "", "path to metadata file in json format")
	payloadsCmd.MarkFlagRequired("meta")
	constructionCmd.AddCommand(payloadsCmd)
}
