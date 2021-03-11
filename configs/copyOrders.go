package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

//CopyFileToAS2 copy a file to AS2 server in PATH_IN_AS2
func CopyFileToAS2(pathOfFile string) {

	path, err1 := os.Getwd()
	if err1 != nil {
		fmt.Println(err1)

		log.Println(err1)
	}
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey(GetUserServerAS2(), path+GetRSAPath(), ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(GetServerAS2(), &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)

		log.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Open a file
	f, err2 := os.Open(pathOfFile)
	if err2 != nil {
		fmt.Println("eROOR OPEN FILE:::", err2)

		log.Println(err2)
	}
	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)

	_, na := filepath.Split(pathOfFile)

	err = client.CopyFile(f, GetPathInAS2()+na, "0655")

	if err != nil {
		fmt.Println("Error while copying file::: ", err)
	}
}
