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

	"golang.org/x/crypto/ssh"
)

type ServerInterface interface {
	GetName() string
	CreateServerChain(servers map[string]ServerInterface) error
	Connect() error
	CloseClient() error
	Download(file, dest string) error
	Upload(file, dest string) error
}

type Server struct {
	Name                 string `mapstructure:"name"`
	Addr                 string `mapstructure:"addr"`
	Port                 int    `mapstructure:"port"`
	AuthenticationMethod string `mapstructure:"authentication_method"`
	User                 string `mapstructure:"user"`
	Pass                 string `mapstructure:"pass"`
	Gateway              string `mapstructure:"gateway"`
	BecomeUser           string `mapstructure:"become_user"`
	serverChain          []Server
	client               *ssh.Client
}

// Getter method for Name field
func (server Server) GetName() string {
	return server.Name
}

// CreateServerChain populates the serverChain field
// with a list of servers that should be used as gateways
// server index 0 is the server closest to our Server
func (server *Server) CreateServerChain(servers map[string]ServerInterface) error {
	var serverChain []Server
	curServer := server

	for curServer.Gateway != "" {
		gatewayServer, ok := servers[curServer.Gateway]
		if !ok {
			return fmt.Errorf("server %s, gateway of %s, is not known", curServer.Gateway, curServer.Name)
		}
		curServer, _ = gatewayServer.(*Server)
		serverChain = append(serverChain, *curServer)
	}

	server.serverChain = serverChain
	return nil
}

// Stringer method for Server struct
func (server Server) String() string {
	str := make([]string, 0)

	if server.Name != "" {
		str = append(str, fmt.Sprintf("Name: %s", server.Name))
	}
	if server.Addr != "" {
		str = append(str, fmt.Sprintf("Addr: %s", server.Addr))
	}
	if server.Port != 0 {
		str = append(str, fmt.Sprintf("Port: %d", server.Port))
	}
	if server.AuthenticationMethod != "" {
		str = append(str, fmt.Sprintf("AuthenticationMethod: %s", server.AuthenticationMethod))
	}
	if server.User != "" {
		str = append(str, fmt.Sprintf("User: %s", server.User))
	}
	if server.Pass != "" {
		str = append(str, fmt.Sprintf("Pass: %s", server.Pass))
	}
	if server.Gateway != "" {
		str = append(str, fmt.Sprintf("Gateway: %s", server.Gateway))
	}
	if server.BecomeUser != "" {
		str = append(str, fmt.Sprintf("BecomeUser: %s", server.BecomeUser))
	}

	return strings.Join(str, "\n")
}
