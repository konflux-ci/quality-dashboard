// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/repository"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/workflows"
)

// WorkflowsCreate is the builder for creating a Workflows entity.
type WorkflowsCreate struct {
	config
	mutation *WorkflowsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetWorkflowID sets the "workflow_id" field.
func (wc *WorkflowsCreate) SetWorkflowID(u uuid.UUID) *WorkflowsCreate {
	wc.mutation.SetWorkflowID(u)
	return wc
}

// SetNillableWorkflowID sets the "workflow_id" field if the given value is not nil.
func (wc *WorkflowsCreate) SetNillableWorkflowID(u *uuid.UUID) *WorkflowsCreate {
	if u != nil {
		wc.SetWorkflowID(*u)
	}
	return wc
}

// SetWorkflowName sets the "workflow_name" field.
func (wc *WorkflowsCreate) SetWorkflowName(s string) *WorkflowsCreate {
	wc.mutation.SetWorkflowName(s)
	return wc
}

// SetBadgeURL sets the "badge_url" field.
func (wc *WorkflowsCreate) SetBadgeURL(s string) *WorkflowsCreate {
	wc.mutation.SetBadgeURL(s)
	return wc
}

// SetHTMLURL sets the "html_url" field.
func (wc *WorkflowsCreate) SetHTMLURL(s string) *WorkflowsCreate {
	wc.mutation.SetHTMLURL(s)
	return wc
}

// SetJobURL sets the "job_url" field.
func (wc *WorkflowsCreate) SetJobURL(s string) *WorkflowsCreate {
	wc.mutation.SetJobURL(s)
	return wc
}

// SetState sets the "state" field.
func (wc *WorkflowsCreate) SetState(s string) *WorkflowsCreate {
	wc.mutation.SetState(s)
	return wc
}

// SetWorkflowsID sets the "workflows" edge to the Repository entity by ID.
func (wc *WorkflowsCreate) SetWorkflowsID(id string) *WorkflowsCreate {
	wc.mutation.SetWorkflowsID(id)
	return wc
}

// SetNillableWorkflowsID sets the "workflows" edge to the Repository entity by ID if the given value is not nil.
func (wc *WorkflowsCreate) SetNillableWorkflowsID(id *string) *WorkflowsCreate {
	if id != nil {
		wc = wc.SetWorkflowsID(*id)
	}
	return wc
}

// SetWorkflows sets the "workflows" edge to the Repository entity.
func (wc *WorkflowsCreate) SetWorkflows(r *Repository) *WorkflowsCreate {
	return wc.SetWorkflowsID(r.ID)
}

// Mutation returns the WorkflowsMutation object of the builder.
func (wc *WorkflowsCreate) Mutation() *WorkflowsMutation {
	return wc.mutation
}

// Save creates the Workflows in the database.
func (wc *WorkflowsCreate) Save(ctx context.Context) (*Workflows, error) {
	wc.defaults()
	return withHooks[*Workflows, WorkflowsMutation](ctx, wc.sqlSave, wc.mutation, wc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (wc *WorkflowsCreate) SaveX(ctx context.Context) *Workflows {
	v, err := wc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wc *WorkflowsCreate) Exec(ctx context.Context) error {
	_, err := wc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wc *WorkflowsCreate) ExecX(ctx context.Context) {
	if err := wc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wc *WorkflowsCreate) defaults() {
	if _, ok := wc.mutation.WorkflowID(); !ok {
		v := workflows.DefaultWorkflowID()
		wc.mutation.SetWorkflowID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wc *WorkflowsCreate) check() error {
	if _, ok := wc.mutation.WorkflowID(); !ok {
		return &ValidationError{Name: "workflow_id", err: errors.New(`db: missing required field "Workflows.workflow_id"`)}
	}
	if _, ok := wc.mutation.WorkflowName(); !ok {
		return &ValidationError{Name: "workflow_name", err: errors.New(`db: missing required field "Workflows.workflow_name"`)}
	}
	if v, ok := wc.mutation.WorkflowName(); ok {
		if err := workflows.WorkflowNameValidator(v); err != nil {
			return &ValidationError{Name: "workflow_name", err: fmt.Errorf(`db: validator failed for field "Workflows.workflow_name": %w`, err)}
		}
	}
	if _, ok := wc.mutation.BadgeURL(); !ok {
		return &ValidationError{Name: "badge_url", err: errors.New(`db: missing required field "Workflows.badge_url"`)}
	}
	if _, ok := wc.mutation.HTMLURL(); !ok {
		return &ValidationError{Name: "html_url", err: errors.New(`db: missing required field "Workflows.html_url"`)}
	}
	if _, ok := wc.mutation.JobURL(); !ok {
		return &ValidationError{Name: "job_url", err: errors.New(`db: missing required field "Workflows.job_url"`)}
	}
	if _, ok := wc.mutation.State(); !ok {
		return &ValidationError{Name: "state", err: errors.New(`db: missing required field "Workflows.state"`)}
	}
	return nil
}

func (wc *WorkflowsCreate) sqlSave(ctx context.Context) (*Workflows, error) {
	if err := wc.check(); err != nil {
		return nil, err
	}
	_node, _spec := wc.createSpec()
	if err := sqlgraph.CreateNode(ctx, wc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	wc.mutation.id = &_node.ID
	wc.mutation.done = true
	return _node, nil
}

func (wc *WorkflowsCreate) createSpec() (*Workflows, *sqlgraph.CreateSpec) {
	var (
		_node = &Workflows{config: wc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: workflows.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workflows.FieldID,
			},
		}
	)
	_spec.OnConflict = wc.conflict
	if value, ok := wc.mutation.WorkflowID(); ok {
		_spec.SetField(workflows.FieldWorkflowID, field.TypeUUID, value)
		_node.WorkflowID = value
	}
	if value, ok := wc.mutation.WorkflowName(); ok {
		_spec.SetField(workflows.FieldWorkflowName, field.TypeString, value)
		_node.WorkflowName = value
	}
	if value, ok := wc.mutation.BadgeURL(); ok {
		_spec.SetField(workflows.FieldBadgeURL, field.TypeString, value)
		_node.BadgeURL = value
	}
	if value, ok := wc.mutation.HTMLURL(); ok {
		_spec.SetField(workflows.FieldHTMLURL, field.TypeString, value)
		_node.HTMLURL = value
	}
	if value, ok := wc.mutation.JobURL(); ok {
		_spec.SetField(workflows.FieldJobURL, field.TypeString, value)
		_node.JobURL = value
	}
	if value, ok := wc.mutation.State(); ok {
		_spec.SetField(workflows.FieldState, field.TypeString, value)
		_node.State = value
	}
	if nodes := wc.mutation.WorkflowsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   workflows.WorkflowsTable,
			Columns: []string{workflows.WorkflowsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: repository.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.repository_workflows = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Workflows.Create().
//		SetWorkflowID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.WorkflowsUpsert) {
//			SetWorkflowID(v+v).
//		}).
//		Exec(ctx)
func (wc *WorkflowsCreate) OnConflict(opts ...sql.ConflictOption) *WorkflowsUpsertOne {
	wc.conflict = opts
	return &WorkflowsUpsertOne{
		create: wc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Workflows.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (wc *WorkflowsCreate) OnConflictColumns(columns ...string) *WorkflowsUpsertOne {
	wc.conflict = append(wc.conflict, sql.ConflictColumns(columns...))
	return &WorkflowsUpsertOne{
		create: wc,
	}
}

type (
	// WorkflowsUpsertOne is the builder for "upsert"-ing
	//  one Workflows node.
	WorkflowsUpsertOne struct {
		create *WorkflowsCreate
	}

	// WorkflowsUpsert is the "OnConflict" setter.
	WorkflowsUpsert struct {
		*sql.UpdateSet
	}
)

// SetWorkflowID sets the "workflow_id" field.
func (u *WorkflowsUpsert) SetWorkflowID(v uuid.UUID) *WorkflowsUpsert {
	u.Set(workflows.FieldWorkflowID, v)
	return u
}

// UpdateWorkflowID sets the "workflow_id" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateWorkflowID() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldWorkflowID)
	return u
}

// SetWorkflowName sets the "workflow_name" field.
func (u *WorkflowsUpsert) SetWorkflowName(v string) *WorkflowsUpsert {
	u.Set(workflows.FieldWorkflowName, v)
	return u
}

// UpdateWorkflowName sets the "workflow_name" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateWorkflowName() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldWorkflowName)
	return u
}

// SetBadgeURL sets the "badge_url" field.
func (u *WorkflowsUpsert) SetBadgeURL(v string) *WorkflowsUpsert {
	u.Set(workflows.FieldBadgeURL, v)
	return u
}

// UpdateBadgeURL sets the "badge_url" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateBadgeURL() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldBadgeURL)
	return u
}

// SetHTMLURL sets the "html_url" field.
func (u *WorkflowsUpsert) SetHTMLURL(v string) *WorkflowsUpsert {
	u.Set(workflows.FieldHTMLURL, v)
	return u
}

// UpdateHTMLURL sets the "html_url" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateHTMLURL() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldHTMLURL)
	return u
}

// SetJobURL sets the "job_url" field.
func (u *WorkflowsUpsert) SetJobURL(v string) *WorkflowsUpsert {
	u.Set(workflows.FieldJobURL, v)
	return u
}

// UpdateJobURL sets the "job_url" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateJobURL() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldJobURL)
	return u
}

// SetState sets the "state" field.
func (u *WorkflowsUpsert) SetState(v string) *WorkflowsUpsert {
	u.Set(workflows.FieldState, v)
	return u
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *WorkflowsUpsert) UpdateState() *WorkflowsUpsert {
	u.SetExcluded(workflows.FieldState)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Workflows.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *WorkflowsUpsertOne) UpdateNewValues() *WorkflowsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Workflows.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *WorkflowsUpsertOne) Ignore() *WorkflowsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *WorkflowsUpsertOne) DoNothing() *WorkflowsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the WorkflowsCreate.OnConflict
// documentation for more info.
func (u *WorkflowsUpsertOne) Update(set func(*WorkflowsUpsert)) *WorkflowsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&WorkflowsUpsert{UpdateSet: update})
	}))
	return u
}

