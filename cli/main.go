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
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ntsanov/yarc/cli/cmd"
	"github.com/spf13/viper"
)

var (
	version   string
	commit    string
	buildtime string
)

func main() {
	// We need this for static builds
	// Can't use musl because of dependencies and with gcc works only if
	// libc versions match. This forces to use go native dns queries
	os.Setenv("GODEBUG", "netdns=go")
	// rosetta-sdk-go uses logger to print the error as well as returning it
	// We use our own json error wrapper so we don't really need it. Might be
	// useful to turn it on for debugging though
	viper.Set("version", version)
	viper.Set("commit", commit)
	viper.Set("build_time", buildtime)
	log.SetOutput(ioutil.Discard)
	cmd.Execute()
}
