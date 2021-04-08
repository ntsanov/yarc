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

// mempoolCmd represents the mempool command
var mempoolCmd = &cobra.Command{
	Use:   "mempool",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	},
}

func init() {
	rootCmd.AddCommand(mempoolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mempoolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mempoolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
