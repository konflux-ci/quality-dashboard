package client

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/redhat-appstudio/quality-studio/api/apis/prow/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

func (d *Database) GetSuitesFailureFrequency(gitOrg string, repoName string, startDate string, endDate string) ([]v1alpha1.SuitesFailureFrequency, error) {
	var suitesFailure []v1alpha1.SuitesFailureFrequency
	var stats []v1alpha1.SuitesFailureFrequency

	repository, err := d.client.Repository.Query().
		Where(repository.RepositoryName(repoName)).Where(repository.GitOrganization(gitOrg)).Only(context.TODO())
	if err != nil {
		fmt.Println(err)
		return nil, convertDBError("get repository: %w", err)
	}

	err = d.client.Repository.QueryProwSuites(repository).
		Where(func(s *sql.Selector) { // "merged_at BETWEEN ? AND 2022-08-17", "2022-08-16"
			//s.Where(sql.ExprP(fmt.Sprintf("created_at BETWEEN '%s' AND '%s'", "2023-10-25", fmt.Sprintf("%s:23:59:59", "2023-11-13"))))
		}).
		GroupBy(prowsuites.FieldName, prowsuites.FieldStatus, prowsuites.FieldSuiteName).
		Aggregate(db.Count()).
		Scan(context.Background(), &suitesFailure)

	if err != nil {
		return nil, convertDBError("get suites: %w", err)
	}

	for _, occ := range suitesFailure {
		var msg []v1alpha1.ErrorMessage

		suite, err := d.client.Repository.QueryProwSuites(repository).Where(prowsuites.Name(occ.Name)).All(context.Background())
		if err != nil {
			continue
		}
		for _, c := range suite {
			msg = append(msg, v1alpha1.ErrorMessage{
				JobId:   c.JobID,
				JobURL:  c.JobURL,
				Message: *c.ErrorMessage,
			})
		}

		stats = append(stats, v1alpha1.SuitesFailureFrequency{
			Name:         occ.Name,
			Count:        occ.Count,
			Status:       occ.Status,
			SuiteName:    occ.SuiteName,
			ErrorMessage: msg,
		})

	}

	return stats, nil
}
