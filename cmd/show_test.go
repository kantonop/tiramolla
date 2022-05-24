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
	"io"
	"os"
	"testing"

	"github.com/kantonop/tiramolla/pkg/remote"

	"github.com/spf13/cobra"
)

func TestShow(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		servers    map[string]remote.Server
		matchExpr  string
		expErr     error
		expPrinted string
	}{
		{
			name:       "ShowServers",
			args:       []string{"servers"},
			servers:    map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			expErr:     nil,
			expPrinted: "bar\nfoo\n",
		},
		{
			name:       "ShowServersValidMatch",
			args:       []string{"servers"},
			servers:    map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			matchExpr:  "ar",
			expErr:     nil,
			expPrinted: "bar\n",
		},
		{
			name:      "ShowServersNoMatchingServer",
			args:      []string{"servers"},
			servers:   map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			matchExpr: "foobar",
			expErr:    fmt.Errorf("NoMatch"),
		},
		{
			name:      "ShowServersInvalidMatch",
			args:      []string{"servers"},
			servers:   map[string]remote.Server{"foo": {Name: "foo"}, "bar": {Name: "bar"}},
			matchExpr: "\\d(+",
			expErr:    fmt.Errorf("cannot parse match expr"),
		},
		{
			name:    "UnknownServer",
			args:    []string{"unknownServer"},
			servers: map[string]remote.Server{"foo": {Name: "foo"}},
			expErr:  fmt.Errorf("unknownServer"),
		},
		{
			name:       "ShowServerDetails",
			args:       []string{"foo"},
			servers:    map[string]remote.Server{"foo": {Name: "foo", Addr: "1.1.1.1", Gateway: "bar"}, "bar": {Name: "bar"}},
			expErr:     nil,
			expPrinted: "Name: foo\nAddr: 1.1.1.1\nGateway: bar\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			origStdout := os.Stdout

			r, w, setupErr := os.Pipe()
			if setupErr != nil {
				t.Fatalf("FAILED: %s\ncouldn't create pipe", testCase.name)
			}
			os.Stdout = w

			cmd := cobra.Command{}
			for key, val := range testCase.servers {
				s := val
				servers[key] = &s
			}
			matchExpr = testCase.matchExpr
			err := show(&cmd, testCase.args)
			w.Close()

			printed, setupErr := io.ReadAll(r)
			if setupErr != nil {
				t.Fatalf("FAILED: %s\ncouldn't read from pipe", testCase.name)
			}
			r.Close()
			fmt.Println(string(printed))
			os.Stdout = origStdout // restore original Stdout

			if testCase.expErr != nil {
				if err == nil {
					t.Fatalf("FAILED: %s\nexpected error '%v', got '%v'", testCase.name, testCase.expErr, err)
				}
				return
			}
			if testCase.expErr == nil && err != nil {
				t.Fatalf("FAILED: %s\nexpected error '%v', got '%v'", testCase.name, testCase.expErr, err)
			}

			if string(printed) != testCase.expPrinted {
				t.Fatalf("FAILED: %s\nexpected printed message '%s', got '%s'", testCase.name, testCase.expPrinted, printed)
			}
		})
	}
}
