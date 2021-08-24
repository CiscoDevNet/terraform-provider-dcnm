package dcnm

import (
	"fmt"
	"hash/crc32"
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

func toStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func getErrorFromContainer(cont *container.Container, err error) error {
	if contErr := stripQuotes(cont.S("error", "detail").String()); cont != nil && contErr != "null" {
		return fmt.Errorf(contErr)
	}
	return err
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

func contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a.(string) == e {
			return true
		}
	}
	return false
}

func hashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func setDifference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
