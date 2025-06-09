package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SftpCopyFiles(email, password string, files []string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid email format")
		os.Exit(1)
	}
	username := parts[0]
	host := parts[1]

	// SSH client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect
	addr := host + ":22"
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to connect via SSH: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start SFTP session
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Failed to start SFTP: %v\n", err)
		os.Exit(1)
	}
	defer sftpClient.Close()

	// Upload files
	for _, file := range files {
		srcFile, err := os.Open(file)
		if err != nil {
			fmt.Printf("Failed to open file %s: %v\n", file, err)
			os.Exit(1)
		}
		defer srcFile.Close()

		filename := filepath.Base(file)
		fmt.Printf("Copying %s to %s:~\n", file, email)

		dstFile, err := sftpClient.Create(filename)
		if err != nil {
			fmt.Printf("Failed to create remote file %s: %v\n", filename, err)
			os.Exit(1)
		}

		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close()
		if err != nil {
			fmt.Printf("Failed to upload file %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
}

func SshPrintFiles(email, password string, files []string, deleteAfter bool) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid email format")
		os.Exit(1)
	}
	username := parts[0]
	host := parts[1]

	// SSH config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect
	addr := host + ":22"
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("SSH connection failed: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// whoami verification
	fmt.Println("\nVerifying remote login with 'whoami'...")
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create SSH session:", err)
		os.Exit(1)
	}
	output, err := session.Output("whoami")
	session.Close()
	if err != nil {
		fmt.Println("Failed to run whoami:", err)
		os.Exit(1)
	}
	fmt.Printf("Connected as: %s\n\n", strings.TrimSpace(string(output)))

	// Choose printer
	fmt.Println("Choose a printer:")
	for i, p := range PRINTERS_LIST {
		fmt.Printf("%d: %s (%s)\n", i+1, p.Name, p.Colour)
	}
	fmt.Print("Enter choice number: ")

	var choice int
	fmt.Scanln(&choice)
	if choice < 1 || choice > len(PRINTERS_LIST) {
		fmt.Println("Invalid printer choice.")
		os.Exit(1)
	}
	selectedPrinter := PRINTERS_LIST[choice-1].Name
	fmt.Printf("Selected printer: %s\n", selectedPrinter)

	// Print files
	for _, file := range files {
		filename := "'" + filepath.Base(file) + "'"
		lprCmd := fmt.Sprintf("lpr -P %s %s%s", selectedPrinter,
			func() string {
				if deleteAfter {
					return "-r "
				}
				return ""
			}(), filename)

		fmt.Printf("\n(Now Printing) Command executed: %s\n", lprCmd)

		session, err := client.NewSession()
		if err != nil {
			fmt.Println("Failed to create SSH session:", err)
			continue
		}
		stdout, err := session.CombinedOutput(lprCmd)
		session.Close()
		if err != nil {
			fmt.Printf("Error printing %s: %v\n", filename, err)
			fmt.Println("Output:", string(stdout))
			continue
		}
		fmt.Printf("Printed %s to %s.\n", filename, selectedPrinter)
	}
}
