package tools

func ArrayContains[T comparable](needle T, haystack []T) bool {
	for t := range haystack {
		if haystack[t] == needle {
			return true
		}
	}
	return false
}
