package datastructs

// StringInSlice returns true/false depending on whether or not the string exists in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
