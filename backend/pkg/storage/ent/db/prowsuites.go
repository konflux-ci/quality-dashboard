// Code generated by ent, DO NOT EDIT.

package db

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowsuites"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

// ProwSuites is the model entity for the ProwSuites schema.
type ProwSuites struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// JobID holds the value of the "job_id" field.
	JobID string `json:"job_id,omitempty"`
	// JobURL holds the value of the "job_url" field.
	JobURL string `json:"job_url,omitempty"`
	// JobName holds the value of the "job_name" field.
	JobName string `json:"job_name,omitempty"`
	// SuiteName holds the value of the "suite_name" field.
	SuiteName string `json:"suite_name,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty"`
	// ErrorMessage holds the value of the "error_message" field.
	ErrorMessage *string `json:"error_message,omitempty"`
	// ExternalServicesImpact holds the value of the "external_services_impact" field.
	ExternalServicesImpact *bool `json:"external_services_impact,omitempty"`
	// Time holds the value of the "time" field.
	Time float64 `json:"time,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProwSuitesQuery when eager-loading is set.
	Edges                  ProwSuitesEdges `json:"edges"`
	repository_prow_suites *string
	selectValues           sql.SelectValues
}

// ProwSuitesEdges holds the relations/edges for other nodes in the graph.
type ProwSuitesEdges struct {
	// ProwSuites holds the value of the prow_suites edge.
	ProwSuites *Repository `json:"prow_suites,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ProwSuitesOrErr returns the ProwSuites value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProwSuitesEdges) ProwSuitesOrErr() (*Repository, error) {
	if e.ProwSuites != nil {
		return e.ProwSuites, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: repository.Label}
	}
	return nil, &NotLoadedError{edge: "prow_suites"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProwSuites) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case prowsuites.FieldExternalServicesImpact:
			values[i] = new(sql.NullBool)
		case prowsuites.FieldTime:
			values[i] = new(sql.NullFloat64)
		case prowsuites.FieldID:
			values[i] = new(sql.NullInt64)
		case prowsuites.FieldJobID, prowsuites.FieldJobURL, prowsuites.FieldJobName, prowsuites.FieldSuiteName, prowsuites.FieldName, prowsuites.FieldStatus, prowsuites.FieldErrorMessage:
			values[i] = new(sql.NullString)
		case prowsuites.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case prowsuites.ForeignKeys[0]: // repository_prow_suites
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProwSuites fields.
func (ps *ProwSuites) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case prowsuites.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ps.ID = int(value.Int64)
		case prowsuites.FieldJobID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field job_id", values[i])
			} else if value.Valid {
				ps.JobID = value.String
			}
		case prowsuites.FieldJobURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field job_url", values[i])
			} else if value.Valid {
				ps.JobURL = value.String
			}
		case prowsuites.FieldJobName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field job_name", values[i])
			} else if value.Valid {
				ps.JobName = value.String
			}
		case prowsuites.FieldSuiteName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field suite_name", values[i])
			} else if value.Valid {
				ps.SuiteName = value.String
			}
		case prowsuites.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ps.Name = value.String
			}
		case prowsuites.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ps.Status = value.String
			}
		case prowsuites.FieldErrorMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field error_message", values[i])
			} else if value.Valid {
				ps.ErrorMessage = new(string)
				*ps.ErrorMessage = value.String
			}
		case prowsuites.FieldExternalServicesImpact:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field external_services_impact", values[i])
			} else if value.Valid {
				ps.ExternalServicesImpact = new(bool)
				*ps.ExternalServicesImpact = value.Bool
			}
		case prowsuites.FieldTime:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field time", values[i])
			} else if value.Valid {
				ps.Time = value.Float64
			}
		case prowsuites.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ps.CreatedAt = new(time.Time)
				*ps.CreatedAt = value.Time
			}
		case prowsuites.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field repository_prow_suites", values[i])
			} else if value.Valid {
				ps.repository_prow_suites = new(string)
				*ps.repository_prow_suites = value.String
			}
		default:
			ps.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProwSuites.
// This includes values selected through modifiers, order, etc.
func (ps *ProwSuites) Value(name string) (ent.Value, error) {
	return ps.selectValues.Get(name)
}

// QueryProwSuites queries the "prow_suites" edge of the ProwSuites entity.
func (ps *ProwSuites) QueryProwSuites() *RepositoryQuery {
	return NewProwSuitesClient(ps.config).QueryProwSuites(ps)
}

// Update returns a builder for updating this ProwSuites.
// Note that you need to call ProwSuites.Unwrap() before calling this method if this ProwSuites
// was returned from a transaction, and the transaction was committed or rolled back.
func (ps *ProwSuites) Update() *ProwSuitesUpdateOne {
	return NewProwSuitesClient(ps.config).UpdateOne(ps)
}

// Unwrap unwraps the ProwSuites entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ps *ProwSuites) Unwrap() *ProwSuites {
	_tx, ok := ps.config.driver.(*txDriver)
	if !ok {
		panic("db: ProwSuites is not a transactional entity")
	}
	ps.config.driver = _tx.drv
	return ps
}

// String implements the fmt.Stringer.
func (ps *ProwSuites) String() string {
	var builder strings.Builder
	builder.WriteString("ProwSuites(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ps.ID))
	builder.WriteString("job_id=")
	builder.WriteString(ps.JobID)
	builder.WriteString(", ")
	builder.WriteString("job_url=")
	builder.WriteString(ps.JobURL)
	builder.WriteString(", ")
	builder.WriteString("job_name=")
	builder.WriteString(ps.JobName)
	builder.WriteString(", ")
	builder.WriteString("suite_name=")
	builder.WriteString(ps.SuiteName)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(ps.Name)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(ps.Status)
	builder.WriteString(", ")
	if v := ps.ErrorMessage; v != nil {
		builder.WriteString("error_message=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := ps.ExternalServicesImpact; v != nil {
		builder.WriteString("external_services_impact=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("time=")
	builder.WriteString(fmt.Sprintf("%v", ps.Time))
	builder.WriteString(", ")
	if v := ps.CreatedAt; v != nil {
		builder.WriteString("created_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// ProwSuitesSlice is a parsable slice of ProwSuites.
type ProwSuitesSlice []*ProwSuites
