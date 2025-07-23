package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Encrypt encrypts a plaintext string using AES-256-GCM
func Encrypt(plaintext, key string) (string, error) {
	// Convert key to bytes (should be 32 bytes for AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// Pad key if it's too short
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		// Truncate key if it's too long
		keyBytes = keyBytes[:32]
	}

	// Create cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// Decrypt decrypts an encrypted string using AES-256-GCM
func Decrypt(encryptedText, key string) (string, error) {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Convert key to bytes (should be 32 bytes for AES-256)
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// Pad key if it's too short
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		// Truncate key if it's too long
		keyBytes = keyBytes[:32]
	}

	// Create cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// SearchableEncrypt encrypts text in a way that allows for searchable encryption
// This is a simplified implementation - in production, you might want to use
// more sophisticated searchable encryption schemes
func SearchableEncrypt(plaintext, key string) (string, error) {
	// For searchable encryption, we'll use a deterministic encryption
	// This is not as secure as random encryption but allows for searching
	// In production, consider using techniques like:
	// - Encrypted Bloom filters
	// - Searchable symmetric encryption (SSE)
	// - Homomorphic encryption for specific use cases

	// For now, we'll use a simple approach with a fixed nonce
	// WARNING: This is not recommended for production use
	
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Use a deterministic nonce (NOT recommended for production)
	nonce := make([]byte, gcm.NonceSize())
	// In production, use a proper deterministic nonce generation

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// GenerateKey generates a random encryption key
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
} 