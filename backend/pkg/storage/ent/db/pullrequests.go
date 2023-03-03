// Code generated by ent, DO NOT EDIT.

package db

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/pullrequests"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

// PullRequests is the model entity for the PullRequests schema.
type PullRequests struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// PrID holds the value of the "pr_id" field.
	PrID uuid.UUID `json:"pr_id,omitempty"`
	// RepositoryName holds the value of the "repository_name" field.
	RepositoryName string `json:"repository_name,omitempty"`
	// RepositoryOrganization holds the value of the "repository_organization" field.
	RepositoryOrganization string `json:"repository_organization,omitempty"`
	// Number holds the value of the "number" field.
	Number int `json:"number,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// ClosedAt holds the value of the "closed_at" field.
	ClosedAt time.Time `json:"closed_at,omitempty"`
	// MergedAt holds the value of the "merged_at" field.
	MergedAt time.Time `json:"merged_at,omitempty"`
	// State holds the value of the "state" field.
	State string `json:"state,omitempty"`
	// Author holds the value of the "author" field.
	Author string `json:"author,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PullRequestsQuery when eager-loading is set.
	Edges          PullRequestsEdges `json:"edges"`
	repository_prs *uuid.UUID
}

// PullRequestsEdges holds the relations/edges for other nodes in the graph.
type PullRequestsEdges struct {
	// Prs holds the value of the prs edge.
	Prs *Repository `json:"prs,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PrsOrErr returns the Prs value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PullRequestsEdges) PrsOrErr() (*Repository, error) {
	if e.loadedTypes[0] {
		if e.Prs == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: repository.Label}
		}
		return e.Prs, nil
	}
	return nil, &NotLoadedError{edge: "prs"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PullRequests) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case pullrequests.FieldID, pullrequests.FieldNumber:
			values[i] = new(sql.NullInt64)
		case pullrequests.FieldRepositoryName, pullrequests.FieldRepositoryOrganization, pullrequests.FieldState, pullrequests.FieldAuthor, pullrequests.FieldTitle:
			values[i] = new(sql.NullString)
		case pullrequests.FieldCreatedAt, pullrequests.FieldClosedAt, pullrequests.FieldMergedAt:
			values[i] = new(sql.NullTime)
		case pullrequests.FieldPrID:
			values[i] = new(uuid.UUID)
		case pullrequests.ForeignKeys[0]: // repository_prs
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			return nil, fmt.Errorf("unexpected column %q for type PullRequests", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PullRequests fields.
func (pr *PullRequests) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case pullrequests.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pr.ID = int(value.Int64)
		case pullrequests.FieldPrID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field pr_id", values[i])
			} else if value != nil {
				pr.PrID = *value
			}
		case pullrequests.FieldRepositoryName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field repository_name", values[i])
			} else if value.Valid {
				pr.RepositoryName = value.String
			}
		case pullrequests.FieldRepositoryOrganization:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field repository_organization", values[i])
			} else if value.Valid {
				pr.RepositoryOrganization = value.String
			}
		case pullrequests.FieldNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field number", values[i])
			} else if value.Valid {
				pr.Number = int(value.Int64)
			}
		case pullrequests.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		case pullrequests.FieldClosedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field closed_at", values[i])
			} else if value.Valid {
				pr.ClosedAt = value.Time
			}
		case pullrequests.FieldMergedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field merged_at", values[i])
			} else if value.Valid {
				pr.MergedAt = value.Time
			}
		case pullrequests.FieldState:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field state", values[i])
			} else if value.Valid {
				pr.State = value.String
			}
		case pullrequests.FieldAuthor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field author", values[i])
			} else if value.Valid {
				pr.Author = value.String
			}
		case pullrequests.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				pr.Title = value.String
			}
		case pullrequests.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field repository_prs", values[i])
			} else if value.Valid {
				pr.repository_prs = new(uuid.UUID)
				*pr.repository_prs = *value.S.(*uuid.UUID)
			}
		}
	}
	return nil
}

// QueryPrs queries the "prs" edge of the PullRequests entity.
func (pr *PullRequests) QueryPrs() *RepositoryQuery {
	return NewPullRequestsClient(pr.config).QueryPrs(pr)
}

// Update returns a builder for updating this PullRequests.
// Note that you need to call PullRequests.Unwrap() before calling this method if this PullRequests
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *PullRequests) Update() *PullRequestsUpdateOne {
	return NewPullRequestsClient(pr.config).UpdateOne(pr)
}

// Unwrap unwraps the PullRequests entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *PullRequests) Unwrap() *PullRequests {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("db: PullRequests is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *PullRequests) String() string {
	var builder strings.Builder
	builder.WriteString("PullRequests(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("pr_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.PrID))
	builder.WriteString(", ")
	builder.WriteString("repository_name=")
	builder.WriteString(pr.RepositoryName)
	builder.WriteString(", ")
	builder.WriteString("repository_organization=")
	builder.WriteString(pr.RepositoryOrganization)
	builder.WriteString(", ")
	builder.WriteString("number=")
	builder.WriteString(fmt.Sprintf("%v", pr.Number))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("closed_at=")
	builder.WriteString(pr.ClosedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("merged_at=")
	builder.WriteString(pr.MergedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("state=")
	builder.WriteString(pr.State)
	builder.WriteString(", ")
	builder.WriteString("author=")
	builder.WriteString(pr.Author)
	builder.WriteString(", ")
	builder.WriteString("title=")
	builder.WriteString(pr.Title)
	builder.WriteByte(')')
	return builder.String()
}

// PullRequestsSlice is a parsable slice of PullRequests.
type PullRequestsSlice []*PullRequests

func (pr PullRequestsSlice) config(cfg config) {
	for _i := range pr {
		pr[_i].config = cfg
	}
}
