package utils

import "strings"

func ContainsString(slice []string, search string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}

	return false
}

func ContainsStringLower(slice []string, search string) bool {
	lowerSearch := strings.ToLower(search)
	for _, item := range slice {
		if strings.ToLower(item) == lowerSearch {
			return true
		}
	}

	return false
}
