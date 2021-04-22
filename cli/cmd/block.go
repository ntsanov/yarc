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
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/spf13/cobra"
)

var (
	blockIdentifier *types.PartialBlockIdentifier
)

// blockCmd represents the block command
var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Gets a block by its BlockIdentifier",
	Args:  cobra.MinimumNArgs(1),
	Long: `
	Usage:

	block [<block hash or id>] [--transaction <tx_hash>]`,
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

		if len(args) > 0 {
			block := args[0]
			blockIdentifier = &types.PartialBlockIdentifier{}
			if strings.HasPrefix(block, "0x") {
				blockIdentifier.Hash = &block
			} else {
				blockIdx, err := strconv.ParseInt(block, 10, 64)
				if err != nil {
					HandleError(err, "could not parse block id", 0)
				}
				blockIdentifier.Index = &blockIdx
			}
		}

		block, fetchErr := f.Block(
			ctx,
			network,
			blockIdentifier,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not fetch block", 0)
		}

		PrintResult(block)
	},
}

func init() {
	blockCmd.MarkFlagRequired("block")
	dataCmd.AddCommand(blockCmd)
}
