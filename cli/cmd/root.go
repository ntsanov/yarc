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
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var ErrSelectedNetworkIndexOutOfRange = errors.New("selected network index out of range")

var cfgFile string

// TODO set this from viper
var ChainID = big.NewInt(2)

var (
	fromShard, toShard int
	transactionFlag    string
	blockFlag          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "harmony_cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.harmony_cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// Flags
	rootCmd.PersistentFlags().StringP("server", "s", "", "Server to query")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.BindPFlag("from_shard", preprocessCmd.Flags().Lookup("from_shard"))
	viper.BindPFlag("to_shard", preprocessCmd.Flags().Lookup("to_shard"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	// Defaults
	viper.SetDefault("server", "https://rosetta.s0.b.hmny.io")
	viper.SetDefault("timeout", 10)
	viper.SetDefault("retries", 3)
	viper.SetDefault("network_idx", 0)
	viper.SetDefault("fee_multiplier", 1)
	// Will used PASSPHRASE ENV if it is set and not used as argument
	viper.BindEnv("passphrase")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".harmony_cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".harmony_cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
