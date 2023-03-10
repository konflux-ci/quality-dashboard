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
	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/pullrequests"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

// PullRequestsCreate is the builder for creating a PullRequests entity.
type PullRequestsCreate struct {
	config
	mutation *PullRequestsMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetPrID sets the "pr_id" field.
func (prc *PullRequestsCreate) SetPrID(u uuid.UUID) *PullRequestsCreate {
	prc.mutation.SetPrID(u)
	return prc
}

// SetNillablePrID sets the "pr_id" field if the given value is not nil.
func (prc *PullRequestsCreate) SetNillablePrID(u *uuid.UUID) *PullRequestsCreate {
	if u != nil {
		prc.SetPrID(*u)
	}
	return prc
}

// SetRepositoryName sets the "repository_name" field.
func (prc *PullRequestsCreate) SetRepositoryName(s string) *PullRequestsCreate {
	prc.mutation.SetRepositoryName(s)
	return prc
}

// SetRepositoryOrganization sets the "repository_organization" field.
func (prc *PullRequestsCreate) SetRepositoryOrganization(s string) *PullRequestsCreate {
	prc.mutation.SetRepositoryOrganization(s)
	return prc
}

// SetNumber sets the "number" field.
func (prc *PullRequestsCreate) SetNumber(i int) *PullRequestsCreate {
	prc.mutation.SetNumber(i)
	return prc
}

// SetCreatedAt sets the "created_at" field.
func (prc *PullRequestsCreate) SetCreatedAt(t time.Time) *PullRequestsCreate {
	prc.mutation.SetCreatedAt(t)
	return prc
}

// SetClosedAt sets the "closed_at" field.
func (prc *PullRequestsCreate) SetClosedAt(t time.Time) *PullRequestsCreate {
	prc.mutation.SetClosedAt(t)
	return prc
}

// SetMergedAt sets the "merged_at" field.
func (prc *PullRequestsCreate) SetMergedAt(t time.Time) *PullRequestsCreate {
	prc.mutation.SetMergedAt(t)
	return prc
}

// SetState sets the "state" field.
func (prc *PullRequestsCreate) SetState(s string) *PullRequestsCreate {
	prc.mutation.SetState(s)
	return prc
}

// SetAuthor sets the "author" field.
func (prc *PullRequestsCreate) SetAuthor(s string) *PullRequestsCreate {
	prc.mutation.SetAuthor(s)
	return prc
}

// SetTitle sets the "title" field.
func (prc *PullRequestsCreate) SetTitle(s string) *PullRequestsCreate {
	prc.mutation.SetTitle(s)
	return prc
}

// SetPrsID sets the "prs" edge to the Repository entity by ID.
func (prc *PullRequestsCreate) SetPrsID(id string) *PullRequestsCreate {
	prc.mutation.SetPrsID(id)
	return prc
}

// SetNillablePrsID sets the "prs" edge to the Repository entity by ID if the given value is not nil.
func (prc *PullRequestsCreate) SetNillablePrsID(id *string) *PullRequestsCreate {
	if id != nil {
		prc = prc.SetPrsID(*id)
	}
	return prc
}

// SetPrs sets the "prs" edge to the Repository entity.
func (prc *PullRequestsCreate) SetPrs(r *Repository) *PullRequestsCreate {
	return prc.SetPrsID(r.ID)
}

// Mutation returns the PullRequestsMutation object of the builder.
func (prc *PullRequestsCreate) Mutation() *PullRequestsMutation {
	return prc.mutation
}

// Save creates the PullRequests in the database.
func (prc *PullRequestsCreate) Save(ctx context.Context) (*PullRequests, error) {
	prc.defaults()
	return withHooks[*PullRequests, PullRequestsMutation](ctx, prc.sqlSave, prc.mutation, prc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (prc *PullRequestsCreate) SaveX(ctx context.Context) *PullRequests {
	v, err := prc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (prc *PullRequestsCreate) Exec(ctx context.Context) error {
	_, err := prc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (prc *PullRequestsCreate) ExecX(ctx context.Context) {
	if err := prc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (prc *PullRequestsCreate) defaults() {
	if _, ok := prc.mutation.PrID(); !ok {
		v := pullrequests.DefaultPrID()
		prc.mutation.SetPrID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (prc *PullRequestsCreate) check() error {
	if _, ok := prc.mutation.PrID(); !ok {
		return &ValidationError{Name: "pr_id", err: errors.New(`db: missing required field "PullRequests.pr_id"`)}
	}
	if _, ok := prc.mutation.RepositoryName(); !ok {
		return &ValidationError{Name: "repository_name", err: errors.New(`db: missing required field "PullRequests.repository_name"`)}
	}
	if _, ok := prc.mutation.RepositoryOrganization(); !ok {
		return &ValidationError{Name: "repository_organization", err: errors.New(`db: missing required field "PullRequests.repository_organization"`)}
	}
	if _, ok := prc.mutation.Number(); !ok {
		return &ValidationError{Name: "number", err: errors.New(`db: missing required field "PullRequests.number"`)}
	}
	if _, ok := prc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`db: missing required field "PullRequests.created_at"`)}
	}
	if _, ok := prc.mutation.ClosedAt(); !ok {
		return &ValidationError{Name: "closed_at", err: errors.New(`db: missing required field "PullRequests.closed_at"`)}
	}
	if _, ok := prc.mutation.MergedAt(); !ok {
		return &ValidationError{Name: "merged_at", err: errors.New(`db: missing required field "PullRequests.merged_at"`)}
	}
	if _, ok := prc.mutation.State(); !ok {
		return &ValidationError{Name: "state", err: errors.New(`db: missing required field "PullRequests.state"`)}
	}
	if _, ok := prc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New(`db: missing required field "PullRequests.author"`)}
	}
	if _, ok := prc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`db: missing required field "PullRequests.title"`)}
	}
	return nil
}

