package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func processMarkdown(lines []string) ([]string, bool) {
	if len(lines) == 0 {
		return lines, false
	}

	contentStartIndex := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				contentStartIndex = i + 1
				break
			}
		}
	}

	frontmatter := lines[:contentStartIndex]
	content := lines[contentStartIndex:]

	var processed []string
	h2Seen := false
	hasSeparator := true

	for _, line := range content {
		isH2 := strings.HasPrefix(line, "## ")
		isHR := strings.TrimSpace(line) == "---"

		if isH2 {
			if h2Seen && !hasSeparator {
				for len(processed) > 0 && strings.TrimSpace(processed[len(processed)-1]) == "" {
					processed = processed[:len(processed)-1]
				}
				processed = append(processed, "\n", "---\n", "\n")
			}
			h2Seen = true
			hasSeparator = false
			processed = append(processed, line)
		} else if isHR {
			hasSeparator = true
			processed = append(processed, line)
		} else {
			processed = append(processed, line)
		}
	}

	finalLines := append(frontmatter, processed...)

	if len(finalLines) != len(lines) {
		return finalLines, true
	}
	for i := range finalLines {
		if finalLines[i] != lines[i] {
			return finalLines, true
		}
	}
	return lines, false
}

func main() {
	verbose := flag.Bool("v", false, "Print filenames as they are processed")
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"blog"}
	}

	for _, path := range paths {
		err := filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(fp) != ".md" {
				return nil
			}

			if *verbose {
				fmt.Printf("Processing: %s\n", fp)
			}

			file, err := os.Open(fp)
			if err != nil {
				return err
			}

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text()+"\n")
			}
			file.Close()

			if err := scanner.Err(); err != nil {
				return err
			}

			updatedLines, changed := processMarkdown(lines)
			if changed {
				err = os.WriteFile(fp, []byte(strings.Join(updatedLines, "")), 0644)
				if err != nil {
					return err
				}
				if *verbose {
					fmt.Printf("  Updated %s\n", fp)
				}
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing path %s: %v\n", path, err)
		}
	}
}
