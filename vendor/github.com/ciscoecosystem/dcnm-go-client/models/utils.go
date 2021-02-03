package models

import "strings"

func StripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func A(data map[string]interface{}, key string, value interface{}) {

	if value != "" {
		data[key] = value
	}

	if value == "{}" {
		data[key] = ""
	}

	if value == nil {
		data[key] = ""
	}
}
