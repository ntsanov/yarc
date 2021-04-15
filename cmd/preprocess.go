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

type OperationsInput struct {
	Operations []*types.Operation `json:"operations,omitempty"`
}

type PreprocessResponse struct {
	Options            map[string]interface{}     `json:"options,omitempty"`
	RequiredPublicKeys []*types.AccountIdentifier `json:"required_public_keys,omitempty"`
}

// preprocessCmd represents the preprocess command
var preprocessCmd = &cobra.Command{
	Use:   "preprocess",
	Short: "Preprocess is called prior to payloads",
	Long: `Called prior to the payloads to construct a request for any metadata that is needed for transaction construction given
	Usage:
		preprocess <path to operations in .json format> 
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			operations OperationsInput
			meta       map[string]interface{}
		)

		pathToOp := args[0]
		file, err := os.Open(pathToOp)
		if err != nil {
			HandleError(err, "could not read operations file:"+pathToOp, 0)
		}
		defer file.Close()

		serializedOp, err := ioutil.ReadAll(file)
		if err != nil {
			HandleError(err, "could not serialize op"+pathToOp, 0)
		}

		err = json.Unmarshal(serializedOp, &operations)
		if err != nil {
			HandleError(err, "could not parse operations", 0)
		}

		// TEMPORARY - include flags
		meta = map[string]interface{}{
			"from_shard": viper.GetInt("from_shard"),
			"to_shard":   viper.GetInt("to_shard"),
		}
		// TEMPORARY - include flags

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		opts, accounts, fetchErr := f.ConstructionPreprocess(
			ctx,
			network,
			operations.Operations,
			meta,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not preprocess operations", 0)
		}

		resp := PreprocessResponse{
			Options:            opts,
			RequiredPublicKeys: accounts,
		}

		PrintResult(resp)
	},
}

func init() {
	rootCmd.AddCommand(preprocessCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// preprocessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// preprocessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
