package security

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides input sanitization methods
type Sanitizer struct{}

// NewSanitizer creates a new sanitizer instance
func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

// SanitizeString removes dangerous characters and sanitizes input
func (s *Sanitizer) SanitizeString(input string, maxLength int) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")

	// Remove other control characters except newline, carriage return, and tab
	sanitized = removeControlCharacters(sanitized)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Limit length
	if maxLength > 0 && len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	return sanitized
}

// SanitizeHTML escapes HTML special characters to prevent XSS
func (s *Sanitizer) SanitizeHTML(input string) string {
	return html.EscapeString(input)
}

// SanitizeSearchQuery sanitizes search input to prevent ReDoS and SQL injection
func (s *Sanitizer) SanitizeSearchQuery(query string) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(query, "\x00", "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Remove or escape special regex characters that could cause ReDoS
	// We'll escape: . * + ? ^ $ ( ) [ ] { } | \
	regexSpecialChars := []string{
		"\\", ".", "*", "+", "?", "^", "$",
		"(", ")", "[", "]", "{", "}", "|",
	}

	for _, char := range regexSpecialChars {
		sanitized = strings.ReplaceAll(sanitized, char, "\\"+char)
	}

	// Limit length to prevent DoS
	maxSearchLength := 200
	if len(sanitized) > maxSearchLength {
		sanitized = sanitized[:maxSearchLength]
	}

	return sanitized
}

// SanitizeSQLLike sanitizes input for SQL LIKE queries
func (s *Sanitizer) SanitizeSQLLike(input string) string {
	// Remove null bytes
	sanitized := strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Escape SQL LIKE special characters
	sanitized = strings.ReplaceAll(sanitized, "\\", "\\\\")
	sanitized = strings.ReplaceAll(sanitized, "%", "\\%")
	sanitized = strings.ReplaceAll(sanitized, "_", "\\_")

	// Limit length
	maxLength := 200
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	return sanitized
}

// StripTags removes HTML/XML tags from input
func (s *Sanitizer) StripTags(input string) string {
	// Remove HTML tags using regex
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

// SanitizeFilename sanitizes filename to prevent path traversal
func (s *Sanitizer) SanitizeFilename(filename string) string {
	// Remove directory separators and path traversal attempts
	sanitized := strings.ReplaceAll(filename, "..", "")
	sanitized = strings.ReplaceAll(sanitized, "/", "")
	sanitized = strings.ReplaceAll(sanitized, "\\", "")
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	// Remove control characters
	sanitized = removeControlCharacters(sanitized)

	// Limit length
	maxLength := 255
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	return strings.TrimSpace(sanitized)
}

// ContainsDangerousContent checks if input contains potentially dangerous content
func (s *Sanitizer) ContainsDangerousContent(input string) bool {
	input = strings.ToLower(input)

	// Check for script tags
	if strings.Contains(input, "<script") || strings.Contains(input, "</script>") {
		return true
	}

	// Check for javascript: protocol
	if strings.Contains(input, "javascript:") {
		return true
	}

	// Check for data: protocol (can be used for XSS)
	if strings.Contains(input, "data:text/html") {
		return true
	}

	// Check for event handlers
	dangerousPatterns := []string{
		"onerror", "onload", "onclick", "onmouseover",
		"onfocus", "onblur", "onchange", "onsubmit",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}

	// Check for iframe
	if strings.Contains(input, "<iframe") {
		return true
	}

	// Check for object/embed tags
	if strings.Contains(input, "<object") || strings.Contains(input, "<embed") {
		return true
	}

	return false
}

// removeControlCharacters removes control characters except newline, tab, and carriage return
func removeControlCharacters(input string) string {
	var result strings.Builder
	for _, char := range input {
		// Keep printable characters, newline (\n), tab (\t), and carriage return (\r)
		if char == '\n' || char == '\t' || char == '\r' || !unicode.IsControl(char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}
