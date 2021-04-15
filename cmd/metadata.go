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

type MetaResponse struct {
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	SuggestedFee []*types.Amount        `json:"suggested_fee,omitempty"`
}

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Get any information required to construct a transaction for a specific network",
	Long: `
	Usage:
		metadata <address> <path to preprocess response>	
	
	Options need to be retrieved first by preprocess
	`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network        *types.NetworkIdentifier
			pubKey         *types.PublicKey
			preprocessResp PreprocessResponse
			ctx            = context.Background()
		)

		address := args[0]
		publicKey, err := GetPublicKey(address, "")
		if err != nil {
			HandleError(err, "could not parse public key", 0)
		}
		compressedPkey := crypto.CompressPubkey(publicKey)

		optionsPath := args[1]
		file, err := os.Open(optionsPath)
		if err != nil {
			HandleError(err, "could not read operations file:"+optionsPath, 0)
		}
		defer file.Close()

		serializedOptions, err := ioutil.ReadAll(file)
		if err != nil {
			HandleError(err, "could not serialize op:"+optionsPath, 0)
		}

		err = json.Unmarshal(serializedOptions, &preprocessResp)
		if err != nil {
			HandleError(err, "could not parse operations", 0)
		}

		pubKey = &types.PublicKey{
			Bytes:     compressedPkey,
			CurveType: types.Secp256k1,
		}

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		// TODO
		// Right now cmd expects one key, seems there might be situations where
		// multiple keys are expected, might have to change it
		keys := []*types.PublicKey{
			pubKey,
		}

		metadata, fee, fetchErr := f.ConstructionMetadata(
			ctx,
			network,
			preprocessResp.Options,
			keys,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve metadata", 0)
		}

		resp := MetaResponse{
			Metadata:     metadata,
			SuggestedFee: fee,
		}

		PrintResult(resp)
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metadataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metadataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
