// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowsuites"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
)

// ProwSuitesUpdate is the builder for updating ProwSuites entities.
type ProwSuitesUpdate struct {
	config
	hooks    []Hook
	mutation *ProwSuitesMutation
}

// Where appends a list predicates to the ProwSuitesUpdate builder.
func (psu *ProwSuitesUpdate) Where(ps ...predicate.ProwSuites) *ProwSuitesUpdate {
	psu.mutation.Where(ps...)
	return psu
}

// SetJobID sets the "job_id" field.
func (psu *ProwSuitesUpdate) SetJobID(s string) *ProwSuitesUpdate {
	psu.mutation.SetJobID(s)
	return psu
}

// SetNillableJobID sets the "job_id" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableJobID(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetJobID(*s)
	}
	return psu
}

// SetJobURL sets the "job_url" field.
func (psu *ProwSuitesUpdate) SetJobURL(s string) *ProwSuitesUpdate {
	psu.mutation.SetJobURL(s)
	return psu
}

// SetNillableJobURL sets the "job_url" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableJobURL(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetJobURL(*s)
	}
	return psu
}

// SetJobName sets the "job_name" field.
func (psu *ProwSuitesUpdate) SetJobName(s string) *ProwSuitesUpdate {
	psu.mutation.SetJobName(s)
	return psu
}

// SetNillableJobName sets the "job_name" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableJobName(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetJobName(*s)
	}
	return psu
}

// SetSuiteName sets the "suite_name" field.
func (psu *ProwSuitesUpdate) SetSuiteName(s string) *ProwSuitesUpdate {
	psu.mutation.SetSuiteName(s)
	return psu
}

// SetNillableSuiteName sets the "suite_name" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableSuiteName(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetSuiteName(*s)
	}
	return psu
}

// SetName sets the "name" field.
func (psu *ProwSuitesUpdate) SetName(s string) *ProwSuitesUpdate {
	psu.mutation.SetName(s)
	return psu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableName(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetName(*s)
	}
	return psu
}

// SetStatus sets the "status" field.
func (psu *ProwSuitesUpdate) SetStatus(s string) *ProwSuitesUpdate {
	psu.mutation.SetStatus(s)
	return psu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableStatus(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetStatus(*s)
	}
	return psu
}

// SetErrorMessage sets the "error_message" field.
func (psu *ProwSuitesUpdate) SetErrorMessage(s string) *ProwSuitesUpdate {
	psu.mutation.SetErrorMessage(s)
	return psu
}

// SetNillableErrorMessage sets the "error_message" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableErrorMessage(s *string) *ProwSuitesUpdate {
	if s != nil {
		psu.SetErrorMessage(*s)
	}
	return psu
}

// ClearErrorMessage clears the value of the "error_message" field.
func (psu *ProwSuitesUpdate) ClearErrorMessage() *ProwSuitesUpdate {
	psu.mutation.ClearErrorMessage()
	return psu
}

// SetExternalServicesImpact sets the "external_services_impact" field.
func (psu *ProwSuitesUpdate) SetExternalServicesImpact(b bool) *ProwSuitesUpdate {
	psu.mutation.SetExternalServicesImpact(b)
	return psu
}

// SetNillableExternalServicesImpact sets the "external_services_impact" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableExternalServicesImpact(b *bool) *ProwSuitesUpdate {
	if b != nil {
		psu.SetExternalServicesImpact(*b)
	}
	return psu
}

// ClearExternalServicesImpact clears the value of the "external_services_impact" field.
func (psu *ProwSuitesUpdate) ClearExternalServicesImpact() *ProwSuitesUpdate {
	psu.mutation.ClearExternalServicesImpact()
	return psu
}

