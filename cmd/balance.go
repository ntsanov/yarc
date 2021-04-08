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
	"strconv"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cobra"
)

type BalanceResponse struct {
	BlockIdentifier *types.BlockIdentifier `json:"block_identifier"`
	Balances        []*types.Amount        `json:"balances"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
}

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get balance for an account",
	Long: `Gets an array of all AccountBalances for an AccountIdentifier and the BlockIdentifier
	at which the balance lookup was performed
	Usage: balance account_address [block]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ctx     = context.Background()
			account = &types.AccountIdentifier{
				Address: args[0],
			}
		)

		network, err := GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		var lookupBlock *types.PartialBlockIdentifier
		if len(args) > 1 {
			idx, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				HandleError(err, "could not parse block index", 0)
			}
			lookupBlock = &types.PartialBlockIdentifier{Index: &idx}
		}

		block, balances, meta, fetchErr := f.AccountBalance(
			ctx,
			network,
			account,
			lookupBlock,
			nil,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve account balance", 0)
		}

		response := BalanceResponse{
			BlockIdentifier: block,
			Balances:        balances,
			Meta:            meta,
		}

		PrintResult(response)

	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// balanceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// balanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
