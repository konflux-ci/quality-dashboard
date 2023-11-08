package client

import (
	"context"
	"fmt"
	"math"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	failureV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/failure/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/api/server/router/prow"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/failure"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowjobs"
)

// CreateFailure ...
func (d *Database) CreateFailure(f failureV1Alpha1.Failure, team_id uuid.UUID) error {
	failureAlreadyExists := d.client.Failure.Query().
		Where(failure.JiraKey(f.JiraKey)).
		ExistX(context.TODO())
	if failureAlreadyExists {
		_, err := d.client.Failure.Update().
			Where(predicate.Failure(failure.JiraKey(f.JiraKey))).
			SetJiraStatus(f.JiraStatus).
			SetErrorMessage(f.ErrorMessage).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to update failure: %w", err)
		}
	} else {
		fr, err := d.client.Failure.Create().
			SetJiraKey(f.JiraKey).
			SetJiraStatus(f.JiraStatus).
			SetErrorMessage(f.ErrorMessage).
			Save(context.TODO())
		if err != nil {
			return convertDBError("create failure: %w", err)
		}

		_, err = d.client.Teams.UpdateOneID(team_id).AddFailures(fr).Save(context.TODO())
		if err != nil {
			return convertDBError("create failure: %w", err)
		}
	}

	return nil
}

// GetFrequency gets the frequency of a error message in the range date time provided
func (d *Database) GetFrequency(team *db.Teams, errorMessage, startDate, endDate string) (float64, error) {
	var frequency, total, occurrences float64 = 0, 0, 0

	repositories, err := d.client.Teams.QueryRepositories(team).All(context.TODO())
	if err != nil {
		return 0, convertDBError("list repositories: %w", err)
	}

	for _, repo := range repositories {
		if repo.RepositoryName == "e2e-tests" || repo.RepositoryName == "infra-deployments" {
			prowJobs, err := d.client.Repository.QueryProwJobs(repo).Select().
				Where(prowjobs.State(string(prow.FailureState))).
				Where(func(s *sql.Selector) { // "created_at BETWEEN ? AND 2022-08-17", "2022-08-16"
					s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", startDate, endDate)))
				}).All(context.TODO())
			if err != nil {
				fmt.Printf("failed to get prow jobs for repo %s: %v", err, repo.RepositoryName)
			}

			for _, prowJob := range prowJobs {
				total++

				if prowJob.E2eFailedTestMessages != nil {
					if strings.Contains(*prowJob.E2eFailedTestMessages, errorMessage) {
						occurrences++
					}
				}

				if prowJob.BuildErrorLogs != nil {
					if strings.Contains(*prowJob.BuildErrorLogs, errorMessage) {
						occurrences++
					}
				}
			}
		}
	}

	if occurrences != 0 {
		frequency = (occurrences * 100) / total
		frequency = math.Round(frequency*100) / 100
	}

	return frequency, nil
}

func (d *Database) GetFailuresByDate(team *db.Teams, startDate, endDate string) ([]*failureV1Alpha1.Failure, error) {
	dbFailures, err := d.client.Teams.QueryFailures(team).All(context.TODO())
	if err != nil {
		return nil, convertDBError("list failures: %w", err)
	}

	failures := make([]*failureV1Alpha1.Failure, 0)

	for _, failure := range dbFailures {
		frequency, err := d.GetFrequency(team, failure.ErrorMessage, startDate, endDate)
		if err != nil {
			fmt.Printf("failed to get frequency: %v", err)
			continue
		}

		failures = append(failures, &failureV1Alpha1.Failure{
			TeamID:       team.ID,
			TeamName:     team.TeamName,
			JiraID:       failure.ID,
			JiraKey:      failure.JiraKey,
			JiraStatus:   failure.JiraStatus,
			ErrorMessage: failure.ErrorMessage,
			Frequency:    frequency,
		})
	}

	return failures, nil
}

func (d *Database) GetAllFailures(team *db.Teams) ([]*db.Failure, error) {
	failures, err := d.client.Teams.QueryFailures(team).All(context.TODO())
	if err != nil {
		return nil, convertDBError("list failures: %w", err)
	}

	return failures, nil
}

func (d *Database) DeleteFailure(teamID, failureID uuid.UUID) error {
	_, err := d.client.Teams.UpdateOneID(teamID).RemoveFailureIDs(failureID).Save(context.TODO())

	if err != nil {
		return convertDBError("create failure: %w", err)
	}

	return nil
}
