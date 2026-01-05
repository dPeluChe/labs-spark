package updater

import (
	"regexp"
	"strings"
)

// Version parsing patterns
var (
	// Semantic versioning: 1.2.3, v1.2.3, 1.2.3-beta, 1.2.3+build
	semverPattern = regexp.MustCompile(`v?(\d+\.\d+\.\d+[\w\-\+]*)`)

	// Major.Minor: 20.11, v16.0
	majorMinorPattern = regexp.MustCompile(`v?(\d+\.\d+)`)

	// Date-based versions: 2024.1.15
	dateVersionPattern = regexp.MustCompile(`(\d{4}\.\d+\.\d+)`)

	// Git commit hashes: abc123f (7+ hex chars)
	gitHashPattern = regexp.MustCompile(`\b([a-f0-9]{7,40})\b`)

	// Simple number version: Version 123
	simpleNumberPattern = regexp.MustCompile(`\b(\d+)\b`)
)

// CleanVersionString extracts and cleans version information from command output
func CleanVersionString(output string) string {
	if output == "" {
		return "Unknown"
	}

	// Handle special cases first
	output = strings.TrimSpace(output)
	lines := strings.Split(output, "\n")
	firstLine := lines[0]

	// Try semantic versioning first (most common)
	if version := extractPattern(firstLine, semverPattern); version != "" {
		return cleanVersion(version)
	}

	// Try major.minor pattern
	if version := extractPattern(firstLine, majorMinorPattern); version != "" {
		return cleanVersion(version)
	}

	// Try date-based versions
	if version := extractPattern(firstLine, dateVersionPattern); version != "" {
		return cleanVersion(version)
	}

	// Try git hash (for tools like oh-my-zsh)
	if version := extractPattern(firstLine, gitHashPattern); version != "" {
		return version[:7] // Return short hash
	}

	// Fallback: look for any number that looks like a version
	if version := extractPattern(firstLine, simpleNumberPattern); version != "" {
		return version
	}

	// Last resort: return first word that looks version-like
	fields := strings.Fields(firstLine)
	for _, field := range fields {
		if isVersionLike(field) {
			return cleanVersion(field)
		}
	}

	// If all else fails, return the first line trimmed
	if len(firstLine) > 30 {
		return firstLine[:30] + "…"
	}
	return firstLine
}

// extractPattern tries to extract version using a regex pattern
func extractPattern(input string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// cleanVersion removes common prefixes and normalizes
func cleanVersion(version string) string {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")

	// Remove trailing dots or dashes
	version = strings.TrimRight(version, ".-")

	return version
}

// isVersionLike checks if a string looks like it could be a version
func isVersionLike(s string) bool {
	// Must start with a digit or 'v'
	if len(s) == 0 {
		return false
	}

	firstChar := s[0]
	if firstChar >= '0' && firstChar <= '9' {
		return true
	}
	if firstChar == 'v' || firstChar == 'V' {
		if len(s) > 1 && s[1] >= '0' && s[1] <= '9' {
			return true
		}
	}

	return false
}

// ParseToolSpecificVersion handles special parsing logic for specific tools
func ParseToolSpecificVersion(toolBinary string, output string) string {
	output = strings.TrimSpace(output)

	switch toolBinary {
	case "aws":
		// aws-cli/2.22.35 Python/3.11.9 Darwin/24.0.0
		parts := strings.Fields(output)
		if len(parts) > 0 {
			// Extract from "aws-cli/2.22.35"
			if strings.Contains(parts[0], "/") {
				segments := strings.Split(parts[0], "/")
				if len(segments) > 1 {
					return segments[1]
				}
			}
		}

	case "go":
		// go version go1.23.4 darwin/arm64
		parts := strings.Fields(output)
		if len(parts) >= 3 {
			return strings.TrimPrefix(parts[2], "go")
		}

	case "python3", "python":
		// Python 3.13.1
		parts := strings.Fields(output)
		for _, part := range parts {
			if semverPattern.MatchString(part) {
				return cleanVersion(part)
			}
		}

	case "node":
		// v20.11.0
		return cleanVersion(output)

	case "npm":
		// 10.2.4
		return cleanVersion(output)

	case "docker":
		// Docker version 24.0.7, build afdd53b
		if strings.Contains(output, "version") {
			parts := strings.Fields(output)
			for i, part := range parts {
				if part == "version" && i+1 < len(parts) {
					// Remove trailing comma if present
					version := strings.TrimSuffix(parts[i+1], ",")
					return cleanVersion(version)
				}
			}
		}

	case "brew":
		// Homebrew 4.2.0
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			return cleanVersion(parts[1])
		}

	case "git":
		// git version 2.43.0
		if strings.Contains(output, "version") {
			parts := strings.Fields(output)
			for i, part := range parts {
				if part == "version" && i+1 < len(parts) {
					return cleanVersion(parts[i+1])
				}
			}
		}
	}

	// Default: use generic cleaner
	return CleanVersionString(output)
}