// SetWorkflowID sets the "workflow_id" field.
func (u *WorkflowsUpsertOne) SetWorkflowID(v uuid.UUID) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetWorkflowID(v)
	})
}

// UpdateWorkflowID sets the "workflow_id" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateWorkflowID() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateWorkflowID()
	})
}

// SetWorkflowName sets the "workflow_name" field.
func (u *WorkflowsUpsertOne) SetWorkflowName(v string) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetWorkflowName(v)
	})
}

// UpdateWorkflowName sets the "workflow_name" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateWorkflowName() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateWorkflowName()
	})
}

// SetBadgeURL sets the "badge_url" field.
func (u *WorkflowsUpsertOne) SetBadgeURL(v string) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetBadgeURL(v)
	})
}

// UpdateBadgeURL sets the "badge_url" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateBadgeURL() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateBadgeURL()
	})
}

// SetHTMLURL sets the "html_url" field.
func (u *WorkflowsUpsertOne) SetHTMLURL(v string) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetHTMLURL(v)
	})
}

// UpdateHTMLURL sets the "html_url" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateHTMLURL() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateHTMLURL()
	})
}

// SetJobURL sets the "job_url" field.
func (u *WorkflowsUpsertOne) SetJobURL(v string) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetJobURL(v)
	})
}

