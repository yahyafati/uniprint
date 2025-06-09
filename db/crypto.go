package db

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

func deriveKey(pin string) []byte {
	hash := sha256.Sum256([]byte(pin))
	return hash[:]
}

func encrypt(data []byte, pin string) ([]byte, error) {
	key := deriveKey(pin)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padding := aes.BlockSize - len(data)%aes.BlockSize
	padded := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	copy(iv, key[:aes.BlockSize]) // for simplicity; not secure for real use

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	return ciphertext, nil
}

func decrypt(ciphertext []byte, pin string) ([]byte, error) {
	key := deriveKey(pin)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	padding := int(ciphertext[len(ciphertext)-1])
	if padding > aes.BlockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	return ciphertext[:len(ciphertext)-padding], nil
}

func promptAndValidatePIN(reader *bufio.Reader, confirm bool) (string, error) {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	for {
		fmt.Print("Enter PIN: ")
		pin, _ := reader.ReadString('\n')
		pin = strings.TrimSpace(pin)

		if !re.MatchString(pin) {
			fmt.Println("PIN must be alphanumeric.")
			continue
		}

		if confirm {
			fmt.Print("Confirm PIN: ")
			confirmPin, _ := reader.ReadString('\n')
			confirmPin = strings.TrimSpace(confirmPin)
			if pin != confirmPin {
				fmt.Println("PINs do not match.")
				continue
			}
		}

		return pin, nil
	}
}
