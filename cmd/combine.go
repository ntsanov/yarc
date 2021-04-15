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
	"github.com/spf13/cobra"
)

// combineCmd represents the combine command
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "Create transaction",
	Long: ` create a transaction from an unsigned transaction 
	and an array of provided signatures

	Usage:
		combine < unsigned tx json>
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network  *types.NetworkIdentifier
			ctx      = context.Background()
			payloads PayloadsResponse
		)

		pathToTx := args[0]
		file, err := os.Open(pathToTx)
		if err != nil {
			HandleError(err, "could not read operations file:"+pathToTx, 0)
		}
		defer file.Close()
		serializedInput, err := ioutil.ReadAll(file)
		if err != nil {
			HandleError(err, "could not serialize op"+pathToTx, 0)
		}
		err = json.Unmarshal(serializedInput, &payloads)
		if err != nil {
			HandleError(err, "could not parse operations", 0)
		}

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		wtx, tx, err := DecodeWrappedTransaction(payloads.UnsignedTransaction)
		if err != nil {
			HandleError(err, "could not unmarshal unsigned transaction", 0)
		}

		signatures := []*types.Signature{}
		for _, payload := range payloads.Payloads {
			addr := payload.AccountIdentifier.Address
			// TODO put this in argument, server expects only 1 signature
			passphrase := os.Getenv("PASSPHRASE_" + addr)
			signature, err := Signature(wtx.From, passphrase, tx)
			if err != nil {
				HandleError(err, "could not sign transaction", 0)
			}
			signatures = append(signatures, signature)

		}

		signedTx, fetchErr := f.ConstructionCombine(
			ctx,
			network,
			payloads.UnsignedTransaction,
			signatures,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not fetch signed tx", 0)
		}

		PrintResult(signedTx)
	},
}

func init() {
	rootCmd.AddCommand(combineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// combineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// combineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
