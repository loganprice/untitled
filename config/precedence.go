package config

import (
	"strings"
)

var (
	priorityOrder = [...]string{"args", "file", "env", "ext"}
)

func setPrecedence(sources map[string]ConfService) []string {
	temp := make([]string, 0, len(sources))
	for source := range sources {
		temp = append(temp, source)
	}

	return precedenceSort(temp)
}

func precedenceSort(input []string) []string {
	var output, sorted []string
	for _, prefix := range priorityOrder {
		for _, item := range input {
			if strings.HasPrefix(item, prefix) {
				sorted = append(sorted, item)
			}
		}
	}
	// TBD if this stays in or not.. Idea Maybe just log
	// that they were dropped
	output = append(difference(sorted, input), sorted...)
	return output
}

// Consider moving this to a utility folder
func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
