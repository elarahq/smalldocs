package utils

import (
	"regexp"
	"strings"
)

var SPACES = regexp.MustCompile(" +")

/**
 * Get formatted name from string
 */
func FormatName(str string) string {
	str = strings.ToLower(strings.Trim(str, " "))
	return SPACES.ReplaceAllString(str, "-")
}
