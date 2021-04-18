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
	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Provides the ability to access blocks, transactions, and balances",
	Long: `    Caller                                              + Data API Implementation
                                 +----------------------------------------------------------------------------------+
                               X                                                      |
                               X      Get Supported Networks +---------------------------------> /network/list
                               X                                                      |                +
                               X  +-----------+--------------------------------------------------------+
      Get supported networks,  X  |           v                                       |
      their supported options, X  |   Get Network Options +------------------------------------> /network/options
      and their status         X  |                                                   |
                               X  +-----------+                                       |
                               X              v                                       |
                               X      Get Network Status +-------------------------------------> /network/status
                               X                                                      |
                                                                                      |
                                   X                                                  |
                                   X  Get Block +----------------------------------------------> /block
                                   X                                                  |             +
                        +---------+X                +-----------------------------------------------+
                        |          X                v                                 |
Ensure balance computed |          X  [Optional] Get Additional Block Transactions +-----------> /block/transaction
from block operations   |          X                                                  |
equals balance on node  |                                                             |
                        |                                                             |
                        +-----------> Get account balance for each +---------------------------> /account/balance
                                      account seen in a block                         |
                                                                                      |
                                                                                      |
                                   X                                                  |
                                   X  Get Mempool Transactions +-------------------------------> /mempool
      Monitor the mempool for      X                                                  |              +
      broadcast transaction status X                 +-----------------------------------------------+
      and incoming deposits        X                 v                                |
                                   X  Get a Specific Mempool Transaction +---------------------> /mempool/transaction
                                   X                                                  |
                                                                                      +`,
}

func init() {
	rootCmd.AddCommand(dataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
