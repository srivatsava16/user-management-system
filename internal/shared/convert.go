package shared

import (
	"reflect"
	"strconv"
	"user-management-system/internal/models"
)

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func EqualIntArrays(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func CompareJSON(a, b models.JSON) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		bVal, exists := b[k]
		if !exists {
			return false
		}
		if !reflect.DeepEqual(v, bVal) {
			return false
		}
	}

	return true
}
