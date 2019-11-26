package main

// ContainsString sees if a given string is in a slice of strings.
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
