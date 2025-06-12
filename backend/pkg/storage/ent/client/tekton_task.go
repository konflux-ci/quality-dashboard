package client

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/tektontasks"
)

// All tasks will be optionally linked to the same ProwJob if jobID is provided.
func (d *Database) CreateTektonTasksBulk(tasks []*db.TektonTasks, jobID *int) error {
	if len(tasks) == 0 {
		return nil
	}

	for _, task := range tasks {
		builder := d.client.TektonTasks.Create().
			SetTaskName(task.TaskName).
			SetDurationSeconds(task.DurationSeconds).
			SetStatus(task.Status)

		if jobID != nil {
			builder = builder.SetTektonTasksID(*jobID)
		}

		if _, err := builder.Save(context.Background()); err != nil {
			continue
		}
	}

	return nil
}

// GetTasksMetrics calculates the success/failure percentage and average duration for all tasks
// within a given date range and returns them as a slice of TaskMetrics structs.
// This version handles duration strings with suffixes (e.g., "123s").
func (d *Database) GetTasksMetrics(startDate string, endDate string) ([]storage.TaskMetrics, error) {
	ctx := context.Background()

	tasks, err := d.client.TektonTasks.Query().
		Where(
			tektontasks.TaskNameIn("provision-kind-cluster", "deploy-konflux"),
			func(s *sql.Selector) { // Apply the date range filter directly.
				// s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
			},
		).
		Select(tektontasks.FieldTaskName, tektontasks.FieldStatus, tektontasks.FieldDurationSeconds).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to query tekton tasks: %w", err)
	}

	if len(tasks) == 0 {
		return []storage.TaskMetrics{}, nil
	}

	durationTotals := make(map[string]float64)
	totalCounts := make(map[string]int)
	successCounts := make(map[string]int)
	failureCounts := make(map[string]int)

	for _, task := range tasks {
		totalCounts[task.TaskName]++

		if task.Status == "success" {
			successCounts[task.TaskName]++
		} else if task.Status == "failure" {
			failureCounts[task.TaskName]++
		}

		durationStr := strings.TrimSuffix(task.DurationSeconds, "s")
		duration, err := strconv.ParseFloat(durationStr, 64)
		if err == nil {
			durationTotals[task.TaskName] += duration
		}
	}

	result := make([]storage.TaskMetrics, 0, len(totalCounts))
	for taskName, count := range totalCounts {
		var avgDuration float64
		if count > 0 {
			avgDuration = durationTotals[taskName] / float64(count)
		}

		var successPercentage float64
		if count > 0 {
			successPercentage = (float64(successCounts[taskName]) / float64(count)) * 100
		}

		var failurePercentage float64
		if count > 0 {
			failurePercentage = (float64(failureCounts[taskName]) / float64(count)) * 100
		}

		result = append(result, storage.TaskMetrics{
			TaskName:          taskName,
			TotalRuns:         count,
			SuccessPercentage: math.Round(successPercentage*100) / 100,
			FailurePercentage: math.Round(failurePercentage*100) / 100,
			AverageDuration:   math.Round(avgDuration*100) / 100,
		})
	}

	return result, nil
}