func (prc *PullRequestsCreate) sqlSave(ctx context.Context) (*PullRequests, error) {
	if err := prc.check(); err != nil {
		return nil, err
	}
	_node, _spec := prc.createSpec()
	if err := sqlgraph.CreateNode(ctx, prc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	prc.mutation.id = &_node.ID
	prc.mutation.done = true
	return _node, nil
}

func (prc *PullRequestsCreate) createSpec() (*PullRequests, *sqlgraph.CreateSpec) {
	var (
		_node = &PullRequests{config: prc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: pullrequests.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: pullrequests.FieldID,
			},
		}
	)
	_spec.OnConflict = prc.conflict
	if value, ok := prc.mutation.PrID(); ok {
		_spec.SetField(pullrequests.FieldPrID, field.TypeUUID, value)
		_node.PrID = value
	}
	if value, ok := prc.mutation.RepositoryName(); ok {
		_spec.SetField(pullrequests.FieldRepositoryName, field.TypeString, value)
		_node.RepositoryName = value
	}
	if value, ok := prc.mutation.RepositoryOrganization(); ok {
		_spec.SetField(pullrequests.FieldRepositoryOrganization, field.TypeString, value)
		_node.RepositoryOrganization = value
	}
	if value, ok := prc.mutation.Number(); ok {
		_spec.SetField(pullrequests.FieldNumber, field.TypeInt, value)
		_node.Number = value
	}
	if value, ok := prc.mutation.CreatedAt(); ok {
		_spec.SetField(pullrequests.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := prc.mutation.ClosedAt(); ok {
		_spec.SetField(pullrequests.FieldClosedAt, field.TypeTime, value)
		_node.ClosedAt = value
	}
	if value, ok := prc.mutation.MergedAt(); ok {
		_spec.SetField(pullrequests.FieldMergedAt, field.TypeTime, value)
		_node.MergedAt = value
	}
	if value, ok := prc.mutation.State(); ok {
		_spec.SetField(pullrequests.FieldState, field.TypeString, value)
		_node.State = value
	}
	if value, ok := prc.mutation.Author(); ok {
		_spec.SetField(pullrequests.FieldAuthor, field.TypeString, value)
		_node.Author = value
	}
	if value, ok := prc.mutation.Title(); ok {
		_spec.SetField(pullrequests.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if nodes := prc.mutation.PrsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   pullrequests.PrsTable,
			Columns: []string{pullrequests.PrsColumn},
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
		_node.repository_prs = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PullRequests.Create().
//		SetPrID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PullRequestsUpsert) {
//			SetPrID(v+v).
//		}).
//		Exec(ctx)
func (prc *PullRequestsCreate) OnConflict(opts ...sql.ConflictOption) *PullRequestsUpsertOne {
	prc.conflict = opts
	return &PullRequestsUpsertOne{
		create: prc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (prc *PullRequestsCreate) OnConflictColumns(columns ...string) *PullRequestsUpsertOne {
	prc.conflict = append(prc.conflict, sql.ConflictColumns(columns...))
	return &PullRequestsUpsertOne{
		create: prc,
	}
}

type (
	// PullRequestsUpsertOne is the builder for "upsert"-ing
	//  one PullRequests node.
	PullRequestsUpsertOne struct {
		create *PullRequestsCreate
	}

	// PullRequestsUpsert is the "OnConflict" setter.
	PullRequestsUpsert struct {
		*sql.UpdateSet
	}
)

// SetRepositoryName sets the "repository_name" field.
func (u *PullRequestsUpsert) SetRepositoryName(v string) *PullRequestsUpsert {
	u.Set(pullrequests.FieldRepositoryName, v)
	return u
}

// UpdateRepositoryName sets the "repository_name" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateRepositoryName() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldRepositoryName)
	return u
}

// SetRepositoryOrganization sets the "repository_organization" field.
func (u *PullRequestsUpsert) SetRepositoryOrganization(v string) *PullRequestsUpsert {
	u.Set(pullrequests.FieldRepositoryOrganization, v)
	return u
}

// UpdateRepositoryOrganization sets the "repository_organization" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateRepositoryOrganization() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldRepositoryOrganization)
	return u
}

// SetNumber sets the "number" field.
func (u *PullRequestsUpsert) SetNumber(v int) *PullRequestsUpsert {
	u.Set(pullrequests.FieldNumber, v)
	return u
}

// UpdateNumber sets the "number" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateNumber() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldNumber)
	return u
}

// AddNumber adds v to the "number" field.
func (u *PullRequestsUpsert) AddNumber(v int) *PullRequestsUpsert {
	u.Add(pullrequests.FieldNumber, v)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *PullRequestsUpsert) SetCreatedAt(v time.Time) *PullRequestsUpsert {
	u.Set(pullrequests.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateCreatedAt() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldCreatedAt)
	return u
}

// SetClosedAt sets the "closed_at" field.
func (u *PullRequestsUpsert) SetClosedAt(v time.Time) *PullRequestsUpsert {
	u.Set(pullrequests.FieldClosedAt, v)
	return u
}

// UpdateClosedAt sets the "closed_at" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateClosedAt() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldClosedAt)
	return u
}

// SetMergedAt sets the "merged_at" field.
func (u *PullRequestsUpsert) SetMergedAt(v time.Time) *PullRequestsUpsert {
	u.Set(pullrequests.FieldMergedAt, v)
	return u
}

// UpdateMergedAt sets the "merged_at" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateMergedAt() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldMergedAt)
	return u
}

// SetState sets the "state" field.
func (u *PullRequestsUpsert) SetState(v string) *PullRequestsUpsert {
	u.Set(pullrequests.FieldState, v)
	return u
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateState() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldState)
	return u
}

// SetAuthor sets the "author" field.
func (u *PullRequestsUpsert) SetAuthor(v string) *PullRequestsUpsert {
	u.Set(pullrequests.FieldAuthor, v)
	return u
}

// UpdateAuthor sets the "author" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateAuthor() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldAuthor)
	return u
}

