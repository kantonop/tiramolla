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

	"github.com/spf13/cobra"
)

var (
	targetServer, mode string
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy /path/to/file /path/to/dest",
	Short: "download or upload a file",
	Long: `Downloads or uploads a file to or from a remote server.

Use absolute paths to avoid unexpected behaviour.
Downloading and uploading is set by the mode flag.`,
	Args:    cobra.ExactArgs(2),
	PreRunE: copyFlagsValidation,
	RunE:    copyFile,
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().StringVar(&targetServer, "server", "", "target server")
	copyCmd.MarkFlagRequired("server")
	copyCmd.Flags().StringVar(&mode, "mode", "", "down or up")
	copyCmd.MarkFlagRequired("mode")
}

// flags validation function
// runs before main copy command
func copyFlagsValidation(cmd *cobra.Command, args []string) error {
	err := validateTargetServer()
	if err != nil {
		return err
	}

	if mode != "down" && mode != "up" {
		return fmt.Errorf("mode should be either 'down' or 'up'")
	}

	return nil
}

// tiramolla copy command
func copyFile(cmd *cobra.Command, args []string) error {
	server := servers[targetServer]
	// chain servers
	err := server.CreateServerChain(servers)
	if err != nil {
		return fmt.Errorf("creation of chain of servers to target server failed with error: %v", err)
	}

	// Connect
	err = server.Connect()
	if err != nil {
		return fmt.Errorf("connect to server failed with error: %v", err)
	}
	defer server.CloseClient()

	// copy
	switch mode {
	case "down":
		err = server.Download(args[0], args[1])
		if err != nil {
			return fmt.Errorf("download failed with error: %v", err)
		}
		fmt.Println("download completed")
	case "up":
		err = server.Upload(args[0], args[1])
		if err != nil {
			return fmt.Errorf("upload failed with error: %v", err)
		}
		fmt.Println("upload completed")
	}

	return nil
}

// Checks if target server exists in configuration
// returns error if not, nil otherwise
func validateTargetServer() error {
	_, ok := servers[targetServer]
	if !ok {
		return fmt.Errorf("server %s not in list of known servers, use 'tiramolla show servers' for the list of available servers", targetServer)
	}
	return nil
}
