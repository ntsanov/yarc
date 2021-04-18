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
	"github.com/spf13/viper"
)

var (
	passphraseFlag string
)

// combineCmd represents the combine command
var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "Create transaction",
	Long: ` create a transaction from an unsigned transaction 
	and an array of provided signatures

	If argument passphrase is not set, Environment variable PASSPHRASE will be used

	Usage:
		combine --from-file <path_to_unsigned_transaction.json> [--passphrase <passphrase>]
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("from-file")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			payloads   PayloadsResponse
			passphrase = viper.GetString("passphrase")
		)

		pathToTx := fromFile
		if len(args) > 1 {
			passphrase = args[1]
		}
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

		signature, err := Signature(wtx.From, passphrase, tx)
		if err != nil {
			HandleError(err, "could not sign transaction", 0)
		}
		signatures := []*types.Signature{signature}

		signedTx, fetchErr := f.ConstructionCombine(
			ctx,
			network,
			payloads.UnsignedTransaction,
			signatures,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not fetch signed tx", 0)
		}

		res := ParseInput{
			SignedTransaction: signedTx,
		}
		PrintResult(res)
	},
}

func init() {
	combineCmd.Flags().StringVar(&passphraseFlag, "passphrase", "", "passphrase for sender account")
	viper.BindPFlag("passphrase", combineCmd.Flags().Lookup("passphrase"))
	constructionCmd.AddCommand(combineCmd)
}
