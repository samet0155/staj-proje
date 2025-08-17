package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"staj/rules"
	"staj/vulnerability"
)

var configExtensions = []string{".conf", ".cfg", ".env", ".ini", ".json", ".yaml", ".yml"}

func ScanDirectory(root string) ([]vulnerability.Vulnerability, error) {
	var results []vulnerability.Vulnerability

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk error for %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if isConfigFile(path) {
			vulns, err := ScanFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "file scan error %s: %v\n", path, err)
				return nil
			}
			results = append(results, vulns...)
		}
		return nil
	})

	return results, err
}

func isConfigFile(filename string) bool {
	lower := strings.ToLower(filename)
	for _, ext := range configExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func ScanFile(filePath string) ([]vulnerability.Vulnerability, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return scanJSONFile(filePath)
	case ".yaml", ".yml":
		return scanYAMLFile(filePath)
	default:
		return scanTextFile(filePath)
	}
}

func scanTextFile(filePath string) ([]vulnerability.Vulnerability, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var vulns []vulnerability.Vulnerability
	sc := bufio.NewScanner(f)
	lineNumber := 1
	for sc.Scan() {
		line := sc.Text()
		lowerLine := strings.ToLower(line)
		for _, r := range rules.Rules {
			if strings.Contains(lowerLine, strings.ToLower(r.Pattern)) {
				vulns = append(vulns, vulnerability.Vulnerability{
					File: filePath,
					Rule: r.Description,
					Line: lineNumber,
				})
			}
		}
		lineNumber++
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return vulns, nil
}

func scanJSONFile(filePath string) ([]vulnerability.Vulnerability, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	var vulns []vulnerability.Vulnerability
	scanJSONRecursive(filePath, content, &vulns)
	return vulns, nil
}

func scanJSONRecursive(filePath string, node interface{}, vulns *[]vulnerability.Vulnerability) {
	switch val := node.(type) {
	case map[string]interface{}:
		for k, v := range val {
			checkRules(filePath, k, vulns)
			checkRules(filePath, fmt.Sprintf("%v", v), vulns)
			scanJSONRecursive(filePath, v, vulns)
		}
	case []interface{}:
		for _, elem := range val {
			scanJSONRecursive(filePath, elem, vulns)
		}
	}
}

func scanYAMLFile(filePath string) ([]vulnerability.Vulnerability, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var content interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	var vulns []vulnerability.Vulnerability
	scanYAMLRecursive(filePath, content, &vulns)
	return vulns, nil
}

func scanYAMLRecursive(filePath string, node interface{}, vulns *[]vulnerability.Vulnerability) {
	switch val := node.(type) {
	case map[string]interface{}:
		for k, v := range val {
			checkRules(filePath, k, vulns)
			checkRules(filePath, fmt.Sprintf("%v", v), vulns)
			scanYAMLRecursive(filePath, v, vulns)
		}
	case []interface{}:
		for _, elem := range val {
			scanYAMLRecursive(filePath, elem, vulns)
		}
	}
}

func checkRules(filePath, text string, vulns *[]vulnerability.Vulnerability) {
	lowerText := strings.ToLower(text)
	for _, r := range rules.Rules {
		if strings.Contains(lowerText, strings.ToLower(r.Pattern)) {
			*vulns = append(*vulns, vulnerability.Vulnerability{
				File: filePath,
				Rule: r.Description,
				Line: 0,
			})
		}
	}
}
