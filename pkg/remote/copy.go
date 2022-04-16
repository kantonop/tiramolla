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
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

// Upload file to Server
func (server Server) Upload(file, dest string) error {
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(server.client)
	if err != nil {
		return fmt.Errorf("error spawning sftp remote session: %v", err)
	}
	defer sftp.Close()

	// open the source file
	srcFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// construct the file path at destination
	// if BecomeUser is not set, prepend the destination path
	// otherwise copy to WD to pick it up later
	filename := filepath.Base(file)
	if server.BecomeUser == "" {
		filename = filepath.Join(dest, filename)
	}

	// create the destination file
	dstFile, err := sftp.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	// write to file
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	// if BecomeUser is set then the upload so far is an intermediate step
	// copy the file to its final destination by becoming the BecomeUser
	if server.BecomeUser != "" {
		// remove the file from its intermediate stop
		defer sftp.Remove(filename)

		wd, err := sftp.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %v", err)
		}

		err = server.copyAsBecomeUser(filepath.Join(wd, filename), dest, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// Download file to Server
func (server Server) Download(file, dest string) error {
	filename := filepath.Base(file)

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(server.client)
	if err != nil {
		return fmt.Errorf("error spawning sftp remote session: %v", err)
	}
	defer sftp.Close()

	// if BecomeUser is set, next steps are:
	// 1. switch user and copy to the working directory of the main user
	// 2. proceed with downloading from that path
	if server.BecomeUser != "" {
		wd, err := sftp.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %v", err)
		}

		err = server.copyAsBecomeUser(file, wd, true)
		if err != nil {
			return err
		}

		// the file is moved to main user WD, removing it as clean up
		defer sftp.Remove(filename)

		// the path is set to bare filename as the file resides in WD of main user
		file = filename
	}

	// open the source file
	srcFile, err := sftp.Open(file)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// create the destination file
	filename = filepath.Join(dest, filename)
	dstFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	// write to file
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

// copy file to destination
// switch user to BecomeUser to do the copy
func (server Server) copyAsBecomeUser(file, dest string, giveReadPerm bool) error {
	// copy file on the same server as the BecomeUser
	cmd := fmt.Sprintf("sudo su - %s -c 'cp %s %s'", server.BecomeUser, file, dest)
	sess, err := server.client.NewSession()
	if err != nil {
		return fmt.Errorf("error spawning remote session: %v", err)
	}
	defer sess.Close()

	err = sess.Run(cmd)
	if err != nil {
		return fmt.Errorf("error running '%s': %v", cmd, err)
	}

	// give read permissions to the file in destination
	// required to allow downloading as the main user
	if giveReadPerm {
		filename := filepath.Base(file)
		cmd = fmt.Sprintf("sudo su - %s -c 'chmod o+r %s/%s'", server.BecomeUser, dest, filename)
		sess, err = server.client.NewSession()
		if err != nil {
			return fmt.Errorf("error spawning remote session: %v", err)
		}
		defer sess.Close()

		err = sess.Run(cmd)
		if err != nil {
			return fmt.Errorf("error running '%s': %v", cmd, err)
		}
	}

	return nil
}
