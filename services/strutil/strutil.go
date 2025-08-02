package strutil

import (
	"regexp"
	"strings"
)

// IsValidName checks if the given string is a valid name.
func IsValidName(s string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	return re.MatchString(s)
}

func JoinStringsFunc(stringArray []string, stepFunc func(string) string, afterFunc func(string) string) string {
	var builder strings.Builder
	for _, str := range stringArray {
		if stepFunc != nil {
			builder.WriteString(stepFunc(str))
		} else {
			builder.WriteString(str)
		}
	}
	fullStr := builder.String()
	if afterFunc != nil {
		return afterFunc(fullStr)
	}
	return fullStr
}
