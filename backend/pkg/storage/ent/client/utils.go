package client

import (
	"fmt"
	"math"
	"time"

	"github.com/redhat-appstudio/quality-studio/pkg/storage"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

func convertDBError(t string, err error) error {
	if db.IsNotFound(err) {
		return storage.ErrNotFound
	}

	if db.IsConstraintError(err) {
		return storage.ErrAlreadyExists
	}

	return fmt.Errorf(t, err)
}

// getDatesBetweenRange gets the dates between a range date (including the start and end dates)
func getDatesBetweenRange(startDate, endDate string) []string {
	start, _ := time.Parse("2006-01-02 15:04:05", startDate)
	end, _ := time.Parse("2006-01-02 15:04:05", endDate)
	difference := int(start.Sub(end).Hours() / 24)
	difference = int(math.Abs(float64(difference)))

	dayArr := make([]string, 0)

	dayArr = append(dayArr, startDate)
	for i := 1; i < difference; i++ {
		dayArr = append(dayArr, start.AddDate(0, 0, +i).Format("2006-01-02 15:04:05"))
	}
	dayArr = append(dayArr, endDate)

	return dayArr
}

// isSameDay verifies if the range is pointing the same day
func isSameDay(startDate, endDate string) bool {
	start, _ := time.Parse("2006-01-02 15:04:05", startDate)
	end, _ := time.Parse("2006-01-02 15:04:05", endDate)

	if start.Year() == end.Year() && start.Month() == end.Month() &&
		start.Day() == end.Day() {
		return true
	}
	return false
}

// getRange gets the start and end date
func getRange(i int, day string, dayArr []string) (string, string) {
	t, _ := time.Parse("2006-01-02 15:04:05", day)
	y, m, dd := t.Date()
	start := ""
	end := ""

	if i == 0 { // first day
		start = day
		end = fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd)
	} else {
		if i == len(dayArr)-1 { // last day
			start = fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd)
			end = day
		} else { // middle days
			start = fmt.Sprintf("%04d-%02d-%02d 00:00:00", y, m, dd)
			end = fmt.Sprintf("%04d-%02d-%02d 23:59:59", y, m, dd)
		}
	}

	return start, end
}
