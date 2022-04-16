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
	"reflect"
	"testing"
)

func TestGetName(t *testing.T) {
	testCases := []struct {
		name   string
		server Server
		expOut string
	}{
		{name: "Foo", server: Server{Name: "foo"}, expOut: "foo"},
		{name: "Empty", server: Server{}, expOut: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := testCase.server
			out := server.GetName()

			if testCase.expOut != out {
				t.Fatalf("expected %s, got %s", testCase.expOut, out)
			}
		})
	}
}

func TestCreateServerChain(t *testing.T) {
	testCases := []struct {
		name      string
		server    Server
		servers   map[string]Server
		expServer Server
		expErr    error
	}{
		{
			name:   "EmptyServerChain",
			server: Server{Name: "foo"},
			servers: map[string]Server{
				"foo": {Name: "foo"},
				"bar": {Name: "bar"},
			},
			expServer: Server{Name: "foo"},
			expErr:    nil,
		},
		{
			name:   "OneHopChain",
			server: Server{Name: "foo", Gateway: "bar"},
			servers: map[string]Server{
				"foo": {Name: "foo", Gateway: "bar"},
				"bar": {Name: "bar"},
			},
			expServer: Server{Name: "foo", Gateway: "bar", serverChain: []Server{{Name: "bar"}}},
			expErr:    nil,
		},
		{
			name:   "TwoHopsChain",
			server: Server{Name: "foo", Gateway: "bar"},
			servers: map[string]Server{
				"foo": {Name: "foo", Gateway: "bar"},
				"bar": {Name: "bar", Gateway: "qux"},
				"qux": {Name: "qux"},
			},
			expServer: Server{Name: "foo", Gateway: "bar", serverChain: []Server{{Name: "bar", Gateway: "qux"}, {Name: "qux"}}},
			expErr:    nil,
		},
		{
			name:   "BrokenChain",
			server: Server{Name: "foo", Gateway: "barz"},
			servers: map[string]Server{
				"foo": {Name: "foo", Gateway: "barz"},
				"bar": {Name: "bar"},
			},
			expErr: fmt.Errorf("unknown gateway server"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := testCase.server
			servers := make(map[string]ServerInterface)
			for key, val := range testCase.servers {
				s := val
				servers[key] = &s
			}
			err := server.CreateServerChain(servers)
			if testCase.expErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
				}
				return
			}

			if testCase.expErr == nil && err != nil {
				t.Fatalf("expected error '%v', got '%v'", testCase.expErr, err)
			}

			if !reflect.DeepEqual(server, testCase.expServer) {
				t.Fatalf("expected '%v', got '%v'", testCase.expServer, server)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name   string
		server Server
		expOut string
	}{
		{
			name:   "OnlyName",
			server: Server{Name: "foo"},
			expOut: "Name: foo",
		},
		{
			name:   "NameAddrGateway",
			server: Server{Name: "foo", Addr: "1.1.1.1", Gateway: "bar"},
			expOut: "Name: foo\nAddr: 1.1.1.1\nGateway: bar",
		},
		{
			name:   "AllFields",
			server: Server{Name: "foo", Addr: "1.1.1.1", Port: 22, AuthenticationMethod: "password", User: "foo", Pass: "bar", Gateway: "qux", BecomeUser: "foobar"},
			expOut: "Name: foo\nAddr: 1.1.1.1\nPort: 22\nAuthenticationMethod: password\nUser: foo\nPass: bar\nGateway: qux\nBecomeUser: foobar",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := testCase.server.String()
			if testCase.expOut != out {
				t.Fatalf("expected '%s', got '%s'", testCase.expOut, out)
			}
		})
	}
}
