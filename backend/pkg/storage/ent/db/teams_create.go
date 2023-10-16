// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/bugs"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/failure"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/plugins"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/teams"
)

// TeamsCreate is the builder for creating a Teams entity.
type TeamsCreate struct {
	config
	mutation *TeamsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetTeamName sets the "team_name" field.
func (tc *TeamsCreate) SetTeamName(s string) *TeamsCreate {
	tc.mutation.SetTeamName(s)
	return tc
}

// SetDescription sets the "description" field.
func (tc *TeamsCreate) SetDescription(s string) *TeamsCreate {
	tc.mutation.SetDescription(s)
	return tc
}

// SetJiraKeys sets the "jira_keys" field.
func (tc *TeamsCreate) SetJiraKeys(s string) *TeamsCreate {
	tc.mutation.SetJiraKeys(s)
	return tc
}

// SetID sets the "id" field.
func (tc *TeamsCreate) SetID(u uuid.UUID) *TeamsCreate {
	tc.mutation.SetID(u)
	return tc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tc *TeamsCreate) SetNillableID(u *uuid.UUID) *TeamsCreate {
	if u != nil {
		tc.SetID(*u)
	}
	return tc
}

// AddRepositoryIDs adds the "repositories" edge to the Repository entity by IDs.
func (tc *TeamsCreate) AddRepositoryIDs(ids ...string) *TeamsCreate {
	tc.mutation.AddRepositoryIDs(ids...)
	return tc
}

// AddRepositories adds the "repositories" edges to the Repository entity.
func (tc *TeamsCreate) AddRepositories(r ...*Repository) *TeamsCreate {
	ids := make([]string, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return tc.AddRepositoryIDs(ids...)
}

// AddBugIDs adds the "bugs" edge to the Bugs entity by IDs.
func (tc *TeamsCreate) AddBugIDs(ids ...uuid.UUID) *TeamsCreate {
	tc.mutation.AddBugIDs(ids...)
	return tc
}

// AddBugs adds the "bugs" edges to the Bugs entity.
func (tc *TeamsCreate) AddBugs(b ...*Bugs) *TeamsCreate {
	ids := make([]uuid.UUID, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return tc.AddBugIDs(ids...)
}

// AddFailureIDs adds the "failures" edge to the Failure entity by IDs.
func (tc *TeamsCreate) AddFailureIDs(ids ...uuid.UUID) *TeamsCreate {
	tc.mutation.AddFailureIDs(ids...)
	return tc
}

// AddFailures adds the "failures" edges to the Failure entity.
func (tc *TeamsCreate) AddFailures(f ...*Failure) *TeamsCreate {
	ids := make([]uuid.UUID, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return tc.AddFailureIDs(ids...)
}

// AddPluginIDs adds the "plugins" edge to the Plugins entity by IDs.
func (tc *TeamsCreate) AddPluginIDs(ids ...uuid.UUID) *TeamsCreate {
	tc.mutation.AddPluginIDs(ids...)
	return tc
}

// AddPlugins adds the "plugins" edges to the Plugins entity.
func (tc *TeamsCreate) AddPlugins(p ...*Plugins) *TeamsCreate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return tc.AddPluginIDs(ids...)
}

// Mutation returns the TeamsMutation object of the builder.
func (tc *TeamsCreate) Mutation() *TeamsMutation {
	return tc.mutation
}

// Save creates the Teams in the database.
func (tc *TeamsCreate) Save(ctx context.Context) (*Teams, error) {
	tc.defaults()
	return withHooks[*Teams, TeamsMutation](ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TeamsCreate) SaveX(ctx context.Context) *Teams {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TeamsCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TeamsCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TeamsCreate) defaults() {
	if _, ok := tc.mutation.ID(); !ok {
		v := teams.DefaultID()
		tc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TeamsCreate) check() error {
	if _, ok := tc.mutation.TeamName(); !ok {
		return &ValidationError{Name: "team_name", err: errors.New(`db: missing required field "Teams.team_name"`)}
	}
	if _, ok := tc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`db: missing required field "Teams.description"`)}
	}
	if _, ok := tc.mutation.JiraKeys(); !ok {
		return &ValidationError{Name: "jira_keys", err: errors.New(`db: missing required field "Teams.jira_keys"`)}
	}
	return nil
}

func (tc *TeamsCreate) sqlSave(ctx context.Context) (*Teams, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TeamsCreate) createSpec() (*Teams, *sqlgraph.CreateSpec) {
	var (
		_node = &Teams{config: tc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: teams.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: teams.FieldID,
			},
		}
	)
	_spec.OnConflict = tc.conflict
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.TeamName(); ok {
		_spec.SetField(teams.FieldTeamName, field.TypeString, value)
		_node.TeamName = value
	}
	if value, ok := tc.mutation.Description(); ok {
		_spec.SetField(teams.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := tc.mutation.JiraKeys(); ok {
		_spec.SetField(teams.FieldJiraKeys, field.TypeString, value)
		_node.JiraKeys = value
	}
	if nodes := tc.mutation.RepositoriesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   teams.RepositoriesTable,
			Columns: []string{teams.RepositoriesColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.BugsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   teams.BugsTable,
			Columns: []string{teams.BugsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: bugs.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.FailuresIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   teams.FailuresTable,
			Columns: []string{teams.FailuresColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: failure.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.PluginsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   teams.PluginsTable,
			Columns: teams.PluginsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: plugins.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Teams.Create().
//		SetTeamName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TeamsUpsert) {
//			SetTeamName(v+v).
//		}).
//		Exec(ctx)
func (tc *TeamsCreate) OnConflict(opts ...sql.ConflictOption) *TeamsUpsertOne {
	tc.conflict = opts
	return &TeamsUpsertOne{
		create: tc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Teams.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tc *TeamsCreate) OnConflictColumns(columns ...string) *TeamsUpsertOne {
	tc.conflict = append(tc.conflict, sql.ConflictColumns(columns...))
	return &TeamsUpsertOne{
		create: tc,
	}
}

type (
	// TeamsUpsertOne is the builder for "upsert"-ing
	//  one Teams node.
	TeamsUpsertOne struct {
		create *TeamsCreate
	}

	// TeamsUpsert is the "OnConflict" setter.
	TeamsUpsert struct {
		*sql.UpdateSet
	}
)

// SetTeamName sets the "team_name" field.
func (u *TeamsUpsert) SetTeamName(v string) *TeamsUpsert {
	u.Set(teams.FieldTeamName, v)
	return u
}

// UpdateTeamName sets the "team_name" field to the value that was provided on create.
func (u *TeamsUpsert) UpdateTeamName() *TeamsUpsert {
	u.SetExcluded(teams.FieldTeamName)
	return u
}

// SetDescription sets the "description" field.
func (u *TeamsUpsert) SetDescription(v string) *TeamsUpsert {
	u.Set(teams.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *TeamsUpsert) UpdateDescription() *TeamsUpsert {
	u.SetExcluded(teams.FieldDescription)
	return u
}

// SetJiraKeys sets the "jira_keys" field.
func (u *TeamsUpsert) SetJiraKeys(v string) *TeamsUpsert {
	u.Set(teams.FieldJiraKeys, v)
	return u
}

// UpdateJiraKeys sets the "jira_keys" field to the value that was provided on create.
func (u *TeamsUpsert) UpdateJiraKeys() *TeamsUpsert {
	u.SetExcluded(teams.FieldJiraKeys)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Teams.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(teams.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TeamsUpsertOne) UpdateNewValues() *TeamsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(teams.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Teams.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *TeamsUpsertOne) Ignore() *TeamsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TeamsUpsertOne) DoNothing() *TeamsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TeamsCreate.OnConflict
// documentation for more info.
func (u *TeamsUpsertOne) Update(set func(*TeamsUpsert)) *TeamsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TeamsUpsert{UpdateSet: update})
	}))
	return u
}

// SetTeamName sets the "team_name" field.
func (u *TeamsUpsertOne) SetTeamName(v string) *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.SetTeamName(v)
	})
}

// UpdateTeamName sets the "team_name" field to the value that was provided on create.
func (u *TeamsUpsertOne) UpdateTeamName() *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateTeamName()
	})
}

