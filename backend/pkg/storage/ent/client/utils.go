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
