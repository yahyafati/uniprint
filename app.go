package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

const DEFAULT_HOST = "login.informatik.uni-freiburg.de"

// Dummy printer struct and list
type Printer struct {
	Name          string
	Colour        string
	PrintingSpeed string
	Model         string
}

var PRINTERS_LIST = []Printer{
	{"tfppr1", "b/w", "35 p/m", "Ricoh MP 3554"},
	{"tfppr2", "Color", "38 p/m", "HP Color LaserJet Enterprise M553"},
	{"tfppr3", "b/w", "45 p/m", "HP LaserJet Enterprise MFP M528"},
	{"tfppr4", "Color", "56 p/m", "HP Color LaserJet Enterprise M653"},
}

func loadUserCredentials() (string, string, bool) {
	_ = godotenv.Load()

	email := os.Getenv("UNIPRINT_EMAIL")
	password := os.Getenv("UNIPRINT_PASSWORD")

	if email == "" || password == "" {
		return "", "", false
	}

	return email, password, true
}

func getUserCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	// Email input
	fmt.Print("Enter your email (e.g., user@host): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if !strings.Contains(email, "@") {
		email = email + "@" + DEFAULT_HOST
	}

	fmt.Print("Enter your password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	password := string(passwordBytes)

	return email, password
}

func getFilesToPrint() []string {
	reader := bufio.NewReader(os.Stdin)
	var files []string

	for {
		fmt.Print("Enter file path to print (leave empty to finish): ")
		input, _ := reader.ReadString('\n')
		file := strings.TrimSpace(input)

		if file == "" {
			break
		}

		if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
			files = append(files, file)
		} else {
			fmt.Println("Invalid file. Try again.")
		}
	}

	fmt.Printf("Found %d files to print.\n", len(files))
	return files
}

func confirmInput(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt + ": ")
		input, _ := reader.ReadString('\n')
		answer := strings.TrimSpace(strings.ToLower(input))

		if answer == "y" || answer == "yes" {
			return true
		} else if answer == "n" || answer == "no" || answer == "" {
			return false
		} else {
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}

func askDeleteAfterPrint() bool {
	return confirmInput("Delete files after printing? (y/N)")
}

func sftpCopyFiles(email, password string, files []string) {
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

func sshPrintFiles(email, password string, files []string, deleteAfter bool) {
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

func main() {
	email, password, loadedCredsFromEnv := loadUserCredentials()
	if loadedCredsFromEnv {
		fmt.Println("Loaded credentials of", email, "from environment variables")
		loadedCredsFromEnv = !confirmInput("Do you wish to enter manually instead? [y/N]")
	}
	if !loadedCredsFromEnv {
		email, password = getUserCredentials()
	}
	fmt.Println("Hi,", email)
	files := getFilesToPrint()
	if len(files) == 0 {
		fmt.Println("No files selected!\nBye!")
		return
	}
	deleteLater := askDeleteAfterPrint()
	if deleteLater {
		fmt.Println("Files will be removed from the server once it has been printed.")
	}

	sftpCopyFiles(email, password, files)
	sshPrintFiles(email, password, files, deleteLater)
}
