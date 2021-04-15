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
)

type PayloadsResponse struct {
	UnsignedTransaction string                  `json:"unsigned_transaction,omitempty"`
	Payloads            []*types.SigningPayload `json:"payloads,omitempty"`
}

// payloadsCmd represents the payloads command
var payloadsCmd = &cobra.Command{
	Use:   "payloads",
	Short: "A brief description of your command",
	Long: `Payloads is called with an array of operations and the response from metadata
	Usage:

		payloads <address> <path to operations json> <path to metadata json>
	`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			operations OperationsInput
			metaResp   MetaResponse
		)

		address := args[0]
		publicKey, err := GetPublicKey(address, "")
		if err != nil {
			HandleError(err, "could not parse public key", 0)
		}
		compressedPkey := crypto.CompressPubkey(publicKey)
		pathToOp := args[1]
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

		pathToMetadata := args[2]
		fileMeta, err := os.Open(pathToMetadata)
		if err != nil {
			HandleError(err, "could not read meta file:"+pathToOp, 0)
		}
		defer fileMeta.Close()
		serializedMetaResp, err := ioutil.ReadAll(fileMeta)
		if err != nil {
			HandleError(err, "could not serialize op"+pathToOp, 0)
		}
		err = json.Unmarshal(serializedMetaResp, &metaResp)
		if err != nil {
			HandleError(err, "could not parse metadata", 0)
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
	rootCmd.AddCommand(payloadsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// payloadsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// payloadsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
