// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/prowsuites"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/repository"
)

// ProwSuitesQuery is the builder for querying ProwSuites entities.
type ProwSuitesQuery struct {
	config
	ctx            *QueryContext
	order          []OrderFunc
	inters         []Interceptor
	predicates     []predicate.ProwSuites
	withProwSuites *RepositoryQuery
	withFKs        bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ProwSuitesQuery builder.
func (psq *ProwSuitesQuery) Where(ps ...predicate.ProwSuites) *ProwSuitesQuery {
	psq.predicates = append(psq.predicates, ps...)
	return psq
}

// Limit the number of records to be returned by this query.
func (psq *ProwSuitesQuery) Limit(limit int) *ProwSuitesQuery {
	psq.ctx.Limit = &limit
	return psq
}

// Offset to start from.
func (psq *ProwSuitesQuery) Offset(offset int) *ProwSuitesQuery {
	psq.ctx.Offset = &offset
	return psq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (psq *ProwSuitesQuery) Unique(unique bool) *ProwSuitesQuery {
	psq.ctx.Unique = &unique
	return psq
}

// Order specifies how the records should be ordered.
func (psq *ProwSuitesQuery) Order(o ...OrderFunc) *ProwSuitesQuery {
	psq.order = append(psq.order, o...)
	return psq
}

// QueryProwSuites chains the current query on the "prow_suites" edge.
func (psq *ProwSuitesQuery) QueryProwSuites() *RepositoryQuery {
	query := (&RepositoryClient{config: psq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := psq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := psq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(prowsuites.Table, prowsuites.FieldID, selector),
			sqlgraph.To(repository.Table, repository.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, prowsuites.ProwSuitesTable, prowsuites.ProwSuitesColumn),
		)
		fromU = sqlgraph.SetNeighbors(psq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ProwSuites entity from the query.
// Returns a *NotFoundError when no ProwSuites was found.
func (psq *ProwSuitesQuery) First(ctx context.Context) (*ProwSuites, error) {
	nodes, err := psq.Limit(1).All(setContextOp(ctx, psq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{prowsuites.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (psq *ProwSuitesQuery) FirstX(ctx context.Context) *ProwSuites {
	node, err := psq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ProwSuites ID from the query.
// Returns a *NotFoundError when no ProwSuites ID was found.
func (psq *ProwSuitesQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = psq.Limit(1).IDs(setContextOp(ctx, psq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{prowsuites.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (psq *ProwSuitesQuery) FirstIDX(ctx context.Context) int {
	id, err := psq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ProwSuites entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ProwSuites entity is found.
// Returns a *NotFoundError when no ProwSuites entities are found.
func (psq *ProwSuitesQuery) Only(ctx context.Context) (*ProwSuites, error) {
	nodes, err := psq.Limit(2).All(setContextOp(ctx, psq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{prowsuites.Label}
	default:
		return nil, &NotSingularError{prowsuites.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (psq *ProwSuitesQuery) OnlyX(ctx context.Context) *ProwSuites {
	node, err := psq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ProwSuites ID in the query.
// Returns a *NotSingularError when more than one ProwSuites ID is found.
// Returns a *NotFoundError when no entities are found.
func (psq *ProwSuitesQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = psq.Limit(2).IDs(setContextOp(ctx, psq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{prowsuites.Label}
	default:
		err = &NotSingularError{prowsuites.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (psq *ProwSuitesQuery) OnlyIDX(ctx context.Context) int {
	id, err := psq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ProwSuitesSlice.
func (psq *ProwSuitesQuery) All(ctx context.Context) ([]*ProwSuites, error) {
	ctx = setContextOp(ctx, psq.ctx, "All")
	if err := psq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ProwSuites, *ProwSuitesQuery]()
	return withInterceptors[[]*ProwSuites](ctx, psq, qr, psq.inters)
}

// AllX is like All, but panics if an error occurs.
func (psq *ProwSuitesQuery) AllX(ctx context.Context) []*ProwSuites {
	nodes, err := psq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ProwSuites IDs.
func (psq *ProwSuitesQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	ctx = setContextOp(ctx, psq.ctx, "IDs")
	if err := psq.Select(prowsuites.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (psq *ProwSuitesQuery) IDsX(ctx context.Context) []int {
	ids, err := psq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (psq *ProwSuitesQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, psq.ctx, "Count")
	if err := psq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, psq, querierCount[*ProwSuitesQuery](), psq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (psq *ProwSuitesQuery) CountX(ctx context.Context) int {
	count, err := psq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (psq *ProwSuitesQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, psq.ctx, "Exist")
	switch _, err := psq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("db: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (psq *ProwSuitesQuery) ExistX(ctx context.Context) bool {
	exist, err := psq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ProwSuitesQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (psq *ProwSuitesQuery) Clone() *ProwSuitesQuery {
	if psq == nil {
		return nil
	}
	return &ProwSuitesQuery{
		config:         psq.config,
		ctx:            psq.ctx.Clone(),
		order:          append([]OrderFunc{}, psq.order...),
		inters:         append([]Interceptor{}, psq.inters...),
		predicates:     append([]predicate.ProwSuites{}, psq.predicates...),
		withProwSuites: psq.withProwSuites.Clone(),
		// clone intermediate query.
		sql:  psq.sql.Clone(),
		path: psq.path,
	}
}

// WithProwSuites tells the query-builder to eager-load the nodes that are connected to
// the "prow_suites" edge. The optional arguments are used to configure the query builder of the edge.
func (psq *ProwSuitesQuery) WithProwSuites(opts ...func(*RepositoryQuery)) *ProwSuitesQuery {
	query := (&RepositoryClient{config: psq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	psq.withProwSuites = query
	return psq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		JobID string `json:"job_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ProwSuites.Query().
//		GroupBy(prowsuites.FieldJobID).
//		Aggregate(db.Count()).
//		Scan(ctx, &v)
//
func (psq *ProwSuitesQuery) GroupBy(field string, fields ...string) *ProwSuitesGroupBy {
	psq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ProwSuitesGroupBy{build: psq}
	grbuild.flds = &psq.ctx.Fields
	grbuild.label = prowsuites.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		JobID string `json:"job_id,omitempty"`
//	}
//
//	client.ProwSuites.Query().
//		Select(prowsuites.FieldJobID).
//		Scan(ctx, &v)
//
func (psq *ProwSuitesQuery) Select(fields ...string) *ProwSuitesSelect {
	psq.ctx.Fields = append(psq.ctx.Fields, fields...)
	sbuild := &ProwSuitesSelect{ProwSuitesQuery: psq}
	sbuild.label = prowsuites.Label
	sbuild.flds, sbuild.scan = &psq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ProwSuitesSelect configured with the given aggregations.
func (psq *ProwSuitesQuery) Aggregate(fns ...AggregateFunc) *ProwSuitesSelect {
	return psq.Select().Aggregate(fns...)
}

func (psq *ProwSuitesQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range psq.inters {
		if inter == nil {
			return fmt.Errorf("db: uninitialized interceptor (forgotten import db/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, psq); err != nil {
				return err
			}
		}
	}
	for _, f := range psq.ctx.Fields {
		if !prowsuites.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("db: invalid field %q for query", f)}
		}
	}
	if psq.path != nil {
		prev, err := psq.path(ctx)
		if err != nil {
			return err
		}
		psq.sql = prev
	}
	return nil
}

func (psq *ProwSuitesQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ProwSuites, error) {
	var (
		nodes       = []*ProwSuites{}
		withFKs     = psq.withFKs
		_spec       = psq.querySpec()
		loadedTypes = [1]bool{
			psq.withProwSuites != nil,
		}
	)
	if psq.withProwSuites != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, prowsuites.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ProwSuites).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ProwSuites{config: psq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, psq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := psq.withProwSuites; query != nil {
		if err := psq.loadProwSuites(ctx, query, nodes, nil,
			func(n *ProwSuites, e *Repository) { n.Edges.ProwSuites = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (psq *ProwSuitesQuery) loadProwSuites(ctx context.Context, query *RepositoryQuery, nodes []*ProwSuites, init func(*ProwSuites), assign func(*ProwSuites, *Repository)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ProwSuites)
	for i := range nodes {
		if nodes[i].repository_prow_suites == nil {
			continue
		}
		fk := *nodes[i].repository_prow_suites
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
			return fmt.Errorf(`unexpected foreign-key "repository_prow_suites" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (psq *ProwSuitesQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := psq.querySpec()
	_spec.Node.Columns = psq.ctx.Fields
	if len(psq.ctx.Fields) > 0 {
		_spec.Unique = psq.ctx.Unique != nil && *psq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, psq.driver, _spec)
}

func (psq *ProwSuitesQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   prowsuites.Table,
			Columns: prowsuites.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: prowsuites.FieldID,
			},
		},
		From:   psq.sql,
		Unique: true,
	}
	if unique := psq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := psq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, prowsuites.FieldID)
		for i := range fields {
			if fields[i] != prowsuites.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := psq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := psq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := psq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := psq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (psq *ProwSuitesQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(psq.driver.Dialect())
	t1 := builder.Table(prowsuites.Table)
	columns := psq.ctx.Fields
	if len(columns) == 0 {
		columns = prowsuites.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if psq.sql != nil {
		selector = psq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if psq.ctx.Unique != nil && *psq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range psq.predicates {
		p(selector)
	}
	for _, p := range psq.order {
		p(selector)
	}
	if offset := psq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := psq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ProwSuitesGroupBy is the group-by builder for ProwSuites entities.
type ProwSuitesGroupBy struct {
	selector
	build *ProwSuitesQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (psgb *ProwSuitesGroupBy) Aggregate(fns ...AggregateFunc) *ProwSuitesGroupBy {
	psgb.fns = append(psgb.fns, fns...)
	return psgb
}

// Scan applies the selector query and scans the result into the given value.
func (psgb *ProwSuitesGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, psgb.build.ctx, "GroupBy")
	if err := psgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProwSuitesQuery, *ProwSuitesGroupBy](ctx, psgb.build, psgb, psgb.build.inters, v)
}

func (psgb *ProwSuitesGroupBy) sqlScan(ctx context.Context, root *ProwSuitesQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(psgb.fns))
	for _, fn := range psgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*psgb.flds)+len(psgb.fns))
		for _, f := range *psgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*psgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := psgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ProwSuitesSelect is the builder for selecting fields of ProwSuites entities.
type ProwSuitesSelect struct {
	*ProwSuitesQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pss *ProwSuitesSelect) Aggregate(fns ...AggregateFunc) *ProwSuitesSelect {
	pss.fns = append(pss.fns, fns...)
	return pss
}

// Scan applies the selector query and scans the result into the given value.
func (pss *ProwSuitesSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pss.ctx, "Select")
	if err := pss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProwSuitesQuery, *ProwSuitesSelect](ctx, pss.ProwSuitesQuery, pss, pss.inters, v)
}

func (pss *ProwSuitesSelect) sqlScan(ctx context.Context, root *ProwSuitesQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pss.fns))
	for _, fn := range pss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
