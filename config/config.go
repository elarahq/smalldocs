package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var removeSpaceRegExp = regexp.MustCompile("[ ]*")
var reduceDotRegExp = regexp.MustCompile("[.]+")
var sectionRegExp = regexp.MustCompile(`^\[(.*)\]$`)
var valueRegExp = regexp.MustCompile(`^([^=]+)=([^;#]*)$`)

/**
 * Config is returned when there is a syntax error in an INI file
 */
type ConfigError struct {
	Line   int
	Source string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("Invalid INI syntax on line %d: %s", e.Line, e.Source)
}

/**
 * A section represents a single section of an INI file
 */
type Section map[string]string

type Config struct {
	// sections in INI file
	sections map[string]Section
}

// Get named section
func (cf *Config) Section(name string) Section {
	// Init sections
	if cf.sections == nil {
		cf.sections = make(map[string]Section)
	}

	if section, ok := cf.sections[name]; ok {
		return section
	}
	return nil
}

// Returns value for a key in a section
func (cf *Config) Get(path string) (value string) {
	if sectionName, key, err := splitSection(path); err == nil {
		if section := cf.Section(sectionName); section != nil {
			value, _ = section[key]
		}
	}
	return
}

// Get Int value
func (cf *Config) Integer(path string) (value int) {
	if valueStr := cf.Get(path); valueStr != "" {
		if vInt, err := strconv.Atoi(valueStr); err == nil {
			return vInt
		}
	}
	return
}

// Get Bool value
func (cf *Config) Bool(path string) (ok bool) {
	if valueStr := cf.Get(path); valueStr != "" {
		if vBool, err := strconv.ParseBool(valueStr); err == nil {
			return vBool
		}
	}
	return
}

// Set value for a key in a section
func (cf *Config) Set(path string, value string) (err error) {
	values := make(map[string]string)
	values[path] = value
	return cf.SetMany(values)
}

// Set value for a key in a section
func (cf *Config) SetMany(values map[string]string) (err error) {
	for path, value := range values {
		if sectionName, key, err := splitSection(path); err == nil {
			section := cf.Section(sectionName)
			if section == nil {
				section = make(Section)
				cf.sections[sectionName] = section
			}
			section[key] = value
		}
	}
	return
}

// Load config file
func (cf *Config) Load(file string) (err error) {
	// open config file
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()

	// get buffer reader
	bufin := bufio.NewReader(in)

	section, lineNum, done := "", 0, false
	for !done {
		var line string
		if line, err = bufin.ReadString('\n'); err != nil {
			if err == io.EOF {
				done = true
			} else {
				return
			}
		}
		lineNum++
		line = strings.TrimSpace(line)
		// Skip comments and blank lines
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}

		if groups := valueRegExp.FindStringSubmatch(line); groups != nil {
			key, val := strings.TrimSpace(groups[1]), strings.TrimSpace(groups[2])
			cf.Section(section)[key] = val
		} else if groups := sectionRegExp.FindStringSubmatch(line); groups != nil {
			section = strings.TrimSpace(groups[1])
			if sec := cf.Section(section); sec == nil {
				cf.sections[section] = make(Section)
			}
		} else {
			return &ConfigError{lineNum, line}
		}
	}
	return nil
}

// Split section and key
func splitSection(str string) (section string, key string, err error) {
	// Remove space and multiple `.` if any
	input := removeSpaceRegExp.ReplaceAllString(str, "")
	input = reduceDotRegExp.ReplaceAllString(str, ".")

	tokens := strings.Split(input, ".")
	tokensLength := len(tokens)
	if tokensLength < 2 || len(tokens[tokensLength-1]) < 1 {
		return "", "", fmt.Errorf("error: key does not contain a section: %s", str)
	}

	key = tokens[tokensLength-1]
	section = strings.Join(tokens[:tokensLength-1], ".")
	return
}
