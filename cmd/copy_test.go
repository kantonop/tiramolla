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
	"testing"

	"github.com/kantonop/tiramolla/pkg/remote"

	"github.com/spf13/cobra"
)

var (
	cmd  *cobra.Command
	args []string
)

func TestValidateTargetServer(t *testing.T) {
	testCases := []struct {
		name         string
		targetServer string
		servers      map[string]remote.Server
		expErr       error
	}{
		{
			name:         "KnownServer",
			targetServer: "foo",
			servers:      map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			expErr:       nil,
		},
		{
			name:         "UnknownServer",
			targetServer: "unknownServer",
			servers:      map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			expErr:       fmt.Errorf("unknown server error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			targetServer = testCase.targetServer
			servers = make(map[string]remote.ServerInterface)
			for key, val := range testCase.servers {
				server := val
				servers[key] = &server
			}

			err := validateTargetServer()
			if (testCase.expErr == nil && err != nil) ||
				(testCase.expErr != nil && err == nil) {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}
		})
	}
}

func TestCopyFlagsValidation(t *testing.T) {
	testCases := []struct {
		name    string
		mode    string
		servers map[string]remote.Server
		expErr  error
	}{
		{name: "Download", mode: "down", expErr: nil},
		{name: "Upload", mode: "up", expErr: nil},
		{name: "IncorrectMode", mode: "downAndUp", expErr: fmt.Errorf("download/upload mode error")},
		{name: "UnknownServer", servers: make(map[string]remote.Server), expErr: fmt.Errorf("unknown server")},
	}

	for _, testCase := range testCases {
		targetServer = "foo"
		servers = make(map[string]remote.ServerInterface)
		servers["foo"] = &remote.Server{Name: "foo"}

		t.Run(testCase.name, func(t *testing.T) {
			if testCase.servers != nil {
				delete(servers, "foo")
				for key, val := range testCase.servers {
					server := val
					servers[key] = &server
				}
			}
			mode = testCase.mode
			err := copyFlagsValidation(cmd, args)

			if (testCase.expErr == nil && err != nil) ||
				(testCase.expErr != nil && err == nil) {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}
		})
	}
}

type ServerMock struct {
	name                 string
	createServerChainErr error
	connectErr           error
	closeClientErr       error
	downloadErr          error
	uploadErr            error
}

func (serverMock ServerMock) GetName() string {
	return serverMock.name
}

func (serverMock ServerMock) CreateServerChain(servers map[string]remote.ServerInterface) error {
	return serverMock.createServerChainErr
}

func (serverMock ServerMock) Connect() error {
	return serverMock.connectErr
}

func (serverMock ServerMock) CloseClient() error {
	return serverMock.closeClientErr
}

func (serverMock ServerMock) Download(file, dest string) error {
	return serverMock.downloadErr
}

func (serverMock ServerMock) Upload(file, dest string) error {
	return serverMock.uploadErr
}

func TestCopyFile(t *testing.T) {
	testCases := []struct {
		name    string
		mode    string
		servers map[string]ServerMock
		expErr  error
	}{
		{
			name:    "ChainServersError",
			servers: map[string]ServerMock{"foo": {createServerChainErr: fmt.Errorf("ChainServersError")}},
			expErr:  fmt.Errorf("ChainServersError"),
		},
		{
			name:    "ConnectError",
			servers: map[string]ServerMock{"foo": {connectErr: fmt.Errorf("ConnectError")}},
			expErr:  fmt.Errorf("ConnectError"),
		},
		{
			name:    "DownloadError",
			mode:    "down",
			servers: map[string]ServerMock{"foo": {downloadErr: fmt.Errorf("DownloadError")}},
			expErr:  fmt.Errorf("DownloadError"),
		},
		{
			name:    "UploadError",
			mode:    "up",
			servers: map[string]ServerMock{"foo": {uploadErr: fmt.Errorf("UploadError")}},
			expErr:  fmt.Errorf("UploadError"),
		},
		{
			name:    "DownloadSuccess",
			mode:    "down",
			servers: map[string]ServerMock{"foo": {}},
			expErr:  nil,
		},
		{
			name:    "UploadSuccess",
			mode:    "up",
			servers: map[string]ServerMock{"foo": {}},
			expErr:  nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			targetServer = "foo"
			args = []string{"fileSource", "fileDestination"}
			mode = testCase.mode

			servers = make(map[string]remote.ServerInterface)
			for key, val := range testCase.servers {
				servers[key] = val
			}
			err := copyFile(cmd, args)

			if (testCase.expErr == nil && err != nil) ||
				(testCase.expErr != nil && err == nil) {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}
		})
	}
}
