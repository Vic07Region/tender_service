package utils

import "strconv"

func CheckString(str string, allowedValues []string) bool {
	for _, value := range allowedValues {
		if str == value {
			return true
		}
	}
	return false
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
