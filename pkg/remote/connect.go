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
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

const DefaultPort = 22

// Connect to server
// connects to the server, hopping through all of the servers in the list starting from the last
func (server *Server) Connect() error {
	serverChain := server.serverChain
	// server is directly reachable
	if len(serverChain) == 0 {
		// connect and return
		client, err := server.directConnect()
		if err != nil {
			return err
		}

		server.client = client
		return nil
	}

	// connect to last Server in the list and pop
	nextServer := serverChain[len(server.serverChain)-1]
	serverChain = serverChain[:(len(server.serverChain) - 1)]

	client, err := nextServer.directConnect()
	if err != nil {
		return err
	}

	// check if there are servers in the jump chain
	for len(serverChain) > 0 {
		// pop the server last in the list
		nextServer = serverChain[len(server.serverChain)-1]
		server.serverChain = serverChain[:(len(server.serverChain) - 1)]

		client, err = nextServer.hopConnect(client)
		if err != nil {
			return err
		}
	}

	// hop connect to target server
	client, err = server.hopConnect(client)
	if err != nil {
		return err
	}

	server.client = client
	return nil
}

// Closes the client
func (server Server) CloseClient() error {
	if server.client == nil {
		return fmt.Errorf("Client is not set up")
	}

	return server.client.Close()
}

// connect to server directly
func (server Server) directConnect() (*ssh.Client, error) {
	port := DefaultPort
	if server.Port != 0 {
		port = server.Port
	}
	host := fmt.Sprintf("%s:%d", server.Addr, port)
	clientCFG, err := server.constructClientCFG()
	if err != nil {
		return nil, err
	}

	return ssh.Dial("tcp", host, clientCFG)
}

// connect to server with a hop from a server we already are connected
func (server Server) hopConnect(prevClient *ssh.Client) (*ssh.Client, error) {
	port := DefaultPort
	if server.Port != 0 {
		port = server.Port
	}
	host := fmt.Sprintf("%s:%d", server.Addr, port)
	clientCFG, err := server.constructClientCFG()
	if err != nil {
		return nil, err
	}

	netConn, err := prevClient.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	conn, chans, reqs, err := ssh.NewClientConn(netConn, host, clientCFG)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(conn, chans, reqs), nil
}

// construct the client configuration for the ssh calls
func (server Server) constructClientCFG() (*ssh.ClientConfig, error) {
	authMethods, err := server.constructAuthMethod()
	if err != nil {
		return nil, err
	}

	// if username in yaml file starts with $
	// look for env variable
	user := server.User
	re := regexp.MustCompile(`^\$`)
	if re.MatchString(user) {
		user = os.Getenv(strings.TrimPrefix(server.User, "$"))
	}

	clientCFG := ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &clientCFG, nil
}

// constructs the authentication method for the ClientConfig
func (server Server) constructAuthMethod() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod
	switch server.AuthenticationMethod {
	case "password":
		// if password in yaml file starts with $
		// look for env variable
		pass := server.Pass
		re := regexp.MustCompile(`^\$`)
		if re.MatchString(pass) {
			pass = os.Getenv(strings.TrimPrefix(server.Pass, "$"))
		}

		decodedPass, err := url.QueryUnescape(pass)
		if err != nil {
			return nil, fmt.Errorf("error decoding password %s: %v", pass, err)
		}
		authMethods = append(authMethods, ssh.Password(decodedPass))
	default:
		return nil, fmt.Errorf("authentication method %s is not supported", server.AuthenticationMethod)
	}

	return authMethods, nil
}
