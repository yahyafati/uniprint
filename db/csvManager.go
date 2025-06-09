package db

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type Credential struct {
	Username string
	Password string
}

type CSVManager struct {
	filePath string
	pin      string
	rows     []Credential
}

func newCSVManager(filePath, pin string) (*CSVManager, error) {
	manager := &CSVManager{
		filePath: filePath,
		pin:      pin,
	}

	if err := manager.load(); err != nil {
		return nil, err
	}

	return manager, nil
}

func (c *CSVManager) AddRow(cred Credential) error {
	c.rows = append([]Credential{cred}, c.rows...)
	return c.save()
}

func (c *CSVManager) UpdateRow(username string, newCred Credential) error {
	found := false
	for i, cred := range c.rows {
		if cred.Username == username {
			c.rows = append(c.rows[:i], c.rows[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("no credential found for username: %s", username)
	}
	c.rows = append([]Credential{newCred}, c.rows...)
	return c.save()
}

func (c *CSVManager) AccessRow(username string) (*Credential, error) {
	for i, cred := range c.rows {
		if cred.Username == username {
			// Move to top
			c.rows = append([]Credential{cred}, append(c.rows[:i], c.rows[i+1:]...)...)
			_ = c.save()
			return &cred, nil
		}
	}
	return nil, fmt.Errorf("username not found: %s", username)
}

func (c *CSVManager) load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	plaintext, err := decrypt(data, c.pin)
	if err != nil {
		return errors.New("failed to decrypt file — invalid PIN or corrupted file")
	}

	reader := csv.NewReader(bytes.NewReader(plaintext))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("CSV file is empty")
	}

	// Skip header
	for _, record := range records[1:] {
		if len(record) >= 2 {
			c.rows = append(c.rows, Credential{
				Username: record[0],
				Password: record[1],
			})
		}
	}
	return nil
}

func (c *CSVManager) save() error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write([]string{"username", "password"})

	for _, cred := range c.rows {
		writer.Write([]string{cred.Username, cred.Password})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	encrypted, err := encrypt(buf.Bytes(), c.pin)
	if err != nil {
		return err
	}

	return os.WriteFile(c.filePath, encrypted, 0600)
}
