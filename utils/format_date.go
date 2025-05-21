package utils

import (
	"fmt"
	"time"
)

func GetDate(dateString string) string {
	parsedTime, _ := time.Parse(time.RFC3339, dateString)
	return parsedTime.Format("02-01-2006")
}

func GetDateHour(dateString string) string {
	parsedTime, _ := time.Parse(time.RFC3339, dateString)
	return parsedTime.Format("02-01-2006 15:04:05")
}

func GetHumanReadableTimeDiff(dateString string) string {
	parsedTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil { return "invalid date"}

	now := time.Now()
	diff := now.Sub(parsedTime)
		
	switch {
	case diff < 0:
		if parsedTime.Hour() < now.Hour() {
			return fmt.Sprintf("%d hour ago", now.Hour()-parsedTime.Hour())
		} else if (parsedTime.Hour() == now.Hour()) {
			if parsedTime.Minute() < now.Minute() {
				return fmt.Sprintf("%d minute ago", now.Minute() - parsedTime.Minute())
			} else if parsedTime.Minute() == now.Minute() {
				if parsedTime.Second() < now.Second() {
					return fmt.Sprintf("%d seconds ago", now.Second() - parsedTime.Second())
				} else {
					return "Just now"
				}
			} else {
				return "Just now"
			}
		} else {
			return "Just now"
		}

	case diff < time.Second:
		seconds := int(diff.Seconds())
		return fmt.Sprintf("%d seconds ago", seconds)
		
	case diff < time.Minute:
		seconds := int(diff.Seconds())
		if seconds <= 1 { return "just now" }
		return fmt.Sprintf("%d seconds ago", seconds)
		
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minutes ago", minutes)

	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hours ago", hours)

	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)

	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / (24 * 7))
		return fmt.Sprintf("%d weeks ago", weeks)

	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / (24 * 30))
		return fmt.Sprintf("%d months ago", months)

	default:
		years := int(diff.Hours() / (24 * 365))
		return fmt.Sprintf("%d years ago", years)
	}
}
