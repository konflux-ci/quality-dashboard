// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/users"
)

// UsersDelete is the builder for deleting a Users entity.
type UsersDelete struct {
	config
	hooks    []Hook
	mutation *UsersMutation
}

// Where appends a list predicates to the UsersDelete builder.
func (ud *UsersDelete) Where(ps ...predicate.Users) *UsersDelete {
	ud.mutation.Where(ps...)
	return ud
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ud *UsersDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, ud.sqlExec, ud.mutation, ud.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ud *UsersDelete) ExecX(ctx context.Context) int {
	n, err := ud.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ud *UsersDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(users.Table, sqlgraph.NewFieldSpec(users.FieldID, field.TypeUUID))
	if ps := ud.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ud.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ud.mutation.done = true
	return affected, err
}

// UsersDeleteOne is the builder for deleting a single Users entity.
type UsersDeleteOne struct {
	ud *UsersDelete
}

// Where appends a list predicates to the UsersDelete builder.
func (udo *UsersDeleteOne) Where(ps ...predicate.Users) *UsersDeleteOne {
	udo.ud.mutation.Where(ps...)
	return udo
}

// Exec executes the deletion query.
func (udo *UsersDeleteOne) Exec(ctx context.Context) error {
	n, err := udo.ud.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{users.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (udo *UsersDeleteOne) ExecX(ctx context.Context) {
	if err := udo.Exec(ctx); err != nil {
		panic(err)
	}
}
