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
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/numeric"
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

var (
	oneAsDec    = numeric.NewDec(denominations.One)
	operations  OperationsInput
	amountFlag  string
	opType      string
	fromAddress string
	toAddress   string
	currency    string
)

// preprocessCmd represents the preprocess command
var preprocessCmd = &cobra.Command{
	Use:   "preprocess",
	Short: "Preprocess is called prior to payloads",
	Long: `Called prior to the payloads to construct a request for any metadata that is needed for transaction construction given
	If not explicitly set operation_type is NativeTransfer

	Usage:
		preprocess --from-file <path_to_operations.json> 
		preprocess --from <address> --to <address> --amount <amount> --currency <currency> --from-shard --to-shard [--type <operation_type>]  
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if fromFile != "" {
			pathToOp := fromFile
			file, err := os.Open(pathToOp)
			if err != nil {
				HandleError(err, "could not read operations file:"+pathToOp, 0)
			}
			defer file.Close()

			serializedOp, err := ioutil.ReadAll(file)
			if err != nil {
				HandleError(err, "could not serialize op:"+pathToOp, 0)
			}

			err = json.Unmarshal(serializedOp, &operations)
			if err != nil {
				HandleError(err, "could not parse operations", 0)
			}
		} else {
			for _, flagName := range []string{"from", "to", "amount", "from-shard", "to-shard"} {
				cmd.MarkFlagRequired(flagName)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		var (
			network *types.NetworkIdentifier
			ctx     = context.Background()
			meta    map[string]interface{}
		)

		if fromFile == "" {
			op, err := NewTransferOperation(fromAddress, toAddress, amountFlag, currency)
			if err != nil {
				HandleError(err, "could not create operation", 0)
			}
			operations = OperationsInput{
				Operations: op,
			}
		}

		meta = map[string]interface{}{
			"from_shard": viper.GetInt("from-shard"),
			"to_shard":   viper.GetInt("to-shard"),
		}

		network, err := GetNetwork()
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
	preprocessCmd.Flags().IntVar(&fromShard, "from-shard", 0, "From shard")
	preprocessCmd.Flags().IntVar(&toShard, "to-shard", 0, "To shard")
	preprocessCmd.Flags().StringVar(&amountFlag, "amount", "0", "amount to send")
	preprocessCmd.Flags().StringVar(&opType, "type", NativeTransferOperation, "Opration type")
	preprocessCmd.Flags().StringVar(&fromAddress, "from", "", "sender's address, keystore must exist locally")
	preprocessCmd.Flags().StringVar(&toAddress, "to", "", "the destination address")
	preprocessCmd.Flags().StringVar(&currency, "currency", "ONE", "currency used in the transaction")
	constructionCmd.AddCommand(preprocessCmd)
}
