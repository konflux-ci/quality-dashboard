// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/bugs"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-studio/pkg/storage/ent/db/teams"
)

// BugsQuery is the builder for querying Bugs entities.
type BugsQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.Bugs
	withBugs   *TeamsQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the BugsQuery builder.
func (bq *BugsQuery) Where(ps ...predicate.Bugs) *BugsQuery {
	bq.predicates = append(bq.predicates, ps...)
	return bq
}

// Limit the number of records to be returned by this query.
func (bq *BugsQuery) Limit(limit int) *BugsQuery {
	bq.ctx.Limit = &limit
	return bq
}

// Offset to start from.
func (bq *BugsQuery) Offset(offset int) *BugsQuery {
	bq.ctx.Offset = &offset
	return bq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (bq *BugsQuery) Unique(unique bool) *BugsQuery {
	bq.ctx.Unique = &unique
	return bq
}

// Order specifies how the records should be ordered.
func (bq *BugsQuery) Order(o ...OrderFunc) *BugsQuery {
	bq.order = append(bq.order, o...)
	return bq
}

// QueryBugs chains the current query on the "bugs" edge.
func (bq *BugsQuery) QueryBugs() *TeamsQuery {
	query := (&TeamsClient{config: bq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := bq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := bq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(bugs.Table, bugs.FieldID, selector),
			sqlgraph.To(teams.Table, teams.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, bugs.BugsTable, bugs.BugsColumn),
		)
		fromU = sqlgraph.SetNeighbors(bq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Bugs entity from the query.
// Returns a *NotFoundError when no Bugs was found.
func (bq *BugsQuery) First(ctx context.Context) (*Bugs, error) {
	nodes, err := bq.Limit(1).All(setContextOp(ctx, bq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{bugs.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (bq *BugsQuery) FirstX(ctx context.Context) *Bugs {
	node, err := bq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Bugs ID from the query.
// Returns a *NotFoundError when no Bugs ID was found.
func (bq *BugsQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = bq.Limit(1).IDs(setContextOp(ctx, bq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{bugs.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (bq *BugsQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := bq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Bugs entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Bugs entity is found.
// Returns a *NotFoundError when no Bugs entities are found.
func (bq *BugsQuery) Only(ctx context.Context) (*Bugs, error) {
	nodes, err := bq.Limit(2).All(setContextOp(ctx, bq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{bugs.Label}
	default:
		return nil, &NotSingularError{bugs.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (bq *BugsQuery) OnlyX(ctx context.Context) *Bugs {
	node, err := bq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Bugs ID in the query.
// Returns a *NotSingularError when more than one Bugs ID is found.
// Returns a *NotFoundError when no entities are found.
func (bq *BugsQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = bq.Limit(2).IDs(setContextOp(ctx, bq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{bugs.Label}
	default:
		err = &NotSingularError{bugs.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (bq *BugsQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := bq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of BugsSlice.
func (bq *BugsQuery) All(ctx context.Context) ([]*Bugs, error) {
	ctx = setContextOp(ctx, bq.ctx, "All")
	if err := bq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Bugs, *BugsQuery]()
	return withInterceptors[[]*Bugs](ctx, bq, qr, bq.inters)
}

// AllX is like All, but panics if an error occurs.
func (bq *BugsQuery) AllX(ctx context.Context) []*Bugs {
	nodes, err := bq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Bugs IDs.
func (bq *BugsQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	ctx = setContextOp(ctx, bq.ctx, "IDs")
	if err := bq.Select(bugs.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (bq *BugsQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := bq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (bq *BugsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, bq.ctx, "Count")
	if err := bq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, bq, querierCount[*BugsQuery](), bq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (bq *BugsQuery) CountX(ctx context.Context) int {
	count, err := bq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (bq *BugsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, bq.ctx, "Exist")
	switch _, err := bq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("db: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (bq *BugsQuery) ExistX(ctx context.Context) bool {
	exist, err := bq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the BugsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (bq *BugsQuery) Clone() *BugsQuery {
	if bq == nil {
		return nil
	}
	return &BugsQuery{
		config:     bq.config,
		ctx:        bq.ctx.Clone(),
		order:      append([]OrderFunc{}, bq.order...),
		inters:     append([]Interceptor{}, bq.inters...),
		predicates: append([]predicate.Bugs{}, bq.predicates...),
		withBugs:   bq.withBugs.Clone(),
		// clone intermediate query.
		sql:  bq.sql.Clone(),
		path: bq.path,
	}
}

// WithBugs tells the query-builder to eager-load the nodes that are connected to
// the "bugs" edge. The optional arguments are used to configure the query builder of the edge.
func (bq *BugsQuery) WithBugs(opts ...func(*TeamsQuery)) *BugsQuery {
	query := (&TeamsClient{config: bq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	bq.withBugs = query
	return bq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		JiraKey string `json:"jira_key,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Bugs.Query().
//		GroupBy(bugs.FieldJiraKey).
//		Aggregate(db.Count()).
//		Scan(ctx, &v)
func (bq *BugsQuery) GroupBy(field string, fields ...string) *BugsGroupBy {
	bq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &BugsGroupBy{build: bq}
	grbuild.flds = &bq.ctx.Fields
	grbuild.label = bugs.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		JiraKey string `json:"jira_key,omitempty"`
//	}
//
//	client.Bugs.Query().
//		Select(bugs.FieldJiraKey).
//		Scan(ctx, &v)
func (bq *BugsQuery) Select(fields ...string) *BugsSelect {
	bq.ctx.Fields = append(bq.ctx.Fields, fields...)
	sbuild := &BugsSelect{BugsQuery: bq}
	sbuild.label = bugs.Label
	sbuild.flds, sbuild.scan = &bq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a BugsSelect configured with the given aggregations.
func (bq *BugsQuery) Aggregate(fns ...AggregateFunc) *BugsSelect {
	return bq.Select().Aggregate(fns...)
}

func (bq *BugsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range bq.inters {
		if inter == nil {
			return fmt.Errorf("db: uninitialized interceptor (forgotten import db/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, bq); err != nil {
				return err
			}
		}
	}
	for _, f := range bq.ctx.Fields {
		if !bugs.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("db: invalid field %q for query", f)}
		}
	}
	if bq.path != nil {
		prev, err := bq.path(ctx)
		if err != nil {
			return err
		}
		bq.sql = prev
	}
	return nil
}

func (bq *BugsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Bugs, error) {
	var (
		nodes       = []*Bugs{}
		withFKs     = bq.withFKs
		_spec       = bq.querySpec()
		loadedTypes = [1]bool{
			bq.withBugs != nil,
		}
	)
	if bq.withBugs != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, bugs.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Bugs).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Bugs{config: bq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, bq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := bq.withBugs; query != nil {
		if err := bq.loadBugs(ctx, query, nodes, nil,
			func(n *Bugs, e *Teams) { n.Edges.Bugs = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (bq *BugsQuery) loadBugs(ctx context.Context, query *TeamsQuery, nodes []*Bugs, init func(*Bugs), assign func(*Bugs, *Teams)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*Bugs)
	for i := range nodes {
		if nodes[i].teams_bugs == nil {
			continue
		}
		fk := *nodes[i].teams_bugs
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(teams.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "teams_bugs" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (bq *BugsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := bq.querySpec()
	_spec.Node.Columns = bq.ctx.Fields
	if len(bq.ctx.Fields) > 0 {
		_spec.Unique = bq.ctx.Unique != nil && *bq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, bq.driver, _spec)
}

func (bq *BugsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   bugs.Table,
			Columns: bugs.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: bugs.FieldID,
			},
		},
		From:   bq.sql,
		Unique: true,
	}
	if unique := bq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := bq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, bugs.FieldID)
		for i := range fields {
			if fields[i] != bugs.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := bq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := bq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := bq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := bq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (bq *BugsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(bq.driver.Dialect())
	t1 := builder.Table(bugs.Table)
	columns := bq.ctx.Fields
	if len(columns) == 0 {
		columns = bugs.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if bq.sql != nil {
		selector = bq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if bq.ctx.Unique != nil && *bq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range bq.predicates {
		p(selector)
	}
	for _, p := range bq.order {
		p(selector)
	}
	if offset := bq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := bq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// BugsGroupBy is the group-by builder for Bugs entities.
type BugsGroupBy struct {
	selector
	build *BugsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (bgb *BugsGroupBy) Aggregate(fns ...AggregateFunc) *BugsGroupBy {
	bgb.fns = append(bgb.fns, fns...)
	return bgb
}

// Scan applies the selector query and scans the result into the given value.
func (bgb *BugsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bgb.build.ctx, "GroupBy")
	if err := bgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BugsQuery, *BugsGroupBy](ctx, bgb.build, bgb, bgb.build.inters, v)
}

func (bgb *BugsGroupBy) sqlScan(ctx context.Context, root *BugsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(bgb.fns))
	for _, fn := range bgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*bgb.flds)+len(bgb.fns))
		for _, f := range *bgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*bgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// BugsSelect is the builder for selecting fields of Bugs entities.
type BugsSelect struct {
	*BugsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (bs *BugsSelect) Aggregate(fns ...AggregateFunc) *BugsSelect {
	bs.fns = append(bs.fns, fns...)
	return bs
}

// Scan applies the selector query and scans the result into the given value.
func (bs *BugsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bs.ctx, "Select")
	if err := bs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BugsQuery, *BugsSelect](ctx, bs.BugsQuery, bs, bs.inters, v)
}

func (bs *BugsSelect) sqlScan(ctx context.Context, root *BugsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(bs.fns))
	for _, fn := range bs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*bs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
