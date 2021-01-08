package utilities

import (
	"strings"
)

// Returns a lexicographic check to see if a date falls within a required date range
func CompareDates(lower string, upper string, compare string) bool {
	strings.ReplaceAll(lower, "-", "")

	return (lower <= compare) && (compare <= upper)
}
