package data

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

const allowedChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

func sanitizeBase(input string) string {
	var sanitized strings.Builder
	for _, r := range input {
		if strings.ContainsRune(allowedChars, r) {
			sanitized.WriteRune(r)
		} else {
			sanitized.WriteRune('_')
		}
	}

	parts := strings.Split(sanitized.String(), "_")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	result := strings.Join(filtered, "_")
	return strings.Trim(result, "_")
}

func generateRandomSuffix(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		if err != nil {
			log.Printf("can't generate random suffix: %s", err)
			return "", err
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}

func generateRandomUsername() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(59))
	if err != nil {
		log.Printf("can't generate random username: %s", err)
		return "", err
	}
	length := int(n.Int64()) + 6
	return generateRandomSuffix(length)
}

func generateUsername(base string, suffixMinLength int) (string, error) {
	sanitized := sanitizeBase(base)
	if sanitized == "" {
		return generateRandomUsername()
	}

	if suffixMinLength == 0 {
		return sanitized, nil
	}

	const maxTotalLength = 64
	maxBaseLength := maxTotalLength - suffixMinLength
	if len(sanitized) > maxBaseLength {
		sanitized = sanitized[:maxBaseLength]
	}

	currentBaseLength := len(sanitized)
	suffixLength := max(suffixMinLength, 6-currentBaseLength)
	maxPossibleSuffix := maxTotalLength - currentBaseLength
	suffixLength = min(suffixLength, maxPossibleSuffix)

	suffix, err := generateRandomSuffix(suffixLength)
	if err != nil {
		return "", err
	}

	return sanitized + suffix, nil
}
