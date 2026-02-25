//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func slugify(title string) string {
	// Remove quotes if present
	title = strings.Trim(title, `"'`)
	title = strings.ToLower(title)

	// Remove apostrophes
	title = strings.ReplaceAll(title, "'", "")
	title = strings.ReplaceAll(title, "’", "")

	// Replace non-alphanumeric characters with spaces
	re := regexp.MustCompile(`[^a-z0-9]+`)
	title = re.ReplaceAllString(title, " ")

	// Trim spaces and replace remaining spaces with hyphens
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, " ", "-")

	return title
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run check_title.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	mismatches := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			if checkFile(path) {
				mismatches++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	if mismatches == 0 {
		fmt.Println("All file names match their titles.")
	} else {
		fmt.Printf("Found %d mismatch(es).\n", mismatches)
		os.Exit(1)
	}
}

func checkFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	title := ""

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "---" {
			if inFrontmatter {
				break // End of frontmatter
			} else {
				inFrontmatter = true
				continue
			}
		}

		if inFrontmatter && strings.HasPrefix(line, "title:") {
			title = strings.TrimPrefix(line, "title:")
			title = strings.TrimSpace(title)
			break
		}
	}

	if title == "" {
		return false
	}

	expectedSlug := slugify(title)
	baseName := filepath.Base(path)
	fileNameSlug := strings.TrimSuffix(baseName, ".md")

	if expectedSlug != fileNameSlug {
		fmt.Printf("Mismatch in %s\n", path)
		fmt.Printf("  Title from file  : %s\n", title)
		fmt.Printf("  Expected filename: %s.md\n", expectedSlug)
		fmt.Printf("  Actual filename  : %s\n\n", baseName)
		return true
	}
	return false
}
