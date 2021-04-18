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
	"errors"
	"io/ioutil"
	"os"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submits a pre-signed transaction to the node without waiting for the transaction to be finalized.",
	Long: `
	Usage:
		submit --from-file <signed_transaction.json>
		submit <transaction string>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			parseInput ParseInput
			tx         string
		)
		if fromFile != "" {
			inputFile, err := os.Open(fromFile)
			if err != nil {
				HandleError(err, "could not read signed transaction file:"+fromFile, 0)
			}
			defer inputFile.Close()
			serializedInput, err := ioutil.ReadAll(inputFile)
			if err != nil {
				HandleError(err, "could not read serialize transaction"+fromFile, 0)
			}
			err = json.Unmarshal(serializedInput, &parseInput)
			if err != nil {
				HandleError(err, "could not parse transaction file", 0)
			}
			tx = parseInput.SignedTransaction
		} else {
			if len(args) < 1 {
				HandleError(errors.New("transaction missing"), "must provide transaction as a string or file", 0)
			}
			tx = args[0]
		}

		network, err := GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		txIdentifier, metadata, fetchErr := f.ConstructionSubmit(
			ctx,
			network,
			tx,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve transaction hash", 0)
		}

		resp := TransactionIdentifierResponse{
			TransactionIdentifier: txIdentifier,
			Metadata:              metadata,
		}

		PrintResult(resp)

	},
}

func init() {
	constructionCmd.AddCommand(submitCmd)
}
