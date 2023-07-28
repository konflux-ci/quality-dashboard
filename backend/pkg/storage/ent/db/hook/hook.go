// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

// The BugsFunc type is an adapter to allow the use of ordinary
// function as Bugs mutator.
type BugsFunc func(context.Context, *db.BugsMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f BugsFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.BugsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.BugsMutation", m)
}

// The CodeCovFunc type is an adapter to allow the use of ordinary
// function as CodeCov mutator.
type CodeCovFunc func(context.Context, *db.CodeCovMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f CodeCovFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.CodeCovMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.CodeCovMutation", m)
}

// The FailureFunc type is an adapter to allow the use of ordinary
// function as Failure mutator.
type FailureFunc func(context.Context, *db.FailureMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f FailureFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.FailureMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.FailureMutation", m)
}

// The ProwJobsFunc type is an adapter to allow the use of ordinary
// function as ProwJobs mutator.
type ProwJobsFunc func(context.Context, *db.ProwJobsMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f ProwJobsFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.ProwJobsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.ProwJobsMutation", m)
}

// The ProwSuitesFunc type is an adapter to allow the use of ordinary
// function as ProwSuites mutator.
type ProwSuitesFunc func(context.Context, *db.ProwSuitesMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f ProwSuitesFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.ProwSuitesMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.ProwSuitesMutation", m)
}

// The PullRequestsFunc type is an adapter to allow the use of ordinary
// function as PullRequests mutator.
type PullRequestsFunc func(context.Context, *db.PullRequestsMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f PullRequestsFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.PullRequestsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.PullRequestsMutation", m)
}

// The RepositoryFunc type is an adapter to allow the use of ordinary
// function as Repository mutator.
type RepositoryFunc func(context.Context, *db.RepositoryMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f RepositoryFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.RepositoryMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.RepositoryMutation", m)
}

// The TeamsFunc type is an adapter to allow the use of ordinary
// function as Teams mutator.
type TeamsFunc func(context.Context, *db.TeamsMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f TeamsFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.TeamsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.TeamsMutation", m)
}

// The WorkflowsFunc type is an adapter to allow the use of ordinary
// function as Workflows mutator.
type WorkflowsFunc func(context.Context, *db.WorkflowsMutation) (db.Value, error)

// Mutate calls f(ctx, m).
func (f WorkflowsFunc) Mutate(ctx context.Context, m db.Mutation) (db.Value, error) {
	if mv, ok := m.(*db.WorkflowsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *db.WorkflowsMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, db.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m db.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m db.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m db.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op db.Op) Condition {
	return func(_ context.Context, m db.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m db.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m db.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m db.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk db.Hook, cond Condition) db.Hook {
	return func(next db.Mutator) db.Mutator {
		return db.MutateFunc(func(ctx context.Context, m db.Mutation) (db.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, db.Delete|db.Create)
func On(hk db.Hook, op db.Op) db.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, db.Update|db.UpdateOne)
func Unless(hk db.Hook, op db.Op) db.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) db.Hook {
	return func(db.Mutator) db.Mutator {
		return db.MutateFunc(func(context.Context, db.Mutation) (db.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []db.Hook {
//		return []db.Hook{
//			Reject(db.Delete|db.Update),
//		}
//	}
func Reject(op db.Op) db.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []db.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...db.Hook) Chain {
	return Chain{append([]db.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() db.Hook {
	return func(mutator db.Mutator) db.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...db.Hook) Chain {
	newHooks := make([]db.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