// SetTime sets the "time" field.
func (psu *ProwSuitesUpdate) SetTime(f float64) *ProwSuitesUpdate {
	psu.mutation.ResetTime()
	psu.mutation.SetTime(f)
	return psu
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableTime(f *float64) *ProwSuitesUpdate {
	if f != nil {
		psu.SetTime(*f)
	}
	return psu
}

// AddTime adds f to the "time" field.
func (psu *ProwSuitesUpdate) AddTime(f float64) *ProwSuitesUpdate {
	psu.mutation.AddTime(f)
	return psu
}

// SetCreatedAt sets the "created_at" field.
func (psu *ProwSuitesUpdate) SetCreatedAt(t time.Time) *ProwSuitesUpdate {
	psu.mutation.SetCreatedAt(t)
	return psu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableCreatedAt(t *time.Time) *ProwSuitesUpdate {
	if t != nil {
		psu.SetCreatedAt(*t)
	}
	return psu
}

// ClearCreatedAt clears the value of the "created_at" field.
func (psu *ProwSuitesUpdate) ClearCreatedAt() *ProwSuitesUpdate {
	psu.mutation.ClearCreatedAt()
	return psu
}

// SetProwSuitesID sets the "prow_suites" edge to the Repository entity by ID.
func (psu *ProwSuitesUpdate) SetProwSuitesID(id string) *ProwSuitesUpdate {
	psu.mutation.SetProwSuitesID(id)
	return psu
}

// SetNillableProwSuitesID sets the "prow_suites" edge to the Repository entity by ID if the given value is not nil.
func (psu *ProwSuitesUpdate) SetNillableProwSuitesID(id *string) *ProwSuitesUpdate {
	if id != nil {
		psu = psu.SetProwSuitesID(*id)
	}
	return psu
}

// SetProwSuites sets the "prow_suites" edge to the Repository entity.
func (psu *ProwSuitesUpdate) SetProwSuites(r *Repository) *ProwSuitesUpdate {
	return psu.SetProwSuitesID(r.ID)
}

// Mutation returns the ProwSuitesMutation object of the builder.
func (psu *ProwSuitesUpdate) Mutation() *ProwSuitesMutation {
	return psu.mutation
}

// ClearProwSuites clears the "prow_suites" edge to the Repository entity.
func (psu *ProwSuitesUpdate) ClearProwSuites() *ProwSuitesUpdate {
	psu.mutation.ClearProwSuites()
	return psu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (psu *ProwSuitesUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, psu.sqlSave, psu.mutation, psu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psu *ProwSuitesUpdate) SaveX(ctx context.Context) int {
	affected, err := psu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (psu *ProwSuitesUpdate) Exec(ctx context.Context) error {
	_, err := psu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psu *ProwSuitesUpdate) ExecX(ctx context.Context) {
	if err := psu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (psu *ProwSuitesUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(prowsuites.Table, prowsuites.Columns, sqlgraph.NewFieldSpec(prowsuites.FieldID, field.TypeInt))
	if ps := psu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := psu.mutation.JobID(); ok {
		_spec.SetField(prowsuites.FieldJobID, field.TypeString, value)
	}
	if value, ok := psu.mutation.JobURL(); ok {
		_spec.SetField(prowsuites.FieldJobURL, field.TypeString, value)
	}
	if value, ok := psu.mutation.JobName(); ok {
		_spec.SetField(prowsuites.FieldJobName, field.TypeString, value)
	}
	if value, ok := psu.mutation.SuiteName(); ok {
		_spec.SetField(prowsuites.FieldSuiteName, field.TypeString, value)
	}
	if value, ok := psu.mutation.Name(); ok {
		_spec.SetField(prowsuites.FieldName, field.TypeString, value)
	}
	if value, ok := psu.mutation.Status(); ok {
		_spec.SetField(prowsuites.FieldStatus, field.TypeString, value)
	}
	if value, ok := psu.mutation.ErrorMessage(); ok {
		_spec.SetField(prowsuites.FieldErrorMessage, field.TypeString, value)
	}
	if psu.mutation.ErrorMessageCleared() {
		_spec.ClearField(prowsuites.FieldErrorMessage, field.TypeString)
	}
	if value, ok := psu.mutation.ExternalServicesImpact(); ok {
		_spec.SetField(prowsuites.FieldExternalServicesImpact, field.TypeBool, value)
	}
	if psu.mutation.ExternalServicesImpactCleared() {
		_spec.ClearField(prowsuites.FieldExternalServicesImpact, field.TypeBool)
	}
	if value, ok := psu.mutation.Time(); ok {
		_spec.SetField(prowsuites.FieldTime, field.TypeFloat64, value)
	}
	if value, ok := psu.mutation.AddedTime(); ok {
		_spec.AddField(prowsuites.FieldTime, field.TypeFloat64, value)
	}
	if value, ok := psu.mutation.CreatedAt(); ok {
		_spec.SetField(prowsuites.FieldCreatedAt, field.TypeTime, value)
	}
	if psu.mutation.CreatedAtCleared() {
		_spec.ClearField(prowsuites.FieldCreatedAt, field.TypeTime)
	}
	if psu.mutation.ProwSuitesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   prowsuites.ProwSuitesTable,
			Columns: []string{prowsuites.ProwSuitesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(repository.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psu.mutation.ProwSuitesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   prowsuites.ProwSuitesTable,
			Columns: []string{prowsuites.ProwSuitesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(repository.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, psu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{prowsuites.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	psu.mutation.done = true
	return n, nil
}

// ProwSuitesUpdateOne is the builder for updating a single ProwSuites entity.
type ProwSuitesUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProwSuitesMutation
}

// SetJobID sets the "job_id" field.
func (psuo *ProwSuitesUpdateOne) SetJobID(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetJobID(s)
	return psuo
}

// SetNillableJobID sets the "job_id" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableJobID(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetJobID(*s)
	}
	return psuo
}

// SetJobURL sets the "job_url" field.
func (psuo *ProwSuitesUpdateOne) SetJobURL(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetJobURL(s)
	return psuo
}

// SetNillableJobURL sets the "job_url" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableJobURL(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetJobURL(*s)
	}
	return psuo
}

// SetJobName sets the "job_name" field.
func (psuo *ProwSuitesUpdateOne) SetJobName(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetJobName(s)
	return psuo
}

// SetNillableJobName sets the "job_name" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableJobName(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetJobName(*s)
	}
	return psuo
}

// SetSuiteName sets the "suite_name" field.
func (psuo *ProwSuitesUpdateOne) SetSuiteName(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetSuiteName(s)
	return psuo
}

// SetNillableSuiteName sets the "suite_name" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableSuiteName(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetSuiteName(*s)
	}
	return psuo
}

// SetName sets the "name" field.
func (psuo *ProwSuitesUpdateOne) SetName(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetName(s)
	return psuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableName(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetName(*s)
	}
	return psuo
}

// SetStatus sets the "status" field.
func (psuo *ProwSuitesUpdateOne) SetStatus(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetStatus(s)
	return psuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableStatus(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetStatus(*s)
	}
	return psuo
}

// SetErrorMessage sets the "error_message" field.
func (psuo *ProwSuitesUpdateOne) SetErrorMessage(s string) *ProwSuitesUpdateOne {
	psuo.mutation.SetErrorMessage(s)
	return psuo
}

// SetNillableErrorMessage sets the "error_message" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableErrorMessage(s *string) *ProwSuitesUpdateOne {
	if s != nil {
		psuo.SetErrorMessage(*s)
	}
	return psuo
}

