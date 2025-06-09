package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yahyafati/uniprint/db"
	"github.com/yahyafati/uniprint/system"
	"golang.org/x/term"
)

func GetUserCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	// Email input
	fmt.Print("Enter your email (e.g., user@host): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if !strings.Contains(email, "@") {
		email = email + "@" + system.DEFAULT_HOST
	}

	csvManager := db.GetCSVManager()
	response, err := csvManager.AccessRow(email)

	if err == nil {
		return email, response.Password
	}

	fmt.Print("Enter your password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	password := string(passwordBytes)

	saveToCSV := ConfirmInputWithDefault("Do You want to save the account to storage?", true)

	if saveToCSV {
		response = &db.Credential{
			Username: email,
			Password: password,
		}
		csvManager.AddRow(*response)
	}

	return email, password
}

func GetFilesToPrint() []string {
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

func AskDeleteAfterPrint() bool {
	return ConfirmInput("Delete files after printing? (y/N)")
}
