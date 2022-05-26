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
package remote

import (
	"fmt"
	"strings"
	"testing"
)

func TestConstructClientCFG(t *testing.T) {
	testCases := []struct {
		name         string
		server       Server
		userInEnvVar string
		expErr       error
		expUser      string
	}{
		{
			name:   "UnsupportedAuthMethod",
			server: Server{AuthenticationMethod: "unknown"},
			expErr: fmt.Errorf("UnsupportedAuthMethod"),
		},
		{
			name:    "UserFromYaml",
			server:  Server{User: "foo", AuthenticationMethod: "password"},
			expErr:  nil,
			expUser: "foo",
		},
		{
			name:         "UserFromEnvVar",
			server:       Server{User: "$foo", AuthenticationMethod: "password"},
			userInEnvVar: "fooFromEnvVar",
			expErr:       nil,
			expUser:      "fooFromEnvVar",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.userInEnvVar != "" {
				t.Setenv(strings.TrimPrefix(testCase.server.User, "$"), testCase.userInEnvVar)
			}

			out, err := testCase.server.constructClientCFG()
			if testCase.expErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
				}
				return
			}

			if testCase.expErr == nil && err != nil {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}

			if testCase.expUser != out.User {
				t.Fatalf("expected username '%s', got '%s'", testCase.expUser, out.User)
			}

		})
	}
}

func TestConstructAuthMethod(t *testing.T) {
	testCases := []struct {
		name   string
		server Server
		expErr error
	}{
		{
			name:   "PasswordFromYaml",
			server: Server{Name: "foo", AuthenticationMethod: "password", Pass: "foo"},
			expErr: nil,
		},
		{
			name:   "PasswordFromEnvVar",
			server: Server{Name: "foo", AuthenticationMethod: "password", Pass: "$foo"},
			expErr: nil,
		},
		{
			name:   "PasswordHexEncoded",
			server: Server{Name: "foo", AuthenticationMethod: "password", Pass: "foo%3C3"},
			expErr: nil,
		},
		{
			name:   "PasswordBadEncoding",
			server: Server{Name: "foo", AuthenticationMethod: "password", Pass: "foo%3"},
			expErr: fmt.Errorf("PasswordBadEncoding"),
		},
		{
			name:   "UnknownAuthMethod",
			server: Server{Name: "foo", AuthenticationMethod: "unknown"},
			expErr: fmt.Errorf("UnknownAuthMethod"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testCase.server.constructAuthMethod()
			if testCase.expErr != nil && err == nil ||
				testCase.expErr == nil && err != nil {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}
		})
	}
}
