package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func InsertAt(filePath, search string, before bool, newContent string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create a scanner to read line by line
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var newLines []string

	// Scan each line
	for scanner.Scan() {
		line := scanner.Text()
		if !before {
			newLines = append(newLines, line)
		}

		// If we find the comment, append the new content on the next line
		if strings.TrimSpace(line) == search {
			newLines = append(newLines, newContent)
		} else {
			fmt.Println(strings.TrimSpace(line))
		}

		if before {
			newLines = append(newLines, line)
		}
	}

	// Write back to file
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}
