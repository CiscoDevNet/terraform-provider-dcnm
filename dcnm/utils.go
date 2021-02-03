package dcnm

import (
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/container"
)

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func cleanJsonString(data string) (*container.Container, error) {
	data = strings.ReplaceAll(data, "\\", "")

	cont, err := container.ParseJSON([]byte(data))
	if err != nil {
		return nil, err
	}
	return cont, nil
}

func listToString(data interface{}) string {
	values := data.([]interface{})

	strList := make([]string, 0, 1)
	for _, val := range values {
		strList = append(strList, val.(string))
	}

	return strings.Join(strList, ",")
}

func stringToList(data string) []string {
	strList := make([]string, 0, 1)

	strs := strings.Split(data, ",")
	for _, val := range strs {
		strList = append(strList, strings.Trim(val, " "))
	}

	return strList
}

func interfaceToStrList(data interface{}) []string {
	values := data.([]interface{})

	strList := make([]string, 0, 1)
	for _, val := range values {
		strList = append(strList, val.(string))
	}

	return strList
}

func compareStrLists(first, second []string) bool {
	sort.Strings(first)
	sort.Strings(second)

	if reflect.DeepEqual(first, second) {
		return true
	}
	return false
}
