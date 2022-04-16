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
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	matchExpr string
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show {servers|serverName}",
	Short: "print configured server names or details for a specific server",
	Long: `Shows configuration details of a specific server, if serverName is provided.
Shows list of servers if 'servers' argument is passed, list can be matched against a regular expression if the flag is provided.`,
	Args: cobra.ExactArgs(1),
	RunE: show,
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().StringVar(&matchExpr, "match", "", "expression to match server names")
}

// tiramolla show command
func show(cmd *cobra.Command, args []string) error {
	if args[0] == "servers" {
		re, err := regexp.Compile(matchExpr)
		if err != nil {
			return fmt.Errorf("error in parsing match expression: %v", err)
		}
		serversList := make([]string, 0, len(servers))
		for server := range servers {
			if !re.Match([]byte(server)) {
				continue
			}
			serversList = append(serversList, server)
		}
		sort.Strings(serversList)

		if len(serversList) == 0 {
			return fmt.Errorf("no servers are configured and matching")
		}
		fmt.Println(strings.Join(serversList, "\n"))

		return nil
	}

	server, ok := servers[args[0]]
	if !ok {
		return fmt.Errorf("server %s is not in the list of known servers, use 'tiramolla show servers' for the list of available servers", args[0])
	}

	fmt.Println(server)
	return nil
}
