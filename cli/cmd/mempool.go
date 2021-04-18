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

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cobra"
)

type MempoolResponse struct {
	TransactionIdentifiers []*types.TransactionIdentifier `json:"transaction_identifiers"`
}

type TransactionResponse struct {
	Transaction *types.Transaction
}

// mempoolCmd represents the mempool command
var mempoolCmd = &cobra.Command{
	Use:   "mempool",
	Short: "Gets all Transaction Identifiers in the mempool",
	Long: `
	Usage:
		mempool [--transaction <tx_hash>]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network *types.NetworkIdentifier
			ctx     = context.Background()
		)

		network, err := GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		if transactionFlag != "" {
			tx, _, fetchErr := f.MempoolTransaction(
				ctx,
				network,
				&types.TransactionIdentifier{Hash: transactionFlag},
			)
			if fetchErr != nil {
				HandleError(fetchErr.Err, "could not fetch block", 0)
			}

			PrintResult(TransactionResponse{
				Transaction: tx,
			})

		} else {
			mempool, fetchErr := f.Mempool(
				ctx,
				network,
			)

			if fetchErr != nil {
				HandleError(fetchErr.Err, "could not fetch block", 0)
			}

			PrintResult(MempoolResponse{
				TransactionIdentifiers: mempool,
			})

		}

	},
}

func init() {
	mempoolCmd.Flags().StringVar(&transactionFlag, "transaction", "", "transaction hash")
	dataCmd.AddCommand(mempoolCmd)
}
