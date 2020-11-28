package paypal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Loads a discord token from filename
func loadToken(filename string) (string, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Scan for token
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.TrimSpace(s) != "" {
			return s, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// Token not found
	return "", fmt.Errorf("%v did not contain a token", filename)
}
