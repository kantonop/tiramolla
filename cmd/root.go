/*
Copyright © 2022 Kostas Antonopoulos kost.antonopoulos@gmail.com

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
	"os"

	"github.com/kantonop/tiramolla/pkg/remote"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Servers []remote.Server
}

var servers map[string]remote.ServerInterface

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiramolla",
	Short: "file transferring from/to remote servers",
	Long:  "Utility for file transferring from/to a remote server, possibly at multiple hops distance",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// PersistentPreRun runs after args validation
		// if an error occurs after args validation, don't show usage
		cmd.SilenceUsage = true
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".tiramolla" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".tiramolla")

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	cobra.CheckErr(err)

	// Unmarshal the yaml to create the config structs
	conf := &config{}
	err = viper.Unmarshal(conf)
	cobra.CheckErr(err)
	servers = serverListToMap(conf.Servers)
}

// serverListToMap constructs a map of server interfaces with their name as key
func serverListToMap(servers []remote.Server) map[string]remote.ServerInterface {
	serversMap := make(map[string]remote.ServerInterface)
	for _, s := range servers {
		server := s
		serversMap[s.GetName()] = &server
	}

	return serversMap
}
