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
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/users"
)

// UsersCreate is the builder for creating a Users entity.
type UsersCreate struct {
	config
	mutation *UsersMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUserEmail sets the "user_email" field.
func (uc *UsersCreate) SetUserEmail(s string) *UsersCreate {
	uc.mutation.SetUserEmail(s)
	return uc
}

// SetConfig sets the "config" field.
func (uc *UsersCreate) SetConfig(s string) *UsersCreate {
	uc.mutation.SetConfig(s)
	return uc
}

// SetID sets the "id" field.
func (uc *UsersCreate) SetID(u uuid.UUID) *UsersCreate {
	uc.mutation.SetID(u)
	return uc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (uc *UsersCreate) SetNillableID(u *uuid.UUID) *UsersCreate {
	if u != nil {
		uc.SetID(*u)
	}
	return uc
}

// Mutation returns the UsersMutation object of the builder.
func (uc *UsersCreate) Mutation() *UsersMutation {
	return uc.mutation
}

// Save creates the Users in the database.
func (uc *UsersCreate) Save(ctx context.Context) (*Users, error) {
	uc.defaults()
	return withHooks[*Users, UsersMutation](ctx, uc.sqlSave, uc.mutation, uc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UsersCreate) SaveX(ctx context.Context) *Users {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uc *UsersCreate) Exec(ctx context.Context) error {
	_, err := uc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uc *UsersCreate) ExecX(ctx context.Context) {
	if err := uc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uc *UsersCreate) defaults() {
	if _, ok := uc.mutation.ID(); !ok {
		v := users.DefaultID()
		uc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uc *UsersCreate) check() error {
	if _, ok := uc.mutation.UserEmail(); !ok {
		return &ValidationError{Name: "user_email", err: errors.New(`db: missing required field "Users.user_email"`)}
	}
	if _, ok := uc.mutation.Config(); !ok {
		return &ValidationError{Name: "config", err: errors.New(`db: missing required field "Users.config"`)}
	}
	return nil
}

func (uc *UsersCreate) sqlSave(ctx context.Context) (*Users, error) {
	if err := uc.check(); err != nil {
		return nil, err
	}
	_node, _spec := uc.createSpec()
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
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
	uc.mutation.id = &_node.ID
	uc.mutation.done = true
	return _node, nil
}

func (uc *UsersCreate) createSpec() (*Users, *sqlgraph.CreateSpec) {
	var (
		_node = &Users{config: uc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: users.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: users.FieldID,
			},
		}
	)
	_spec.OnConflict = uc.conflict
	if id, ok := uc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := uc.mutation.UserEmail(); ok {
		_spec.SetField(users.FieldUserEmail, field.TypeString, value)
		_node.UserEmail = value
	}
	if value, ok := uc.mutation.Config(); ok {
		_spec.SetField(users.FieldConfig, field.TypeString, value)
		_node.Config = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Users.Create().
//		SetUserEmail(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UsersUpsert) {
//			SetUserEmail(v+v).
//		}).
//		Exec(ctx)
func (uc *UsersCreate) OnConflict(opts ...sql.ConflictOption) *UsersUpsertOne {
	uc.conflict = opts
	return &UsersUpsertOne{
		create: uc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Users.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (uc *UsersCreate) OnConflictColumns(columns ...string) *UsersUpsertOne {
	uc.conflict = append(uc.conflict, sql.ConflictColumns(columns...))
	return &UsersUpsertOne{
		create: uc,
	}
}

type (
	// UsersUpsertOne is the builder for "upsert"-ing
	//  one Users node.
	UsersUpsertOne struct {
		create *UsersCreate
	}

	// UsersUpsert is the "OnConflict" setter.
	UsersUpsert struct {
		*sql.UpdateSet
	}
)

// SetUserEmail sets the "user_email" field.
func (u *UsersUpsert) SetUserEmail(v string) *UsersUpsert {
	u.Set(users.FieldUserEmail, v)
	return u
}

// UpdateUserEmail sets the "user_email" field to the value that was provided on create.
func (u *UsersUpsert) UpdateUserEmail() *UsersUpsert {
	u.SetExcluded(users.FieldUserEmail)
	return u
}

// SetConfig sets the "config" field.
func (u *UsersUpsert) SetConfig(v string) *UsersUpsert {
	u.Set(users.FieldConfig, v)
	return u
}

// UpdateConfig sets the "config" field to the value that was provided on create.
func (u *UsersUpsert) UpdateConfig() *UsersUpsert {
	u.SetExcluded(users.FieldConfig)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Users.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(users.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *UsersUpsertOne) UpdateNewValues() *UsersUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(users.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Users.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *UsersUpsertOne) Ignore() *UsersUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UsersUpsertOne) DoNothing() *UsersUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UsersCreate.OnConflict
// documentation for more info.
func (u *UsersUpsertOne) Update(set func(*UsersUpsert)) *UsersUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UsersUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserEmail sets the "user_email" field.
func (u *UsersUpsertOne) SetUserEmail(v string) *UsersUpsertOne {
	return u.Update(func(s *UsersUpsert) {
		s.SetUserEmail(v)
	})
}

// UpdateUserEmail sets the "user_email" field to the value that was provided on create.
func (u *UsersUpsertOne) UpdateUserEmail() *UsersUpsertOne {
	return u.Update(func(s *UsersUpsert) {
		s.UpdateUserEmail()
	})
}

// SetConfig sets the "config" field.
func (u *UsersUpsertOne) SetConfig(v string) *UsersUpsertOne {
	return u.Update(func(s *UsersUpsert) {
		s.SetConfig(v)
	})
}

// UpdateConfig sets the "config" field to the value that was provided on create.
func (u *UsersUpsertOne) UpdateConfig() *UsersUpsertOne {
	return u.Update(func(s *UsersUpsert) {
		s.UpdateConfig()
	})
}

// Exec executes the query.
func (u *UsersUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for UsersCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UsersUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *UsersUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("db: UsersUpsertOne.ID is not supported by MySQL driver. Use UsersUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *UsersUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// UsersCreateBulk is the builder for creating many Users entities in bulk.
type UsersCreateBulk struct {
	config
	builders []*UsersCreate
	conflict []sql.ConflictOption
}

// Save creates the Users entities in the database.
func (ucb *UsersCreateBulk) Save(ctx context.Context) ([]*Users, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ucb.builders))
	nodes := make([]*Users, len(ucb.builders))
	mutators := make([]Mutator, len(ucb.builders))
	for i := range ucb.builders {
		func(i int, root context.Context) {
			builder := ucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UsersMutation)
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
					_, err = mutators[i+1].Mutate(root, ucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ucb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ucb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ucb *UsersCreateBulk) SaveX(ctx context.Context) []*Users {
	v, err := ucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ucb *UsersCreateBulk) Exec(ctx context.Context) error {
	_, err := ucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ucb *UsersCreateBulk) ExecX(ctx context.Context) {
	if err := ucb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Users.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UsersUpsert) {
//			SetUserEmail(v+v).
//		}).
//		Exec(ctx)
func (ucb *UsersCreateBulk) OnConflict(opts ...sql.ConflictOption) *UsersUpsertBulk {
	ucb.conflict = opts
	return &UsersUpsertBulk{
		create: ucb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Users.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ucb *UsersCreateBulk) OnConflictColumns(columns ...string) *UsersUpsertBulk {
	ucb.conflict = append(ucb.conflict, sql.ConflictColumns(columns...))
	return &UsersUpsertBulk{
		create: ucb,
	}
}

// UsersUpsertBulk is the builder for "upsert"-ing
// a bulk of Users nodes.
type UsersUpsertBulk struct {
	create *UsersCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Users.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(users.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *UsersUpsertBulk) UpdateNewValues() *UsersUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(users.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Users.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *UsersUpsertBulk) Ignore() *UsersUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UsersUpsertBulk) DoNothing() *UsersUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UsersCreateBulk.OnConflict
// documentation for more info.
func (u *UsersUpsertBulk) Update(set func(*UsersUpsert)) *UsersUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UsersUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserEmail sets the "user_email" field.
func (u *UsersUpsertBulk) SetUserEmail(v string) *UsersUpsertBulk {
	return u.Update(func(s *UsersUpsert) {
		s.SetUserEmail(v)
	})
}

// UpdateUserEmail sets the "user_email" field to the value that was provided on create.
func (u *UsersUpsertBulk) UpdateUserEmail() *UsersUpsertBulk {
	return u.Update(func(s *UsersUpsert) {
		s.UpdateUserEmail()
	})
}

// SetConfig sets the "config" field.
func (u *UsersUpsertBulk) SetConfig(v string) *UsersUpsertBulk {
	return u.Update(func(s *UsersUpsert) {
		s.SetConfig(v)
	})
}

// UpdateConfig sets the "config" field to the value that was provided on create.
func (u *UsersUpsertBulk) UpdateConfig() *UsersUpsertBulk {
	return u.Update(func(s *UsersUpsert) {
		s.UpdateConfig()
	})
}

// Exec executes the query.
func (u *UsersUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("db: OnConflict was set for builder %d. Set it on the UsersCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("db: missing options for UsersCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UsersUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
