// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowsuites"
)

// ProwSuitesDelete is the builder for deleting a ProwSuites entity.
type ProwSuitesDelete struct {
	config
	hooks    []Hook
	mutation *ProwSuitesMutation
}

// Where appends a list predicates to the ProwSuitesDelete builder.
func (psd *ProwSuitesDelete) Where(ps ...predicate.ProwSuites) *ProwSuitesDelete {
	psd.mutation.Where(ps...)
	return psd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (psd *ProwSuitesDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, psd.sqlExec, psd.mutation, psd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (psd *ProwSuitesDelete) ExecX(ctx context.Context) int {
	n, err := psd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (psd *ProwSuitesDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(prowsuites.Table, sqlgraph.NewFieldSpec(prowsuites.FieldID, field.TypeInt))
	if ps := psd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, psd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	psd.mutation.done = true
	return affected, err
}

// ProwSuitesDeleteOne is the builder for deleting a single ProwSuites entity.
type ProwSuitesDeleteOne struct {
	psd *ProwSuitesDelete
}

// Where appends a list predicates to the ProwSuitesDelete builder.
func (psdo *ProwSuitesDeleteOne) Where(ps ...predicate.ProwSuites) *ProwSuitesDeleteOne {
	psdo.psd.mutation.Where(ps...)
	return psdo
}

// Exec executes the deletion query.
func (psdo *ProwSuitesDeleteOne) Exec(ctx context.Context) error {
	n, err := psdo.psd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{prowsuites.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (psdo *ProwSuitesDeleteOne) ExecX(ctx context.Context) {
	if err := psdo.Exec(ctx); err != nil {
		panic(err)
	}
}