// ClearErrorMessage clears the value of the "error_message" field.
func (psuo *ProwSuitesUpdateOne) ClearErrorMessage() *ProwSuitesUpdateOne {
	psuo.mutation.ClearErrorMessage()
	return psuo
}

// SetExternalServicesImpact sets the "external_services_impact" field.
func (psuo *ProwSuitesUpdateOne) SetExternalServicesImpact(b bool) *ProwSuitesUpdateOne {
	psuo.mutation.SetExternalServicesImpact(b)
	return psuo
}

// SetNillableExternalServicesImpact sets the "external_services_impact" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableExternalServicesImpact(b *bool) *ProwSuitesUpdateOne {
	if b != nil {
		psuo.SetExternalServicesImpact(*b)
	}
	return psuo
}

// ClearExternalServicesImpact clears the value of the "external_services_impact" field.
func (psuo *ProwSuitesUpdateOne) ClearExternalServicesImpact() *ProwSuitesUpdateOne {
	psuo.mutation.ClearExternalServicesImpact()
	return psuo
}

// SetTime sets the "time" field.
func (psuo *ProwSuitesUpdateOne) SetTime(f float64) *ProwSuitesUpdateOne {
	psuo.mutation.ResetTime()
	psuo.mutation.SetTime(f)
	return psuo
}

// SetNillableTime sets the "time" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableTime(f *float64) *ProwSuitesUpdateOne {
	if f != nil {
		psuo.SetTime(*f)
	}
	return psuo
}

// AddTime adds f to the "time" field.
func (psuo *ProwSuitesUpdateOne) AddTime(f float64) *ProwSuitesUpdateOne {
	psuo.mutation.AddTime(f)
	return psuo
}

