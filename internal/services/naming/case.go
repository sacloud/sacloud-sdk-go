// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package naming

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// SnakeCase converts a string to snake_case.
func SnakeCase(s string) string {
	return toDelimitedLower(s, '_')
}

// KebabCase converts a string to kebab-case.
func KebabCase(s string) string {
	return toDelimitedLower(s, '-')
}

// PascalCase converts a string to PascalCase.
func PascalCase(s string) string {
	return toCamelCase(s, true)
}

// CamelCase converts a string to camelCase.
func CamelCase(s string) string {
	return toCamelCase(s, false)
}

// FirstRuneToLower lowercases the first rune in a string.
func FirstRuneToLower(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if !unicode.IsUpper(r) {
		return s
	}
	return string(unicode.ToLower(r)) + s[size:]
}

func toDelimitedLower(s string, delimiter rune) string {
	runes := []rune(s)
	var result strings.Builder
	lastOutputIsDelimiter := false
	for i, r := range runes {
		switch {
		case isConnector(r):
			result.WriteRune(delimiter)
			lastOutputIsDelimiter = true
		case unicode.IsPunct(r) && r != '[' && r != ']':
			continue
		default:
			if i > 0 && result.Len() > 0 && needsDelimiter(runes, i) && !lastOutputIsDelimiter {
				result.WriteRune(delimiter)
			}
			result.WriteRune(unicode.ToLower(r))
			lastOutputIsDelimiter = false
		}
	}
	return result.String()
}

func needsDelimiter(runes []rune, i int) bool {
	current := runes[i]
	previous := runes[i-1]
	if isConnector(previous) || unicode.IsPunct(previous) {
		return true
	}

	if unicode.IsUpper(current) {
		if unicode.IsLower(previous) {
			return true
		}
		if unicode.IsUpper(previous) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
			return true
		}
	}

	if unicode.IsNumber(previous) && unicode.IsUpper(current) {
		return true
	}

	if unicode.IsLetter(previous) && unicode.IsNumber(current) {
		end := i
		for end < len(runes) && unicode.IsNumber(runes[end]) {
			end++
		}
		return end < len(runes) && unicode.IsLower(runes[end])
	}

	return false
}

func toCamelCase(s string, upperFirst bool) string {
	runes := []rune(s)
	var result strings.Builder
	capitalize := upperFirst
	for _, r := range runes {
		if isConnector(r) {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(unicode.ToUpper(r))
			capitalize = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func isConnector(r rune) bool {
	return r == '-' || r == '_' || unicode.IsSpace(r)
}
