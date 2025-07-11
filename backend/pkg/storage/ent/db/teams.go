// Code generated by ent, DO NOT EDIT.

package db

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/teams"
)

// Teams is the model entity for the Teams schema.
type Teams struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// TeamName holds the value of the "team_name" field.
	TeamName string `json:"team_name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// JiraKeys holds the value of the "jira_keys" field.
	JiraKeys string `json:"jira_keys,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TeamsQuery when eager-loading is set.
	Edges        TeamsEdges `json:"edges"`
	selectValues sql.SelectValues
}

// TeamsEdges holds the relations/edges for other nodes in the graph.
type TeamsEdges struct {
	// Repositories holds the value of the repositories edge.
	Repositories []*Repository `json:"repositories,omitempty"`
	// Bugs holds the value of the bugs edge.
	Bugs []*Bugs `json:"bugs,omitempty"`
	// Failures holds the value of the failures edge.
	Failures []*Failure `json:"failures,omitempty"`
	// Configuration holds the value of the configuration edge.
	Configuration []*Configuration `json:"configuration,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// RepositoriesOrErr returns the Repositories value or an error if the edge
// was not loaded in eager-loading.
func (e TeamsEdges) RepositoriesOrErr() ([]*Repository, error) {
	if e.loadedTypes[0] {
		return e.Repositories, nil
	}
	return nil, &NotLoadedError{edge: "repositories"}
}

// BugsOrErr returns the Bugs value or an error if the edge
// was not loaded in eager-loading.
func (e TeamsEdges) BugsOrErr() ([]*Bugs, error) {
	if e.loadedTypes[1] {
		return e.Bugs, nil
	}
	return nil, &NotLoadedError{edge: "bugs"}
}

// FailuresOrErr returns the Failures value or an error if the edge
// was not loaded in eager-loading.
func (e TeamsEdges) FailuresOrErr() ([]*Failure, error) {
	if e.loadedTypes[2] {
		return e.Failures, nil
	}
	return nil, &NotLoadedError{edge: "failures"}
}

// ConfigurationOrErr returns the Configuration value or an error if the edge
// was not loaded in eager-loading.
func (e TeamsEdges) ConfigurationOrErr() ([]*Configuration, error) {
	if e.loadedTypes[3] {
		return e.Configuration, nil
	}
	return nil, &NotLoadedError{edge: "configuration"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Teams) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case teams.FieldTeamName, teams.FieldDescription, teams.FieldJiraKeys:
			values[i] = new(sql.NullString)
		case teams.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Teams fields.
func (t *Teams) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case teams.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				t.ID = *value
			}
		case teams.FieldTeamName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field team_name", values[i])
			} else if value.Valid {
				t.TeamName = value.String
			}
		case teams.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				t.Description = value.String
			}
		case teams.FieldJiraKeys:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field jira_keys", values[i])
			} else if value.Valid {
				t.JiraKeys = value.String
			}
		default:
			t.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Teams.
// This includes values selected through modifiers, order, etc.
func (t *Teams) Value(name string) (ent.Value, error) {
	return t.selectValues.Get(name)
}

// QueryRepositories queries the "repositories" edge of the Teams entity.
func (t *Teams) QueryRepositories() *RepositoryQuery {
	return NewTeamsClient(t.config).QueryRepositories(t)
}

// QueryBugs queries the "bugs" edge of the Teams entity.
func (t *Teams) QueryBugs() *BugsQuery {
	return NewTeamsClient(t.config).QueryBugs(t)
}

// QueryFailures queries the "failures" edge of the Teams entity.
func (t *Teams) QueryFailures() *FailureQuery {
	return NewTeamsClient(t.config).QueryFailures(t)
}

// QueryConfiguration queries the "configuration" edge of the Teams entity.
func (t *Teams) QueryConfiguration() *ConfigurationQuery {
	return NewTeamsClient(t.config).QueryConfiguration(t)
}

// Update returns a builder for updating this Teams.
// Note that you need to call Teams.Unwrap() before calling this method if this Teams
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Teams) Update() *TeamsUpdateOne {
	return NewTeamsClient(t.config).UpdateOne(t)
}

// Unwrap unwraps the Teams entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (t *Teams) Unwrap() *Teams {
	_tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("db: Teams is not a transactional entity")
	}
	t.config.driver = _tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Teams) String() string {
	var builder strings.Builder
	builder.WriteString("Teams(")
	builder.WriteString(fmt.Sprintf("id=%v, ", t.ID))
	builder.WriteString("team_name=")
	builder.WriteString(t.TeamName)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(t.Description)
	builder.WriteString(", ")
	builder.WriteString("jira_keys=")
	builder.WriteString(t.JiraKeys)
	builder.WriteByte(')')
	return builder.String()
}

// TeamsSlice is a parsable slice of Teams.
type TeamsSlice []*Teams
