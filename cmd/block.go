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

// blockCmd represents the block command
var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Gets a block by its BlockIdentifier",
	Long: `Usage:
	block block_id [transaction_id]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network *types.NetworkIdentifier
			ctx     = context.Background()
		)

		blockID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			HandleError(err, "could not parse block id", 0)
		}

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		block, fetchErr := f.Block(
			ctx,
			network,
			&types.PartialBlockIdentifier{
				Index: &blockID,
			},
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not fetch block", 0)
		}

		PrintResult(block)
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