// SetTitle sets the "title" field.
func (u *PullRequestsUpsert) SetTitle(v string) *PullRequestsUpsert {
	u.Set(pullrequests.FieldTitle, v)
	return u
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *PullRequestsUpsert) UpdateTitle() *PullRequestsUpsert {
	u.SetExcluded(pullrequests.FieldTitle)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PullRequestsUpsertOne) UpdateNewValues() *PullRequestsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.PrID(); exists {
			s.SetIgnore(pullrequests.FieldPrID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PullRequestsUpsertOne) Ignore() *PullRequestsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PullRequestsUpsertOne) DoNothing() *PullRequestsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PullRequestsCreate.OnConflict
// documentation for more info.
func (u *PullRequestsUpsertOne) Update(set func(*PullRequestsUpsert)) *PullRequestsUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PullRequestsUpsert{UpdateSet: update})
	}))
	return u
}

// SetRepositoryName sets the "repository_name" field.
func (u *PullRequestsUpsertOne) SetRepositoryName(v string) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetRepositoryName(v)
	})
}

// UpdateRepositoryName sets the "repository_name" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateRepositoryName() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateRepositoryName()
	})
}

// SetRepositoryOrganization sets the "repository_organization" field.
func (u *PullRequestsUpsertOne) SetRepositoryOrganization(v string) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetRepositoryOrganization(v)
	})
}

