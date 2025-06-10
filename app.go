package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yahyafati/uniprint/db"
	"github.com/yahyafati/uniprint/prompt"
	"github.com/yahyafati/uniprint/tools"
)

const DEFAULT_HOST = "login.informatik.uni-freiburg.de"

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

func main() {
	email, password, loadedCredsFromEnv := loadUserCredentials()
	if loadedCredsFromEnv {
		fmt.Println("Loaded credentials of", email, "from environment variables")
		loadedCredsFromEnv = !prompt.ConfirmInput("Do you wish to enter manually instead?")
	}
	if !loadedCredsFromEnv {
		_, err := db.InitializeCSV()
		if err != nil {
			log.Fatal("Error while initializing CSV", err.Error())
		}
		email, password = prompt.GetUserCredentials()
	}

	fmt.Println("Hi,", email)
	files := prompt.GetFilesToPrint()
	if len(files) == 0 {
		fmt.Println("No files selected!\nBye!")
		return
	}
	deleteLater := prompt.AskDeleteAfterPrint()
	if deleteLater {
		fmt.Println("Files will be removed from the server once it has been printed.")
	}
	oneSided := prompt.AskOneSidedPrint()

	tools.SftpCopyFiles(email, password, files)
	if prompt.ConfirmInputWithDefault("Copied! Print?", true) {
		tools.SshPrintFiles(email, password, files, deleteLater, oneSided)
	} else {
		fmt.Println("Not Printing")
	}
}
