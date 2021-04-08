package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type ErrorResponse struct {
	Err  string `json:"err,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Code int    `json:"code,omitempty"`
}

// HandleError pretty prints the error and exits setting the exit code to something negative
func HandleError(err error, message string, code int) {
	resp := ErrorResponse{
		Err:  err.Error(),
		Msg:  message,
		Code: code,
	}
	out, err := json.Marshal(resp)
	if err != nil {
		color.New(color.BgRed).Println("Something went really wrong")
	}
	color.New(color.BgRed).Println(string(out))
	os.Exit(-1)
}

// PrintResult should be called only on success
// it does not exit so that should be handled by the cmd
func PrintResult(res interface{}) {
	out, err := json.Marshal(res)
	if err != nil {
		HandleError(err, "could not serialize result", 0)
	}
	color.New(color.FgGreen).Println(string(out))
}

// ListNetworks lists available network on the server
func ListNetworks() (networks *types.NetworkListResponse, err error) {
	var (
		nodeURL string = viper.GetString("node_url")
	)
	ctx := context.Background()
	fetcherOpts := []fetcher.Option{
		fetcher.WithTimeout(time.Duration(viper.GetInt("timeout")) * time.Second),
	}
	f := fetcher.New(
		nodeURL,
		fetcherOpts...,
	)
	networks, fetchErr := f.NetworkList(ctx, nil)
	if fetchErr != nil {
		return nil, fetchErr.Err
	}
	return networks, nil
}

// GetNetwork returns network to be used. Network is chosen by
// the following order.
// 1. Specific network set in config
// 2. Network index set in config
// 3. Default network index(0) is used
func GetNetwork() (network *types.NetworkIdentifier, err error) {
	configNetwork := viper.Get("network")
	// No network is specifically set in config, use idx set in config

	if configNetwork == nil {
		idx := viper.GetInt("network_idx")
		networks, err := ListNetworks()
		if err != nil {
			return nil, err
		}
		if len(networks.NetworkIdentifiers) < idx {
			return nil, ErrSelectedNetworkIndexOutOfRange
		}
		network = networks.NetworkIdentifiers[idx]
	} else {
		network = configNetwork.(*types.NetworkIdentifier)
	}
	return network, nil
}

func NewFetcher(ctx context.Context, network *types.NetworkIdentifier) (*fetcher.Fetcher, error) {

	var nodeURL string = viper.GetString("node_url")

	fetcherOpts := []fetcher.Option{
		fetcher.WithTimeout(time.Duration(viper.GetInt("timeout")) * time.Second),
	}

	f := fetcher.New(
		nodeURL,
		fetcherOpts...,
	)

	_, _, fetchErr := f.InitializeAsserter(ctx, network)
	if fetchErr != nil {
		return nil, fetchErr.Err
	}

	return f, nil

}