// SetCreatedAt sets the "created_at" field.
func (psuo *ProwSuitesUpdateOne) SetCreatedAt(t time.Time) *ProwSuitesUpdateOne {
	psuo.mutation.SetCreatedAt(t)
	return psuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableCreatedAt(t *time.Time) *ProwSuitesUpdateOne {
	if t != nil {
		psuo.SetCreatedAt(*t)
	}
	return psuo
}

// ClearCreatedAt clears the value of the "created_at" field.
func (psuo *ProwSuitesUpdateOne) ClearCreatedAt() *ProwSuitesUpdateOne {
	psuo.mutation.ClearCreatedAt()
	return psuo
}

// SetProwSuitesID sets the "prow_suites" edge to the Repository entity by ID.
func (psuo *ProwSuitesUpdateOne) SetProwSuitesID(id string) *ProwSuitesUpdateOne {
	psuo.mutation.SetProwSuitesID(id)
	return psuo
}

// SetNillableProwSuitesID sets the "prow_suites" edge to the Repository entity by ID if the given value is not nil.
func (psuo *ProwSuitesUpdateOne) SetNillableProwSuitesID(id *string) *ProwSuitesUpdateOne {
	if id != nil {
		psuo = psuo.SetProwSuitesID(*id)
	}
	return psuo
}

// SetProwSuites sets the "prow_suites" edge to the Repository entity.
func (psuo *ProwSuitesUpdateOne) SetProwSuites(r *Repository) *ProwSuitesUpdateOne {
	return psuo.SetProwSuitesID(r.ID)
}

// Mutation returns the ProwSuitesMutation object of the builder.
func (psuo *ProwSuitesUpdateOne) Mutation() *ProwSuitesMutation {
	return psuo.mutation
}

// ClearProwSuites clears the "prow_suites" edge to the Repository entity.
func (psuo *ProwSuitesUpdateOne) ClearProwSuites() *ProwSuitesUpdateOne {
	psuo.mutation.ClearProwSuites()
	return psuo
}

// Where appends a list predicates to the ProwSuitesUpdate builder.
func (psuo *ProwSuitesUpdateOne) Where(ps ...predicate.ProwSuites) *ProwSuitesUpdateOne {
	psuo.mutation.Where(ps...)
	return psuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (psuo *ProwSuitesUpdateOne) Select(field string, fields ...string) *ProwSuitesUpdateOne {
	psuo.fields = append([]string{field}, fields...)
	return psuo
}

// Save executes the query and returns the updated ProwSuites entity.
func (psuo *ProwSuitesUpdateOne) Save(ctx context.Context) (*ProwSuites, error) {
	return withHooks(ctx, psuo.sqlSave, psuo.mutation, psuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (psuo *ProwSuitesUpdateOne) SaveX(ctx context.Context) *ProwSuites {
	node, err := psuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (psuo *ProwSuitesUpdateOne) Exec(ctx context.Context) error {
	_, err := psuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (psuo *ProwSuitesUpdateOne) ExecX(ctx context.Context) {
	if err := psuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (psuo *ProwSuitesUpdateOne) sqlSave(ctx context.Context) (_node *ProwSuites, err error) {
	_spec := sqlgraph.NewUpdateSpec(prowsuites.Table, prowsuites.Columns, sqlgraph.NewFieldSpec(prowsuites.FieldID, field.TypeInt))
	id, ok := psuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`db: missing "ProwSuites.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := psuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, prowsuites.FieldID)
		for _, f := range fields {
			if !prowsuites.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("db: invalid field %q for query", f)}
			}
			if f != prowsuites.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := psuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := psuo.mutation.JobID(); ok {
		_spec.SetField(prowsuites.FieldJobID, field.TypeString, value)
	}
	if value, ok := psuo.mutation.JobURL(); ok {
		_spec.SetField(prowsuites.FieldJobURL, field.TypeString, value)
	}
	if value, ok := psuo.mutation.JobName(); ok {
		_spec.SetField(prowsuites.FieldJobName, field.TypeString, value)
	}
	if value, ok := psuo.mutation.SuiteName(); ok {
		_spec.SetField(prowsuites.FieldSuiteName, field.TypeString, value)
	}
	if value, ok := psuo.mutation.Name(); ok {
		_spec.SetField(prowsuites.FieldName, field.TypeString, value)
	}
	if value, ok := psuo.mutation.Status(); ok {
		_spec.SetField(prowsuites.FieldStatus, field.TypeString, value)
	}
	if value, ok := psuo.mutation.ErrorMessage(); ok {
		_spec.SetField(prowsuites.FieldErrorMessage, field.TypeString, value)
	}
	if psuo.mutation.ErrorMessageCleared() {
		_spec.ClearField(prowsuites.FieldErrorMessage, field.TypeString)
	}
	if value, ok := psuo.mutation.ExternalServicesImpact(); ok {
		_spec.SetField(prowsuites.FieldExternalServicesImpact, field.TypeBool, value)
	}
	if psuo.mutation.ExternalServicesImpactCleared() {
		_spec.ClearField(prowsuites.FieldExternalServicesImpact, field.TypeBool)
	}
	if value, ok := psuo.mutation.Time(); ok {
		_spec.SetField(prowsuites.FieldTime, field.TypeFloat64, value)
	}
	if value, ok := psuo.mutation.AddedTime(); ok {
		_spec.AddField(prowsuites.FieldTime, field.TypeFloat64, value)
	}
	if value, ok := psuo.mutation.CreatedAt(); ok {
		_spec.SetField(prowsuites.FieldCreatedAt, field.TypeTime, value)
	}
	if psuo.mutation.CreatedAtCleared() {
		_spec.ClearField(prowsuites.FieldCreatedAt, field.TypeTime)
	}
	if psuo.mutation.ProwSuitesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   prowsuites.ProwSuitesTable,
			Columns: []string{prowsuites.ProwSuitesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(repository.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := psuo.mutation.ProwSuitesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   prowsuites.ProwSuitesTable,
			Columns: []string{prowsuites.ProwSuitesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(repository.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ProwSuites{config: psuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, psuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{prowsuites.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	psuo.mutation.done = true
	return _node, nil
}
