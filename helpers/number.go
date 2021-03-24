package helpers

import "strconv"

// Convert integer number to a string
func Int64ToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}
