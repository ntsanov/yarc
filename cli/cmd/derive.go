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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deriveCmd represents the derive command
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "Returns the AccountIdentifier associated with a public key",
	Long: `
	Usage:
		derive <address> [--passphrase]`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			network    *types.NetworkIdentifier
			pubKey     *types.PublicKey
			ctx        = context.Background()
			passphrase = viper.GetString("passphrase")
		)

		address := args[0]
		_, publicKey, err := GetKeys(address, passphrase)
		compressedPkey := crypto.CompressPubkey(publicKey)
		// fmt.Println(hex.EncodeToString(compressedPkey))

		if err != nil {
			HandleError(err, "could not parse public key", 0)
		}

		// harmony only uses secp256k1
		// It will be better to set curvetype with a flag
		// to be able to work with other implementations
		pubKey = &types.PublicKey{
			Bytes:     compressedPkey,
			CurveType: types.Secp256k1,
		}

		network, err = GetNetwork()
		if err != nil {
			HandleError(err, "could not retrieve networks", 0)
		}

		f, err := NewFetcher(ctx, network)
		if err != nil {
			HandleError(err, "could not create fetcher", 0)
		}

		account, _, fetchErr := f.ConstructionDerive(
			ctx,
			network,
			pubKey,
			nil,
		)

		if fetchErr != nil {
			HandleError(fetchErr.Err, "could not derive account from public key", 0)
		}

		PrintResult(account)

	},
}

func init() {

	deriveCmd.Flags().StringVar(&passphraseFlag, "passphrase", "", "passphrase for sender account")
	viper.BindPFlag("passphrase", deriveCmd.Flags().Lookup("passphrase"))
	constructionCmd.AddCommand(deriveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deriveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deriveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
