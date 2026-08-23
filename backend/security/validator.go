package security

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator provides input validation methods
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateString validates a string with length constraints
func (v *Validator) ValidateString(field, value string, minLen, maxLen int, required bool) error {
	if required && strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: "is required"}
	}

	if !required && value == "" {
		return nil
	}

	length := len(strings.TrimSpace(value))
	if length < minLen {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must be at least %d characters", minLen)}
	}

	if maxLen > 0 && length > maxLen {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must not exceed %d characters", maxLen)}
	}

	return nil
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "is required"}
	}

	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "is invalid"}
	}

	if len(email) > 255 {
		return &ValidationError{Field: "email", Message: "must not exceed 255 characters"}
	}

	return nil
}

// ValidateUsername validates username format
func (v *Validator) ValidateUsername(username string) error {
	if username == "" {
		return &ValidationError{Field: "username", Message: "is required"}
	}

	// Length check
	if len(username) < 3 || len(username) > 50 {
		return &ValidationError{Field: "username", Message: "must be between 3 and 50 characters"}
	}

	// No spaces
	if strings.Contains(username, " ") {
		return &ValidationError{Field: "username", Message: "cannot contain spaces"}
	}

	// No Turkish characters
	turkishChars := "ğüşıöçĞÜŞİÖÇ"
	for _, char := range username {
		if strings.ContainsRune(turkishChars, char) {
			return &ValidationError{Field: "username", Message: "cannot contain Turkish characters"}
		}
	}

	// Must start with alphanumeric
	if !unicode.IsLetter(rune(username[0])) && !unicode.IsNumber(rune(username[0])) {
		return &ValidationError{Field: "username", Message: "must start with a letter or number"}
	}

	// Only alphanumeric and underscore/hyphen allowed
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !usernameRegex.MatchString(username) {
		return &ValidationError{Field: "username", Message: "can only contain letters, numbers, underscores, and hyphens"}
	}

	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return &ValidationError{Field: "password", Message: "is required"}
	}

	// Minimum 8 characters
	if len(password) < 8 {
		return &ValidationError{Field: "password", Message: "must be at least 8 characters"}
	}

	// Maximum 128 characters (prevent DoS)
	if len(password) > 128 {
		return &ValidationError{Field: "password", Message: "must not exceed 128 characters"}
	}

	// At least one uppercase letter
	hasUpper := false
	// At least one lowercase letter
	hasLower := false
	// At least one digit
	hasDigit := false
	// At least one special character
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return &ValidationError{Field: "password", Message: "must contain at least one uppercase letter"}
	}

	if !hasLower {
		return &ValidationError{Field: "password", Message: "must contain at least one lowercase letter"}
	}

	if !hasDigit {
		return &ValidationError{Field: "password", Message: "must contain at least one number"}
	}

	if !hasSpecial {
		return &ValidationError{Field: "password", Message: "must contain at least one special character"}
	}

	return nil
}

// ValidateInteger validates integer within range
func (v *Validator) ValidateInteger(field string, value, min, max int) error {
	if value < min {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must be at least %d", min)}
	}

	if max > 0 && value > max {
		return &ValidationError{Field: field, Message: fmt.Sprintf("must not exceed %d", max)}
	}

	return nil
}

// ValidateRoleID validates role_id
func (v *Validator) ValidateRoleID(roleID int) error {
	validRoles := map[int]bool{1: true, 2: true, 3: true}
	if !validRoles[roleID] {
		return &ValidationError{Field: "role_id", Message: "must be 1 (Admin), 2 (Member), or 3 (User)"}
	}
	return nil
}

// ValidateCommunity validates community value
func (v *Validator) ValidateCommunity(community int) error {
	if community != 1 && community != 2 {
		return &ValidationError{Field: "community", Message: "must be 1 (private) or 2 (public)"}
	}
	return nil
}

// ValidateBookCommunity validates a book's visibility level.
// Kitaba özel: ValidateCommunity şiirler tarafından da kullanılıyor ve
// şiir filtresi 3'ü tanımıyor, o yüzden ayrı tutuluyor.
func (v *Validator) ValidateBookCommunity(community int) error {
	if community != 1 && community != 2 && community != 3 {
		return &ValidationError{Field: "community", Message: "must be 1 (private), 2 (public) or 3 (selected users)"}
	}
	return nil
}

// NormalizeISBN, girilen ISBN/EAN'ı temizler, doğrular ve ISBN-13'e çevirir.
// Checksum doğrulaması dış API çağrısı açılmadan önceki en ucuz filtredir:
// rastgele 13 hane burada büyük olasılıkla elenir.
func (v *Validator) NormalizeISBN(raw string) (string, error) {
	cleaned := strings.Map(func(r rune) rune {
		if r == '-' || r == ' ' || r == '.' {
			return -1
		}
		return unicode.ToUpper(r)
	}, strings.TrimSpace(raw))

	switch len(cleaned) {
	case 10:
		if !isValidISBN10(cleaned) {
			return "", &ValidationError{Field: "isbn", Message: "invalid ISBN-10 checksum"}
		}
		return isbn10To13(cleaned), nil
	case 13:
		for _, r := range cleaned {
			if r < '0' || r > '9' {
				return "", &ValidationError{Field: "isbn", Message: "ISBN-13 must contain digits only"}
			}
		}
		if !strings.HasPrefix(cleaned, "978") && !strings.HasPrefix(cleaned, "979") {
			return "", &ValidationError{Field: "isbn", Message: "ISBN-13 must start with 978 or 979"}
		}
		if !isValidISBN13(cleaned) {
			return "", &ValidationError{Field: "isbn", Message: "invalid ISBN-13 checksum"}
		}
		return cleaned, nil
	default:
		return "", &ValidationError{Field: "isbn", Message: "ISBN must be 10 or 13 characters"}
	}
}

// isValidISBN10 — mod 11, son hane 'X' olabilir.
func isValidISBN10(s string) bool {
	sum := 0
	for i := 0; i < 10; i++ {
		c := s[i]
		var digit int
		switch {
		case c >= '0' && c <= '9':
			digit = int(c - '0')
		case c == 'X' && i == 9:
			digit = 10
		default:
			return false
		}
		sum += digit * (10 - i)
	}
	return sum%11 == 0
}

// isValidISBN13 — mod 10, ağırlıklar 1 ve 3 dönüşümlü.
func isValidISBN13(s string) bool {
	sum := 0
	for i := 0; i < 13; i++ {
		digit := int(s[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	return sum%10 == 0
}

// isbn10To13, 978 öneki ekleyip kontrol hanesini yeniden hesaplar.
func isbn10To13(isbn10 string) string {
	core := "978" + isbn10[:9]
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(core[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	check := (10 - sum%10) % 10
	return fmt.Sprintf("%s%d", core, check)
}

// ValidatePermission validates permission value
func (v *Validator) ValidatePermission(permission int) error {
	if permission < 1 || permission > 3 {
		return &ValidationError{Field: "permission", Message: "must be 1 (Admin), 2 (User), or 3 (Guest)"}
	}
	return nil
}

// ValidateID validates that ID is positive
func (v *Validator) ValidateID(field string, id uint) error {
	if id == 0 {
		return &ValidationError{Field: field, Message: "must be a valid ID"}
	}
	return nil
}
