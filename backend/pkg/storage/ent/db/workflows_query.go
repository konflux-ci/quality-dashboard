// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/repository"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/workflows"
)

// WorkflowsQuery is the builder for querying Workflows entities.
type WorkflowsQuery struct {
	config
	ctx           *QueryContext
	order         []OrderFunc
	inters        []Interceptor
	predicates    []predicate.Workflows
	withWorkflows *RepositoryQuery
	withFKs       bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the WorkflowsQuery builder.
func (wq *WorkflowsQuery) Where(ps ...predicate.Workflows) *WorkflowsQuery {
	wq.predicates = append(wq.predicates, ps...)
	return wq
}

// Limit the number of records to be returned by this query.
func (wq *WorkflowsQuery) Limit(limit int) *WorkflowsQuery {
	wq.ctx.Limit = &limit
	return wq
}

// Offset to start from.
func (wq *WorkflowsQuery) Offset(offset int) *WorkflowsQuery {
	wq.ctx.Offset = &offset
	return wq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (wq *WorkflowsQuery) Unique(unique bool) *WorkflowsQuery {
	wq.ctx.Unique = &unique
	return wq
}

// Order specifies how the records should be ordered.
func (wq *WorkflowsQuery) Order(o ...OrderFunc) *WorkflowsQuery {
	wq.order = append(wq.order, o...)
	return wq
}

// QueryWorkflows chains the current query on the "workflows" edge.
func (wq *WorkflowsQuery) QueryWorkflows() *RepositoryQuery {
	query := (&RepositoryClient{config: wq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(workflows.Table, workflows.FieldID, selector),
			sqlgraph.To(repository.Table, repository.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, workflows.WorkflowsTable, workflows.WorkflowsColumn),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Workflows entity from the query.
// Returns a *NotFoundError when no Workflows was found.
func (wq *WorkflowsQuery) First(ctx context.Context) (*Workflows, error) {
	nodes, err := wq.Limit(1).All(setContextOp(ctx, wq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{workflows.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wq *WorkflowsQuery) FirstX(ctx context.Context) *Workflows {
	node, err := wq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Workflows ID from the query.
// Returns a *NotFoundError when no Workflows ID was found.
func (wq *WorkflowsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wq.Limit(1).IDs(setContextOp(ctx, wq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{workflows.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (wq *WorkflowsQuery) FirstIDX(ctx context.Context) int {
	id, err := wq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Workflows entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Workflows entity is found.
// Returns a *NotFoundError when no Workflows entities are found.
func (wq *WorkflowsQuery) Only(ctx context.Context) (*Workflows, error) {
	nodes, err := wq.Limit(2).All(setContextOp(ctx, wq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{workflows.Label}
	default:
		return nil, &NotSingularError{workflows.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wq *WorkflowsQuery) OnlyX(ctx context.Context) *Workflows {
	node, err := wq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Workflows ID in the query.
// Returns a *NotSingularError when more than one Workflows ID is found.
// Returns a *NotFoundError when no entities are found.
func (wq *WorkflowsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wq.Limit(2).IDs(setContextOp(ctx, wq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{workflows.Label}
	default:
		err = &NotSingularError{workflows.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (wq *WorkflowsQuery) OnlyIDX(ctx context.Context) int {
	id, err := wq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of WorkflowsSlice.
func (wq *WorkflowsQuery) All(ctx context.Context) ([]*Workflows, error) {
	ctx = setContextOp(ctx, wq.ctx, "All")
	if err := wq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Workflows, *WorkflowsQuery]()
	return withInterceptors[[]*Workflows](ctx, wq, qr, wq.inters)
}

// AllX is like All, but panics if an error occurs.
func (wq *WorkflowsQuery) AllX(ctx context.Context) []*Workflows {
	nodes, err := wq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Workflows IDs.
func (wq *WorkflowsQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	ctx = setContextOp(ctx, wq.ctx, "IDs")
	if err := wq.Select(workflows.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wq *WorkflowsQuery) IDsX(ctx context.Context) []int {
	ids, err := wq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wq *WorkflowsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, wq.ctx, "Count")
	if err := wq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, wq, querierCount[*WorkflowsQuery](), wq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (wq *WorkflowsQuery) CountX(ctx context.Context) int {
	count, err := wq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wq *WorkflowsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, wq.ctx, "Exist")
	switch _, err := wq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("db: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (wq *WorkflowsQuery) ExistX(ctx context.Context) bool {
	exist, err := wq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the WorkflowsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wq *WorkflowsQuery) Clone() *WorkflowsQuery {
	if wq == nil {
		return nil
	}
	return &WorkflowsQuery{
		config:        wq.config,
		ctx:           wq.ctx.Clone(),
		order:         append([]OrderFunc{}, wq.order...),
		inters:        append([]Interceptor{}, wq.inters...),
		predicates:    append([]predicate.Workflows{}, wq.predicates...),
		withWorkflows: wq.withWorkflows.Clone(),
		// clone intermediate query.
		sql:  wq.sql.Clone(),
		path: wq.path,
	}
}

// WithWorkflows tells the query-builder to eager-load the nodes that are connected to
// the "workflows" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WorkflowsQuery) WithWorkflows(opts ...func(*RepositoryQuery)) *WorkflowsQuery {
	query := (&RepositoryClient{config: wq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	wq.withWorkflows = query
	return wq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		WorkflowID uuid.UUID `json:"workflow_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Workflows.Query().
//		GroupBy(workflows.FieldWorkflowID).
//		Aggregate(db.Count()).
//		Scan(ctx, &v)
func (wq *WorkflowsQuery) GroupBy(field string, fields ...string) *WorkflowsGroupBy {
	wq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &WorkflowsGroupBy{build: wq}
	grbuild.flds = &wq.ctx.Fields
	grbuild.label = workflows.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		WorkflowID uuid.UUID `json:"workflow_id,omitempty"`
//	}
//
//	client.Workflows.Query().
//		Select(workflows.FieldWorkflowID).
//		Scan(ctx, &v)
func (wq *WorkflowsQuery) Select(fields ...string) *WorkflowsSelect {
	wq.ctx.Fields = append(wq.ctx.Fields, fields...)
	sbuild := &WorkflowsSelect{WorkflowsQuery: wq}
	sbuild.label = workflows.Label
	sbuild.flds, sbuild.scan = &wq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a WorkflowsSelect configured with the given aggregations.
func (wq *WorkflowsQuery) Aggregate(fns ...AggregateFunc) *WorkflowsSelect {
	return wq.Select().Aggregate(fns...)
}

func (wq *WorkflowsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range wq.inters {
		if inter == nil {
			return fmt.Errorf("db: uninitialized interceptor (forgotten import db/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, wq); err != nil {
				return err
			}
		}
	}
	for _, f := range wq.ctx.Fields {
		if !workflows.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("db: invalid field %q for query", f)}
		}
	}
	if wq.path != nil {
		prev, err := wq.path(ctx)
		if err != nil {
			return err
		}
		wq.sql = prev
	}
	return nil
}

func (wq *WorkflowsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Workflows, error) {
	var (
		nodes       = []*Workflows{}
		withFKs     = wq.withFKs
		_spec       = wq.querySpec()
		loadedTypes = [1]bool{
			wq.withWorkflows != nil,
		}
	)
	if wq.withWorkflows != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, workflows.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Workflows).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Workflows{config: wq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, wq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := wq.withWorkflows; query != nil {
		if err := wq.loadWorkflows(ctx, query, nodes, nil,
			func(n *Workflows, e *Repository) { n.Edges.Workflows = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (wq *WorkflowsQuery) loadWorkflows(ctx context.Context, query *RepositoryQuery, nodes []*Workflows, init func(*Workflows), assign func(*Workflows, *Repository)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*Workflows)
	for i := range nodes {
		if nodes[i].repository_workflows == nil {
			continue
		}
		fk := *nodes[i].repository_workflows
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(repository.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "repository_workflows" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (wq *WorkflowsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wq.querySpec()
	_spec.Node.Columns = wq.ctx.Fields
	if len(wq.ctx.Fields) > 0 {
		_spec.Unique = wq.ctx.Unique != nil && *wq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, wq.driver, _spec)
}

func (wq *WorkflowsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   workflows.Table,
			Columns: workflows.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: workflows.FieldID,
			},
		},
		From:   wq.sql,
		Unique: true,
	}
	if unique := wq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := wq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, workflows.FieldID)
		for i := range fields {
			if fields[i] != workflows.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := wq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wq *WorkflowsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(wq.driver.Dialect())
	t1 := builder.Table(workflows.Table)
	columns := wq.ctx.Fields
	if len(columns) == 0 {
		columns = workflows.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if wq.sql != nil {
		selector = wq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if wq.ctx.Unique != nil && *wq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range wq.predicates {
		p(selector)
	}
	for _, p := range wq.order {
		p(selector)
	}
	if offset := wq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WorkflowsGroupBy is the group-by builder for Workflows entities.
type WorkflowsGroupBy struct {
	selector
	build *WorkflowsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wgb *WorkflowsGroupBy) Aggregate(fns ...AggregateFunc) *WorkflowsGroupBy {
	wgb.fns = append(wgb.fns, fns...)
	return wgb
}

// Scan applies the selector query and scans the result into the given value.
func (wgb *WorkflowsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, wgb.build.ctx, "GroupBy")
	if err := wgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WorkflowsQuery, *WorkflowsGroupBy](ctx, wgb.build, wgb, wgb.build.inters, v)
}

func (wgb *WorkflowsGroupBy) sqlScan(ctx context.Context, root *WorkflowsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(wgb.fns))
	for _, fn := range wgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*wgb.flds)+len(wgb.fns))
		for _, f := range *wgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*wgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// WorkflowsSelect is the builder for selecting fields of Workflows entities.
type WorkflowsSelect struct {
	*WorkflowsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ws *WorkflowsSelect) Aggregate(fns ...AggregateFunc) *WorkflowsSelect {
	ws.fns = append(ws.fns, fns...)
	return ws
}

// Scan applies the selector query and scans the result into the given value.
func (ws *WorkflowsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ws.ctx, "Select")
	if err := ws.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WorkflowsQuery, *WorkflowsSelect](ctx, ws.WorkflowsQuery, ws, ws.inters, v)
}

func (ws *WorkflowsSelect) sqlScan(ctx context.Context, root *WorkflowsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ws.fns))
	for _, fn := range ws.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ws.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ws.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