// SetDescription sets the "description" field.
func (u *TeamsUpsertOne) SetDescription(v string) *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *TeamsUpsertOne) UpdateDescription() *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateDescription()
	})
}

// SetJiraKeys sets the "jira_keys" field.
func (u *TeamsUpsertOne) SetJiraKeys(v string) *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.SetJiraKeys(v)
	})
}

// UpdateJiraKeys sets the "jira_keys" field to the value that was provided on create.
func (u *TeamsUpsertOne) UpdateJiraKeys() *TeamsUpsertOne {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateJiraKeys()
	})
}

// Exec executes the query.
func (u *TeamsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for TeamsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TeamsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *TeamsUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("db: TeamsUpsertOne.ID is not supported by MySQL driver. Use TeamsUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *TeamsUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// TeamsCreateBulk is the builder for creating many Teams entities in bulk.
type TeamsCreateBulk struct {
	config
	builders []*TeamsCreate
	conflict []sql.ConflictOption
}

// Save creates the Teams entities in the database.
func (tcb *TeamsCreateBulk) Save(ctx context.Context) ([]*Teams, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Teams, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TeamsMutation)
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
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = tcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TeamsCreateBulk) SaveX(ctx context.Context) []*Teams {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TeamsCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TeamsCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Teams.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TeamsUpsert) {
//			SetTeamName(v+v).
//		}).
//		Exec(ctx)
func (tcb *TeamsCreateBulk) OnConflict(opts ...sql.ConflictOption) *TeamsUpsertBulk {
	tcb.conflict = opts
	return &TeamsUpsertBulk{
		create: tcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Teams.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tcb *TeamsCreateBulk) OnConflictColumns(columns ...string) *TeamsUpsertBulk {
	tcb.conflict = append(tcb.conflict, sql.ConflictColumns(columns...))
	return &TeamsUpsertBulk{
		create: tcb,
	}
}

// TeamsUpsertBulk is the builder for "upsert"-ing
// a bulk of Teams nodes.
type TeamsUpsertBulk struct {
	create *TeamsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Teams.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(teams.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TeamsUpsertBulk) UpdateNewValues() *TeamsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(teams.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Teams.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *TeamsUpsertBulk) Ignore() *TeamsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TeamsUpsertBulk) DoNothing() *TeamsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TeamsCreateBulk.OnConflict
// documentation for more info.
func (u *TeamsUpsertBulk) Update(set func(*TeamsUpsert)) *TeamsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TeamsUpsert{UpdateSet: update})
	}))
	return u
}

// SetTeamName sets the "team_name" field.
func (u *TeamsUpsertBulk) SetTeamName(v string) *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.SetTeamName(v)
	})
}

// UpdateTeamName sets the "team_name" field to the value that was provided on create.
func (u *TeamsUpsertBulk) UpdateTeamName() *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateTeamName()
	})
}

// SetDescription sets the "description" field.
func (u *TeamsUpsertBulk) SetDescription(v string) *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *TeamsUpsertBulk) UpdateDescription() *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateDescription()
	})
}

// SetJiraKeys sets the "jira_keys" field.
func (u *TeamsUpsertBulk) SetJiraKeys(v string) *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.SetJiraKeys(v)
	})
}

// UpdateJiraKeys sets the "jira_keys" field to the value that was provided on create.
func (u *TeamsUpsertBulk) UpdateJiraKeys() *TeamsUpsertBulk {
	return u.Update(func(s *TeamsUpsert) {
		s.UpdateJiraKeys()
	})
}

// Exec executes the query.
func (u *TeamsUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("db: OnConflict was set for builder %d. Set it on the TeamsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for TeamsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TeamsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
