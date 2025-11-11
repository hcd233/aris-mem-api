// Package client provides client utilities
package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	// TokenFileName is the local token file name
	TokenFileName = ".aris/aris-mem-token"
	// EncryptionKey is the encryption key (should use more secure key management in production)
	EncryptionKey = "aris-mem-api-client-secret-key-817125"
)

// TokenData represents locally stored token data
type TokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// getTokenFilePath returns the path to the token file
func getTokenFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, TokenFileName), nil
}

// LoadToken loads and decrypts token from local storage
func LoadToken() (*TokenData, error) {
	filePath, err := getTokenFilePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New("token file not found")
	}

	// Read encrypted data
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	decryptedData, err := decrypt(encryptedData, []byte(deriveKey(EncryptionKey)))
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var tokenData TokenData
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, err
	}

	return &tokenData, nil
}

// SaveToken encrypts and saves token to local storage
func SaveToken(tokenData *TokenData) error {
	filePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return err
	}

	// Encrypt data
	encryptedData, err := encrypt(jsonData, []byte(deriveKey(EncryptionKey)))
	if err != nil {
		return err
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(filePath, encryptedData, 0o600); err != nil {
		return err
	}

	return nil
}

// DeleteToken deletes the local token file
func DeleteToken() error {
	filePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// deriveKey derives a 32-byte key from the input string
func deriveKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return string(hash[:])
}

// encrypt encrypts data using AES-GCM
func encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Base64 encode for readability
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(encoded, ciphertext)

	return encoded, nil
}

// decrypt decrypts data using AES-GCM
func decrypt(ciphertext, key []byte) ([]byte, error) {
	// Decode base64
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))
	n, err := base64.StdEncoding.Decode(decoded, ciphertext)
	if err != nil {
		return nil, err
	}
	decoded = decoded[:n]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherData := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
