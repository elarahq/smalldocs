package utils

import (
	"regexp"
	"strings"
)

// Regexp to remove special character
var SLUGIFY = regexp.MustCompile(`[^a-zA-Z0-9\.-]+`)
var TITLE = regexp.MustCompile(`[^a-zA-Z0-9 \.-]+`)

//
// Get formatted name from string
//
func Slug(str string) string {
	str = SLUGIFY.ReplaceAllString(str, " ")
	str = strings.ToLower(strings.Trim(str, " "))
	str = strings.Replace(str, " ", "-", -1)
	return str
}

//
// Formate title
//
func Title(str string) string {
	str = TITLE.ReplaceAllString(str, "")
	return strings.Trim(str, " ")
}

// Get matched parameters
func GetMatchedParams(str string, re *regexp.Regexp) map[string]string {
	match := re.FindStringSubmatch(str)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		result[name] = match[i]
	}
	return result
}
