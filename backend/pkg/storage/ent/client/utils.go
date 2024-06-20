package client

import (
	"fmt"
	"math"
	"time"

	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	util "github.com/konflux-ci/quality-dashboard/pkg/utils"
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
func getRangesInISO(startDate, endDate string) []string {
	start, _ := time.Parse(time.RFC3339, startDate)
	end, _ := time.Parse(time.RFC3339, endDate)
	difference := int(math.Ceil(end.Sub(start).Hours() / 24))

	dayArr := make([]string, 0)

	for i := 0; i < difference; i++ {
		dayArr = append(dayArr, start.AddDate(0, 0, +i).Format(time.RFC3339))
	}
	dayArr = append(dayArr, endDate)

	return dayArr
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

// getDaysBetweenDates gets the number of days between two dates
func getDaysBetweenDates(firstDate, secondDate time.Time) float64 {
	diff := secondDate.Sub(firstDate).Hours()

	// diff in days
	diff = diff / 24

	// round diff to 2 decimal places
	diff = math.Round(diff*100) / 100

	return diff
}

// getWorkingDays only includes business days by excluding weekend days
func getWorkingDays(fromDate, toDate time.Time) float64 {
	var workingDays float64 = 0
	previousDate := fromDate
	nextDate := previousDate.Add(time.Hour * 24)

	for {
		if previousDate.Equal(toDate) || previousDate.Equal(nextDate) || previousDate.After(toDate) || previousDate.After(nextDate) {
			break
		}
		if previousDate.Weekday() != 6 && previousDate.Weekday() != 0 {
			if toDate.Before(nextDate) {
				workingDays += getDaysBetweenDates(previousDate, toDate)
			} else {
				workingDays += getDaysBetweenDates(previousDate, nextDate)
			}
		}
		previousDate = nextDate
		nextDate = previousDate.Add(time.Hour * 24)
	}

	return workingDays
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CalculatePercentage(x float64, y float64) float64 {
	result := util.RoundTo(x/y*100, 2)

	if math.IsNaN(result) {
		return 0
	}

	return result
}
