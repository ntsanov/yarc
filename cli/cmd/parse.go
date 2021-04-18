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
)

type ParseResponse struct {
	Operations               []*types.Operation         `json:"operations,omitempty"`
	AccountIdentifierSigners []*types.AccountIdentifier `json:"account_identifier_signers,omitempty"`
}

type ParseInput struct {
	UnsignedTransaction string                  `json:"unsigned_transaction,omitempty"`
	SignedTransaction   string                  `json:"signed_transaction,omitempty"`
	Payloads            []*types.SigningPayload `json:"payloads,omitempty"`
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Sanity check",
	Long: `Called on both unsigned and signed transactions 
	to understand the intent of the formulated transaction
	
	Usage:
		parse --from-file <path to json tx>
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("from-file")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			ctx        = context.Background()
			parseInput ParseInput
			signed     bool
			tx         string
		)
		pathToPayloads := fromFile
		filePayloads, err := os.Open(pathToPayloads)
		if err != nil {
			HandleError(err, "could not read payloads file:"+pathToPayloads, 0)
		}
		defer filePayloads.Close()
		serializedPayloads, err := ioutil.ReadAll(filePayloads)
		if err != nil {
			HandleError(err, "could not serialize payloads"+pathToPayloads, 0)
		}
		err = json.Unmarshal(serializedPayloads, &parseInput)
		if err != nil {
			HandleError(err, "could not parse operations", 0)
		}

		if parseInput.SignedTransaction != "" {
			tx = parseInput.SignedTransaction
			signed = true
		} else {
			tx = parseInput.UnsignedTransaction
		}

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		// TODO add flag weather it is signed or not
		operations, accounts, _, fetchErr := f.ConstructionParse(
			ctx,
			network,
			signed,
			tx,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not retrieve parse result", 0)
		}

		resp := ParseResponse{
			Operations:               operations,
			AccountIdentifierSigners: accounts,
		}

		PrintResult(resp)

	},
}

func init() {
	constructionCmd.AddCommand(parseCmd)
}
