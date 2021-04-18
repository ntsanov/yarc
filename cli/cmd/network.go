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

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Query network",
	Long:  `Will query networks on rosetta server - endpoints at /network`,
}

var networkCmdLs = &cobra.Command{
	Use:   "list",
	Short: "List network identifiers",
	Long: `Lists all networkIdentifiers retrieved from /network/list,
	also retrieving status from /network/status and options from /network/options`,
	Run: func(cmd *cobra.Command, args []string) {
		networks, err := ListNetworks()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}
		PrintResult(networks)
	},
}

var networkCmdOptions = &cobra.Command{
	Use:   "options",
	Short: "Version information & allowed types",
	Long: `Returns the version information and allowed network-specific types for a shard's NetworkIdentifier
	Note that not all shards support the same types/options.`,
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

		networkOpts, fetchErr := f.NetworkOptions(
			ctx,
			network,
			nil,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve network options", 0)
		}

		PrintResult(networkOpts)
	},
}

var networkCmdStatus = &cobra.Command{
	Use:   "status",
	Short: "Current status of the network",
	Long:  `Current status of the network`,
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

		networkStatus, fetchErr := f.NetworkStatus(
			ctx,
			network,
			nil,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve network options", 0)
		}

		PrintResult(networkStatus)
	},
}

func init() {
	networkCmd.AddCommand(networkCmdLs)
	networkCmd.AddCommand(networkCmdOptions)
	networkCmd.AddCommand(networkCmdStatus)
	dataCmd.AddCommand(networkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
