package strutil

import "regexp"

// IsValidName checks if the given string is a valid name.
func IsValidName(s string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	return re.MatchString(s)
}