// UpdateJobURL sets the "job_url" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateJobURL() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateJobURL()
	})
}

// SetState sets the "state" field.
func (u *WorkflowsUpsertOne) SetState(v string) *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *WorkflowsUpsertOne) UpdateState() *WorkflowsUpsertOne {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateState()
	})
}

// Exec executes the query.
func (u *WorkflowsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for WorkflowsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *WorkflowsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *WorkflowsUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *WorkflowsUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// WorkflowsCreateBulk is the builder for creating many Workflows entities in bulk.
type WorkflowsCreateBulk struct {
	config
	builders []*WorkflowsCreate
	conflict []sql.ConflictOption
}

// Save creates the Workflows entities in the database.
func (wcb *WorkflowsCreateBulk) Save(ctx context.Context) ([]*Workflows, error) {
	specs := make([]*sqlgraph.CreateSpec, len(wcb.builders))
	nodes := make([]*Workflows, len(wcb.builders))
	mutators := make([]Mutator, len(wcb.builders))
	for i := range wcb.builders {
		func(i int, root context.Context) {
			builder := wcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*WorkflowsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, wcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = wcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, wcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, wcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (wcb *WorkflowsCreateBulk) SaveX(ctx context.Context) []*Workflows {
	v, err := wcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wcb *WorkflowsCreateBulk) Exec(ctx context.Context) error {
	_, err := wcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wcb *WorkflowsCreateBulk) ExecX(ctx context.Context) {
	if err := wcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Workflows.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.WorkflowsUpsert) {
//			SetWorkflowID(v+v).
//		}).
//		Exec(ctx)
func (wcb *WorkflowsCreateBulk) OnConflict(opts ...sql.ConflictOption) *WorkflowsUpsertBulk {
	wcb.conflict = opts
	return &WorkflowsUpsertBulk{
		create: wcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Workflows.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (wcb *WorkflowsCreateBulk) OnConflictColumns(columns ...string) *WorkflowsUpsertBulk {
	wcb.conflict = append(wcb.conflict, sql.ConflictColumns(columns...))
	return &WorkflowsUpsertBulk{
		create: wcb,
	}
}

// WorkflowsUpsertBulk is the builder for "upsert"-ing
// a bulk of Workflows nodes.
type WorkflowsUpsertBulk struct {
	create *WorkflowsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Workflows.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *WorkflowsUpsertBulk) UpdateNewValues() *WorkflowsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Workflows.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *WorkflowsUpsertBulk) Ignore() *WorkflowsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *WorkflowsUpsertBulk) DoNothing() *WorkflowsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the WorkflowsCreateBulk.OnConflict
// documentation for more info.
func (u *WorkflowsUpsertBulk) Update(set func(*WorkflowsUpsert)) *WorkflowsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&WorkflowsUpsert{UpdateSet: update})
	}))
	return u
}

