package k8s

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	MaxKubernetesNameLength       = 253
	MaxKubernetesLabelValueLength = 63
	DefaultSuffixLength           = 5
	DefaultSeparator              = "-"
	Separators                    = "-_."
	charset                       = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// FormatKubernetesName formats a name to be compliant with Kubernetes naming rules.
// It ensures the name is not longer than maxLength and adds a unique suffix.
// If maxLength is 0 or greater than MaxKubernetesNameLength, it uses MaxKubernetesNameLength.
// If suffixLength is 0, it uses DefaultSuffixLength.
func FormatKubernetesName(name string, maxLength, suffixLength int) string {
	if maxLength == 0 || maxLength > MaxKubernetesNameLength {
		maxLength = MaxKubernetesNameLength
	}
	if suffixLength == 0 {
		suffixLength = DefaultSuffixLength
	}

	// Generate a random suffix using the full alphabet and numbers
	suffix, _ := GenerateRandomString(suffixLength)

	// Calculate the maximum length for the original name
	maxNameLength := maxLength - suffixLength - len(DefaultSeparator)

	// Truncate the original name if necessary
	if len(name) > maxNameLength {
		name = name[:maxNameLength]
	}

	name = strings.TrimRight(name, Separators)

	return name + DefaultSeparator + suffix
}

// FormatKubernetesNameCustomSuffix formats a name to be compliant with Kubernetes naming rules.
func FormatKubernetesNameCustomSuffix(name string, maxLength int, suffix string) string {
	if maxLength == 0 || maxLength > MaxKubernetesNameLength {
		maxLength = MaxKubernetesNameLength
	}

	// Calculate the maximum length for the original name
	maxNameLength := maxLength - len(suffix) - len(DefaultSeparator)

	// Truncate the original name if necessary
	if len(name) > maxNameLength {
		name = name[:maxNameLength]
	}

	name = strings.TrimRight(name, Separators)

	return name + DefaultSeparator + suffix
}

// TruncateAndCleanName formats a name to not exceed the provided maxLength
// and trims any trailing separator characters (-, _, .).
func TruncateAndCleanName(name string, maxLength int) string {
	// Truncate the original name if necessary
	if len(name) > maxLength {
		name = name[:maxLength]
	}

	name = strings.TrimRight(name, Separators)

	return name
}

// GenerateRandomString generates a random suffix using the full alphabet and numbers
func GenerateRandomString(length int) (string, error) {
	suffix := make([]byte, length)
	for i := range suffix {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		suffix[i] = charset[n.Int64()]
	}
	return string(suffix), nil
}
