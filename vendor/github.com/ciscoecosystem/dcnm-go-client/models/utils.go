package models

import (
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/container"
)

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

func IsService(path string) bool {
	return strings.Contains(path, "elastic-service") || strings.Contains(path, "elasticservice")
}

func G(cont *container.Container, key string) string {
	return StripQuotes(cont.S(key).String())
}