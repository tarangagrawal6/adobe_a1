package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// outlineEntry represents a single entry in the outline
type outlineEntry struct {
	Level string `json:"level"`
	Text  string `json:"text"`
	Page  int    `json:"page"`
}

// pdfData represents the structured data extracted from a PDF
type pdfData struct {
	Title   string         `json:"title"`
	Outline []outlineEntry `json:"outline"`
}

// processPDF processes a PDF file and returns its structured data
func processPDF(pdfPath, _ string) (pdfData, error) {
	var data pdfData

	// Simulate pdftotext by reading OCR content (in real use, replace with actual pdftotext call)
	cmd := exec.Command("pdftotext", "-layout", pdfPath, "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return pdfData{}, fmt.Errorf("failed to extract text from %s: %v", pdfPath, err)
	}

	// Split text into pages (assuming \f as page break, adjust for OCR input)
	pages := strings.Split(out.String(), "\f")
	if len(pages) == 0 {
		return pdfData{}, fmt.Errorf("no pages found in %s", pdfPath)
	}

	// Extract title dynamically
	data.Title = extractTitle(pages)

	// Extract outline from content
	data.Outline = extractOutline(pages)

	return data, nil
}

// extractTitle finds a prominent title from the pages
func extractTitle(pages []string) string {
	if len(pages) == 0 {
		return ""
	}

	// Look at the first page for a title
	lines := strings.Split(pages[0], "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 5 && len(line) < 100 { // Heuristic: non-empty, reasonable length
			if isProminent(line) {
				return cleanText(line)
			}
		}
	}

	// Fallback: Look for repeated text across pages (e.g., headers)
	repeated := findRepeatedText(pages)
	if repeated != "" {
		return cleanText(repeated)
	}

	return "Untitled"
}

// isProminent checks if a line looks like a title (e.g., mostly uppercase or long enough)
func isProminent(line string) bool {
	upperCount := 0
	for _, r := range line {
		if unicode.IsUpper(r) {
			upperCount++
		}
	}
	return upperCount > len(line)/2 || len(line) > 20
}

// findRepeatedText detects text that repeats across pages (e.g., headers)
func findRepeatedText(pages []string) string {
	if len(pages) < 2 {
		return ""
	}
	counts := make(map[string]int)
	for _, page := range pages {
		lines := strings.Split(page, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && len(line) > 5 && len(line) < 100 {
				counts[line]++
			}
		}
	}
	for text, count := range counts {
		if count > len(pages)/2 { // Appears on more than half the pages
			return text
		}
	}
	return ""
}

// extractOutline detects headings from all pages
func extractOutline(pages []string) []outlineEntry {
	var outline []outlineEntry
	seen := make(map[string]bool) // Avoid duplicates

	// Regex for potential headings (e.g., "Section 1", all caps, or standalone phrases)
	reHeading := regexp.MustCompile(`^(Section|Chapter|Part|Appendix)\s+[A-Z\d]+|^[A-Z][A-Za-z\s\-\:]{5,}$`)

	for pageNum, page := range pages {
		scanner := bufio.NewScanner(strings.NewReader(page))
		prevIndent := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || isNoise(line, pages) {
				continue
			}

			if reHeading.MatchString(line) && !seen[line] {
				level := determineLevel(line, prevIndent)
				outline = append(outline, outlineEntry{
					Level: level,
					Text:  cleanText(line),
					Page:  pageNum,
				})
				seen[line] = true
				prevIndent = countIndent(line)
			}
		}
	}
	return outline
}

// determineLevel infers heading level based on text and indentation
func determineLevel(line string, prevIndent int) string {
	indent := countIndent(line)
	if indent > prevIndent {
		return "H2" // Sub-level if indented more
	} else if indent < prevIndent {
		return "H1" // Higher level if less indented
	}

	// Use text clues
	if strings.Contains(line, ":") {
		return "H3" // Colons often indicate sub-sections
	}
	if regexp.MustCompile(`^[A-Z\s]+$`).MatchString(line) {
		return "H1" // All caps for top-level
	}
	return "H2" // Default mid-level
}

// countIndent counts leading spaces for level inference
func countIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// isNoise detects repetitive or irrelevant text
func isNoise(line string, pages []string) bool {
	count := 0
	for _, page := range pages {
		if strings.Contains(page, line) {
			count++
		}
	}
	return count > len(pages)*2 // Appears too often (e.g., "TOPJUMP")
}

// cleanText removes extra spaces and special characters
func cleanText(text string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(text, " "))
}

func main() {
	const inputDir = "/app/input"
	const outputDir = "/app/output"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Error reading input directory: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	var errs []error
	var mu sync.Mutex

	for _, file := range files {
		if !file.IsDir() && strings.ToLower(filepath.Ext(file.Name())) == ".pdf" {
			wg.Add(1)
			go func(file os.DirEntry) {
				defer wg.Done()
				jsonFileName := strings.TrimSuffix(file.Name(), ".pdf") + ".json"
				jsonPath := filepath.Join(outputDir, jsonFileName)
				pdfPath := filepath.Join(inputDir, file.Name())

				data, err := processPDF(pdfPath, file.Name())
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("error processing %s: %v", file.Name(), err))
					mu.Unlock()
					return
				}

				jsonData, err := json.MarshalIndent(data, "", "    ")
				if err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("error marshaling JSON for %s: %v", file.Name(), err))
					mu.Lock()
					return
				}

				if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("error writing JSON file for %s: %v", file.Name(), err))
					mu.Unlock()
					return
				}

				fmt.Printf("Processed %s -> %s\n", file.Name(), jsonFileName)
			}(file)
		}
	}

	wg.Wait()

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}