// SetWorkflowID sets the "workflow_id" field.
func (u *WorkflowsUpsertBulk) SetWorkflowID(v uuid.UUID) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetWorkflowID(v)
	})
}

// UpdateWorkflowID sets the "workflow_id" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateWorkflowID() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateWorkflowID()
	})
}

// SetWorkflowName sets the "workflow_name" field.
func (u *WorkflowsUpsertBulk) SetWorkflowName(v string) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetWorkflowName(v)
	})
}

// UpdateWorkflowName sets the "workflow_name" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateWorkflowName() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateWorkflowName()
	})
}

// SetBadgeURL sets the "badge_url" field.
func (u *WorkflowsUpsertBulk) SetBadgeURL(v string) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetBadgeURL(v)
	})
}

// UpdateBadgeURL sets the "badge_url" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateBadgeURL() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateBadgeURL()
	})
}

// SetHTMLURL sets the "html_url" field.
func (u *WorkflowsUpsertBulk) SetHTMLURL(v string) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetHTMLURL(v)
	})
}

// UpdateHTMLURL sets the "html_url" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateHTMLURL() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateHTMLURL()
	})
}

// SetJobURL sets the "job_url" field.
func (u *WorkflowsUpsertBulk) SetJobURL(v string) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetJobURL(v)
	})
}

// UpdateJobURL sets the "job_url" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateJobURL() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateJobURL()
	})
}

// SetState sets the "state" field.
func (u *WorkflowsUpsertBulk) SetState(v string) *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *WorkflowsUpsertBulk) UpdateState() *WorkflowsUpsertBulk {
	return u.Update(func(s *WorkflowsUpsert) {
		s.UpdateState()
	})
}

// Exec executes the query.
func (u *WorkflowsUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("db: OnConflict was set for builder %d. Set it on the WorkflowsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for WorkflowsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *WorkflowsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
