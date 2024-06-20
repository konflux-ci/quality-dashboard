// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/failure"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/predicate"
)

// FailureDelete is the builder for deleting a Failure entity.
type FailureDelete struct {
	config
	hooks    []Hook
	mutation *FailureMutation
}

// Where appends a list predicates to the FailureDelete builder.
func (fd *FailureDelete) Where(ps ...predicate.Failure) *FailureDelete {
	fd.mutation.Where(ps...)
	return fd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (fd *FailureDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, FailureMutation](ctx, fd.sqlExec, fd.mutation, fd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (fd *FailureDelete) ExecX(ctx context.Context) int {
	n, err := fd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (fd *FailureDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: failure.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: failure.FieldID,
			},
		},
	}
	if ps := fd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, fd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	fd.mutation.done = true
	return affected, err
}

// FailureDeleteOne is the builder for deleting a single Failure entity.
type FailureDeleteOne struct {
	fd *FailureDelete
}

// Where appends a list predicates to the FailureDelete builder.
func (fdo *FailureDeleteOne) Where(ps ...predicate.Failure) *FailureDeleteOne {
	fdo.fd.mutation.Where(ps...)
	return fdo
}

// Exec executes the deletion query.
func (fdo *FailureDeleteOne) Exec(ctx context.Context) error {
	n, err := fdo.fd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{failure.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (fdo *FailureDeleteOne) ExecX(ctx context.Context) {
	if err := fdo.Exec(ctx); err != nil {
		panic(err)
	}
}
