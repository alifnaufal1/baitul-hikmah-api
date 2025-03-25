package utils

import "time"

func FromTimestamp(dateString string) string {
	parsedTime, _ := time.Parse(time.RFC3339, dateString)
	return parsedTime.Format("02-01-2006 15:04:05")
}