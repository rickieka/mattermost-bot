package utils

// StringInSlice evaluates whether the given string is found in the given list
// and returns true/false and the index where the string was found
func StringInSlice(a string, list []string) (bool, int) {
	for i, b := range list {
		if b == a {
			return true, i
		}
	}
	return false, -1
}