// UpdateRepositoryOrganization sets the "repository_organization" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateRepositoryOrganization() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateRepositoryOrganization()
	})
}

// SetNumber sets the "number" field.
func (u *PullRequestsUpsertOne) SetNumber(v int) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetNumber(v)
	})
}

// AddNumber adds v to the "number" field.
func (u *PullRequestsUpsertOne) AddNumber(v int) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.AddNumber(v)
	})
}

// UpdateNumber sets the "number" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateNumber() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateNumber()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *PullRequestsUpsertOne) SetCreatedAt(v time.Time) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateCreatedAt() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetClosedAt sets the "closed_at" field.
func (u *PullRequestsUpsertOne) SetClosedAt(v time.Time) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetClosedAt(v)
	})
}

// UpdateClosedAt sets the "closed_at" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateClosedAt() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateClosedAt()
	})
}

// SetMergedAt sets the "merged_at" field.
func (u *PullRequestsUpsertOne) SetMergedAt(v time.Time) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetMergedAt(v)
	})
}

// UpdateMergedAt sets the "merged_at" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateMergedAt() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateMergedAt()
	})
}

// SetState sets the "state" field.
func (u *PullRequestsUpsertOne) SetState(v string) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateState() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateState()
	})
}

// SetAuthor sets the "author" field.
func (u *PullRequestsUpsertOne) SetAuthor(v string) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetAuthor(v)
	})
}

// UpdateAuthor sets the "author" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateAuthor() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateAuthor()
	})
}

// SetTitle sets the "title" field.
func (u *PullRequestsUpsertOne) SetTitle(v string) *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *PullRequestsUpsertOne) UpdateTitle() *PullRequestsUpsertOne {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateTitle()
	})
}

