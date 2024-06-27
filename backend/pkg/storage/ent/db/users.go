// Code generated by ent, DO NOT EDIT.

package db

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/users"
)

// Users is the model entity for the Users schema.
type Users struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// UserEmail holds the value of the "user_email" field.
	UserEmail string `json:"user_email,omitempty"`
	// Config holds the value of the "config" field.
	Config string `json:"config,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Users) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case users.FieldUserEmail, users.FieldConfig:
			values[i] = new(sql.NullString)
		case users.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Users", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Users fields.
func (u *Users) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case users.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				u.ID = *value
			}
		case users.FieldUserEmail:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_email", values[i])
			} else if value.Valid {
				u.UserEmail = value.String
			}
		case users.FieldConfig:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field config", values[i])
			} else if value.Valid {
				u.Config = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Users.
// Note that you need to call Users.Unwrap() before calling this method if this Users
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *Users) Update() *UsersUpdateOne {
	return NewUsersClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the Users entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *Users) Unwrap() *Users {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("db: Users is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *Users) String() string {
	var builder strings.Builder
	builder.WriteString("Users(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("user_email=")
	builder.WriteString(u.UserEmail)
	builder.WriteString(", ")
	builder.WriteString("config=")
	builder.WriteString(u.Config)
	builder.WriteByte(')')
	return builder.String()
}

// UsersSlice is a parsable slice of Users.
type UsersSlice []*Users

func (u UsersSlice) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
