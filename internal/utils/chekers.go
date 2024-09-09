package utils

func CheckString(str string, allowedValues []string) bool {
	for _, value := range allowedValues {
		if str == value {
			return true
		}
	}
	return false
}
