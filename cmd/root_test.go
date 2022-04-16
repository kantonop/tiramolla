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
	"github/kantonop/tiramolla/pkg/remote"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestInitConfig(t *testing.T) {
	testCases := []struct {
		name       string
		config     string
		expServers map[string]remote.Server
	}{
		{
			name:   "OneServer",
			config: "servers:\n  - name: foo\n    addr: 1.1.1.1\n    gateway: bar",
			expServers: map[string]remote.Server{
				"foo": {Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
			},
		},
		{
			name:   "TwoServers",
			config: "servers:\n  - name: foo\n    addr: 1.1.1.1\n    gateway: bar\n  - name: bar\n    addr: 2.2.2.2",
			expServers: map[string]remote.Server{
				"foo": {Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
				"bar": {Name: "bar", Addr: "2.2.2.2"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			home, err := os.UserHomeDir()
			if err != nil {
				t.Fatalf("error getting home directory: %v", err)
			}
			filename := filepath.Join(home, ".tiramolla.yaml")
			err, cleanup := testFileCreator(testCase.config, filename)
			if err != nil {
				t.Fatalf("error creating test config file: %v", err)
			}
			t.Cleanup(cleanup)

			expServers := make(map[string]remote.ServerInterface)
			for key, val := range testCase.expServers {
				server := val
				expServers[key] = &server
			}
			initConfig()
			if !reflect.DeepEqual(expServers, servers) {
				t.Fatalf("expected '%v', got '%v'", expServers, servers)
			}
		})
	}
}

func TestServerListToMap(t *testing.T) {
	testCases := []struct {
		name    string
		servers []remote.Server
		expOut  map[string]remote.Server
	}{
		{
			name:    "OneServer",
			servers: []remote.Server{{Name: "foo", Addr: "1.1.1.1", Gateway: "bar"}},
			expOut: map[string]remote.Server{
				"foo": {Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
			},
		},
		{
			name: "TwoServers",
			servers: []remote.Server{
				{Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
				{Name: "bar", Addr: "2.2.2.2"},
			},
			expOut: map[string]remote.Server{
				"foo": {Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
				"bar": {Name: "bar", Addr: "2.2.2.2"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			expOut := make(map[string]remote.ServerInterface)
			for key, val := range testCase.expOut {
				server := val
				expOut[key] = &server
			}
			out := serverListToMap(testCase.servers)
			if !reflect.DeepEqual(expOut, out) {
				t.Fatalf("expected '%v', got '%v'", expOut, out)
			}
		})
	}
}

// helper function to create file with test content
// returns a cleanup function to use with t.Cleanup
func testFileCreator(fileContent, filename string) (error, func()) {
	originalFile, err := os.Open(filename)
	if err != nil {
		return err, nil
	}
	defer originalFile.Close()

	originalConfig, err := io.ReadAll(originalFile)

	file, err := os.Create(filename)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	_, err = io.WriteString(file, fileContent)
	if err != nil {
		return err, nil
	}

	cleanup := func() {
		file, _ := os.Create(filename)
		defer file.Close()
		io.WriteString(file, string(originalConfig))
	}
	return nil, cleanup
}
