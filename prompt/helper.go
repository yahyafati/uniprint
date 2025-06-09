package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfirmInput prompts the user with a yes/no question. Defaults to "n" if empty.
func ConfirmInput(prompt string) bool {
	return ConfirmInputWithDefault(prompt, false)
}

// ConfirmInputWithDefault allows specifying a default answer (true = yes, false = no)
func ConfirmInputWithDefault(prompt string, onEmptyDefault bool) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		defaultString := "[y/N]"
		if onEmptyDefault {
			defaultString = "[Y/n]"
		}
		fmt.Fprintf(os.Stderr, "%s %s: ", prompt, defaultString) // Print to stderr as it's more appropriate for prompts
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		answer := strings.TrimSpace(strings.ToLower(input))

		if answer == "y" || answer == "yes" {
			return true
		} else if answer == "n" || answer == "no" {
			return false
		} else if answer == "" {
			// Use default if input is empty
			return onEmptyDefault
		} else {
			fmt.Fprintln(os.Stderr, "Please enter 'y' or 'n'.")
		}
	}
}