// Exec executes the query.
func (u *PullRequestsUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for PullRequestsCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PullRequestsUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PullRequestsUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PullRequestsUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PullRequestsCreateBulk is the builder for creating many PullRequests entities in bulk.
type PullRequestsCreateBulk struct {
	config
	builders []*PullRequestsCreate
	conflict []sql.ConflictOption
}

// Save creates the PullRequests entities in the database.
func (prcb *PullRequestsCreateBulk) Save(ctx context.Context) ([]*PullRequests, error) {
	specs := make([]*sqlgraph.CreateSpec, len(prcb.builders))
	nodes := make([]*PullRequests, len(prcb.builders))
	mutators := make([]Mutator, len(prcb.builders))
	for i := range prcb.builders {
		func(i int, root context.Context) {
			builder := prcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PullRequestsMutation)
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
					_, err = mutators[i+1].Mutate(root, prcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = prcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, prcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, prcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (prcb *PullRequestsCreateBulk) SaveX(ctx context.Context) []*PullRequests {
	v, err := prcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (prcb *PullRequestsCreateBulk) Exec(ctx context.Context) error {
	_, err := prcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (prcb *PullRequestsCreateBulk) ExecX(ctx context.Context) {
	if err := prcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PullRequests.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PullRequestsUpsert) {
//			SetPrID(v+v).
//		}).
//		Exec(ctx)
func (prcb *PullRequestsCreateBulk) OnConflict(opts ...sql.ConflictOption) *PullRequestsUpsertBulk {
	prcb.conflict = opts
	return &PullRequestsUpsertBulk{
		create: prcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (prcb *PullRequestsCreateBulk) OnConflictColumns(columns ...string) *PullRequestsUpsertBulk {
	prcb.conflict = append(prcb.conflict, sql.ConflictColumns(columns...))
	return &PullRequestsUpsertBulk{
		create: prcb,
	}
}

// PullRequestsUpsertBulk is the builder for "upsert"-ing
// a bulk of PullRequests nodes.
type PullRequestsUpsertBulk struct {
	create *PullRequestsCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *PullRequestsUpsertBulk) UpdateNewValues() *PullRequestsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.PrID(); exists {
				s.SetIgnore(pullrequests.FieldPrID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PullRequests.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PullRequestsUpsertBulk) Ignore() *PullRequestsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PullRequestsUpsertBulk) DoNothing() *PullRequestsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PullRequestsCreateBulk.OnConflict
// documentation for more info.
func (u *PullRequestsUpsertBulk) Update(set func(*PullRequestsUpsert)) *PullRequestsUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PullRequestsUpsert{UpdateSet: update})
	}))
	return u
}

// SetRepositoryName sets the "repository_name" field.
func (u *PullRequestsUpsertBulk) SetRepositoryName(v string) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetRepositoryName(v)
	})
}

// UpdateRepositoryName sets the "repository_name" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateRepositoryName() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateRepositoryName()
	})
}

// SetRepositoryOrganization sets the "repository_organization" field.
func (u *PullRequestsUpsertBulk) SetRepositoryOrganization(v string) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetRepositoryOrganization(v)
	})
}

// UpdateRepositoryOrganization sets the "repository_organization" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateRepositoryOrganization() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateRepositoryOrganization()
	})
}

// SetNumber sets the "number" field.
func (u *PullRequestsUpsertBulk) SetNumber(v int) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetNumber(v)
	})
}

// AddNumber adds v to the "number" field.
func (u *PullRequestsUpsertBulk) AddNumber(v int) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.AddNumber(v)
	})
}

// UpdateNumber sets the "number" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateNumber() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateNumber()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *PullRequestsUpsertBulk) SetCreatedAt(v time.Time) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateCreatedAt() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetClosedAt sets the "closed_at" field.
func (u *PullRequestsUpsertBulk) SetClosedAt(v time.Time) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetClosedAt(v)
	})
}

// UpdateClosedAt sets the "closed_at" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateClosedAt() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateClosedAt()
	})
}

// SetMergedAt sets the "merged_at" field.
func (u *PullRequestsUpsertBulk) SetMergedAt(v time.Time) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetMergedAt(v)
	})
}

// UpdateMergedAt sets the "merged_at" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateMergedAt() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateMergedAt()
	})
}

// SetState sets the "state" field.
func (u *PullRequestsUpsertBulk) SetState(v string) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateState() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateState()
	})
}

// SetAuthor sets the "author" field.
func (u *PullRequestsUpsertBulk) SetAuthor(v string) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetAuthor(v)
	})
}

// UpdateAuthor sets the "author" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateAuthor() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateAuthor()
	})
}

// SetTitle sets the "title" field.
func (u *PullRequestsUpsertBulk) SetTitle(v string) *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.SetTitle(v)
	})
}

// UpdateTitle sets the "title" field to the value that was provided on create.
func (u *PullRequestsUpsertBulk) UpdateTitle() *PullRequestsUpsertBulk {
	return u.Update(func(s *PullRequestsUpsert) {
		s.UpdateTitle()
	})
}

// Exec executes the query.
func (u *PullRequestsUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("db: OnConflict was set for builder %d. Set it on the PullRequestsCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for PullRequestsCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PullRequestsUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
