// Code generated by ent, DO NOT EDIT.

package db

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/predicate"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/prowjobs"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/repository"
	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db/tektontasks"
)

// ProwJobsQuery is the builder for querying ProwJobs entities.
type ProwJobsQuery struct {
	config
	ctx             *QueryContext
	order           []prowjobs.OrderOption
	inters          []Interceptor
	predicates      []predicate.ProwJobs
	withRepository  *RepositoryQuery
	withTektonTasks *TektonTasksQuery
	withFKs         bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ProwJobsQuery builder.
func (pjq *ProwJobsQuery) Where(ps ...predicate.ProwJobs) *ProwJobsQuery {
	pjq.predicates = append(pjq.predicates, ps...)
	return pjq
}

// Limit the number of records to be returned by this query.
func (pjq *ProwJobsQuery) Limit(limit int) *ProwJobsQuery {
	pjq.ctx.Limit = &limit
	return pjq
}

// Offset to start from.
func (pjq *ProwJobsQuery) Offset(offset int) *ProwJobsQuery {
	pjq.ctx.Offset = &offset
	return pjq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pjq *ProwJobsQuery) Unique(unique bool) *ProwJobsQuery {
	pjq.ctx.Unique = &unique
	return pjq
}

// Order specifies how the records should be ordered.
func (pjq *ProwJobsQuery) Order(o ...prowjobs.OrderOption) *ProwJobsQuery {
	pjq.order = append(pjq.order, o...)
	return pjq
}

// QueryRepository chains the current query on the "repository" edge.
func (pjq *ProwJobsQuery) QueryRepository() *RepositoryQuery {
	query := (&RepositoryClient{config: pjq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pjq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pjq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(prowjobs.Table, prowjobs.FieldID, selector),
			sqlgraph.To(repository.Table, repository.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, prowjobs.RepositoryTable, prowjobs.RepositoryColumn),
		)
		fromU = sqlgraph.SetNeighbors(pjq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTektonTasks chains the current query on the "tekton_tasks" edge.
func (pjq *ProwJobsQuery) QueryTektonTasks() *TektonTasksQuery {
	query := (&TektonTasksClient{config: pjq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pjq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pjq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(prowjobs.Table, prowjobs.FieldID, selector),
			sqlgraph.To(tektontasks.Table, tektontasks.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, prowjobs.TektonTasksTable, prowjobs.TektonTasksColumn),
		)
		fromU = sqlgraph.SetNeighbors(pjq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ProwJobs entity from the query.
// Returns a *NotFoundError when no ProwJobs was found.
func (pjq *ProwJobsQuery) First(ctx context.Context) (*ProwJobs, error) {
	nodes, err := pjq.Limit(1).All(setContextOp(ctx, pjq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{prowjobs.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pjq *ProwJobsQuery) FirstX(ctx context.Context) *ProwJobs {
	node, err := pjq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ProwJobs ID from the query.
// Returns a *NotFoundError when no ProwJobs ID was found.
func (pjq *ProwJobsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pjq.Limit(1).IDs(setContextOp(ctx, pjq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{prowjobs.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pjq *ProwJobsQuery) FirstIDX(ctx context.Context) int {
	id, err := pjq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ProwJobs entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ProwJobs entity is found.
// Returns a *NotFoundError when no ProwJobs entities are found.
func (pjq *ProwJobsQuery) Only(ctx context.Context) (*ProwJobs, error) {
	nodes, err := pjq.Limit(2).All(setContextOp(ctx, pjq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{prowjobs.Label}
	default:
		return nil, &NotSingularError{prowjobs.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pjq *ProwJobsQuery) OnlyX(ctx context.Context) *ProwJobs {
	node, err := pjq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ProwJobs ID in the query.
// Returns a *NotSingularError when more than one ProwJobs ID is found.
// Returns a *NotFoundError when no entities are found.
func (pjq *ProwJobsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pjq.Limit(2).IDs(setContextOp(ctx, pjq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{prowjobs.Label}
	default:
		err = &NotSingularError{prowjobs.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pjq *ProwJobsQuery) OnlyIDX(ctx context.Context) int {
	id, err := pjq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ProwJobsSlice.
func (pjq *ProwJobsQuery) All(ctx context.Context) ([]*ProwJobs, error) {
	ctx = setContextOp(ctx, pjq.ctx, ent.OpQueryAll)
	if err := pjq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ProwJobs, *ProwJobsQuery]()
	return withInterceptors[[]*ProwJobs](ctx, pjq, qr, pjq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pjq *ProwJobsQuery) AllX(ctx context.Context) []*ProwJobs {
	nodes, err := pjq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ProwJobs IDs.
func (pjq *ProwJobsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if pjq.ctx.Unique == nil && pjq.path != nil {
		pjq.Unique(true)
	}
	ctx = setContextOp(ctx, pjq.ctx, ent.OpQueryIDs)
	if err = pjq.Select(prowjobs.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pjq *ProwJobsQuery) IDsX(ctx context.Context) []int {
	ids, err := pjq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pjq *ProwJobsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pjq.ctx, ent.OpQueryCount)
	if err := pjq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pjq, querierCount[*ProwJobsQuery](), pjq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pjq *ProwJobsQuery) CountX(ctx context.Context) int {
	count, err := pjq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pjq *ProwJobsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pjq.ctx, ent.OpQueryExist)
	switch _, err := pjq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("db: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pjq *ProwJobsQuery) ExistX(ctx context.Context) bool {
	exist, err := pjq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ProwJobsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pjq *ProwJobsQuery) Clone() *ProwJobsQuery {
	if pjq == nil {
		return nil
	}
	return &ProwJobsQuery{
		config:          pjq.config,
		ctx:             pjq.ctx.Clone(),
		order:           append([]prowjobs.OrderOption{}, pjq.order...),
		inters:          append([]Interceptor{}, pjq.inters...),
		predicates:      append([]predicate.ProwJobs{}, pjq.predicates...),
		withRepository:  pjq.withRepository.Clone(),
		withTektonTasks: pjq.withTektonTasks.Clone(),
		// clone intermediate query.
		sql:  pjq.sql.Clone(),
		path: pjq.path,
	}
}

// WithRepository tells the query-builder to eager-load the nodes that are connected to
// the "repository" edge. The optional arguments are used to configure the query builder of the edge.
func (pjq *ProwJobsQuery) WithRepository(opts ...func(*RepositoryQuery)) *ProwJobsQuery {
	query := (&RepositoryClient{config: pjq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pjq.withRepository = query
	return pjq
}

// WithTektonTasks tells the query-builder to eager-load the nodes that are connected to
// the "tekton_tasks" edge. The optional arguments are used to configure the query builder of the edge.
func (pjq *ProwJobsQuery) WithTektonTasks(opts ...func(*TektonTasksQuery)) *ProwJobsQuery {
	query := (&TektonTasksClient{config: pjq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pjq.withTektonTasks = query
	return pjq
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
//	client.ProwJobs.Query().
//		GroupBy(prowjobs.FieldJobID).
//		Aggregate(db.Count()).
//		Scan(ctx, &v)
func (pjq *ProwJobsQuery) GroupBy(field string, fields ...string) *ProwJobsGroupBy {
	pjq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ProwJobsGroupBy{build: pjq}
	grbuild.flds = &pjq.ctx.Fields
	grbuild.label = prowjobs.Label
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
//	client.ProwJobs.Query().
//		Select(prowjobs.FieldJobID).
//		Scan(ctx, &v)
func (pjq *ProwJobsQuery) Select(fields ...string) *ProwJobsSelect {
	pjq.ctx.Fields = append(pjq.ctx.Fields, fields...)
	sbuild := &ProwJobsSelect{ProwJobsQuery: pjq}
	sbuild.label = prowjobs.Label
	sbuild.flds, sbuild.scan = &pjq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ProwJobsSelect configured with the given aggregations.
func (pjq *ProwJobsQuery) Aggregate(fns ...AggregateFunc) *ProwJobsSelect {
	return pjq.Select().Aggregate(fns...)
}

func (pjq *ProwJobsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pjq.inters {
		if inter == nil {
			return fmt.Errorf("db: uninitialized interceptor (forgotten import db/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pjq); err != nil {
				return err
			}
		}
	}
	for _, f := range pjq.ctx.Fields {
		if !prowjobs.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("db: invalid field %q for query", f)}
		}
	}
	if pjq.path != nil {
		prev, err := pjq.path(ctx)
		if err != nil {
			return err
		}
		pjq.sql = prev
	}
	return nil
}

func (pjq *ProwJobsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ProwJobs, error) {
	var (
		nodes       = []*ProwJobs{}
		withFKs     = pjq.withFKs
		_spec       = pjq.querySpec()
		loadedTypes = [2]bool{
			pjq.withRepository != nil,
			pjq.withTektonTasks != nil,
		}
	)
	if pjq.withRepository != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, prowjobs.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ProwJobs).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ProwJobs{config: pjq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pjq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := pjq.withRepository; query != nil {
		if err := pjq.loadRepository(ctx, query, nodes, nil,
			func(n *ProwJobs, e *Repository) { n.Edges.Repository = e }); err != nil {
			return nil, err
		}
	}
	if query := pjq.withTektonTasks; query != nil {
		if err := pjq.loadTektonTasks(ctx, query, nodes,
			func(n *ProwJobs) { n.Edges.TektonTasks = []*TektonTasks{} },
			func(n *ProwJobs, e *TektonTasks) { n.Edges.TektonTasks = append(n.Edges.TektonTasks, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (pjq *ProwJobsQuery) loadRepository(ctx context.Context, query *RepositoryQuery, nodes []*ProwJobs, init func(*ProwJobs), assign func(*ProwJobs, *Repository)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ProwJobs)
	for i := range nodes {
		if nodes[i].repository_prow_jobs == nil {
			continue
		}
		fk := *nodes[i].repository_prow_jobs
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
			return fmt.Errorf(`unexpected foreign-key "repository_prow_jobs" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pjq *ProwJobsQuery) loadTektonTasks(ctx context.Context, query *TektonTasksQuery, nodes []*ProwJobs, init func(*ProwJobs), assign func(*ProwJobs, *TektonTasks)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*ProwJobs)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.TektonTasks(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(prowjobs.TektonTasksColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.prow_jobs_tekton_tasks
		if fk == nil {
			return fmt.Errorf(`foreign-key "prow_jobs_tekton_tasks" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "prow_jobs_tekton_tasks" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (pjq *ProwJobsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pjq.querySpec()
	_spec.Node.Columns = pjq.ctx.Fields
	if len(pjq.ctx.Fields) > 0 {
		_spec.Unique = pjq.ctx.Unique != nil && *pjq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pjq.driver, _spec)
}

func (pjq *ProwJobsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(prowjobs.Table, prowjobs.Columns, sqlgraph.NewFieldSpec(prowjobs.FieldID, field.TypeInt))
	_spec.From = pjq.sql
	if unique := pjq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pjq.path != nil {
		_spec.Unique = true
	}
	if fields := pjq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, prowjobs.FieldID)
		for i := range fields {
			if fields[i] != prowjobs.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pjq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pjq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pjq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pjq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pjq *ProwJobsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pjq.driver.Dialect())
	t1 := builder.Table(prowjobs.Table)
	columns := pjq.ctx.Fields
	if len(columns) == 0 {
		columns = prowjobs.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pjq.sql != nil {
		selector = pjq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pjq.ctx.Unique != nil && *pjq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range pjq.predicates {
		p(selector)
	}
	for _, p := range pjq.order {
		p(selector)
	}
	if offset := pjq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pjq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ProwJobsGroupBy is the group-by builder for ProwJobs entities.
type ProwJobsGroupBy struct {
	selector
	build *ProwJobsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pjgb *ProwJobsGroupBy) Aggregate(fns ...AggregateFunc) *ProwJobsGroupBy {
	pjgb.fns = append(pjgb.fns, fns...)
	return pjgb
}

// Scan applies the selector query and scans the result into the given value.
func (pjgb *ProwJobsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pjgb.build.ctx, ent.OpQueryGroupBy)
	if err := pjgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProwJobsQuery, *ProwJobsGroupBy](ctx, pjgb.build, pjgb, pjgb.build.inters, v)
}

func (pjgb *ProwJobsGroupBy) sqlScan(ctx context.Context, root *ProwJobsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pjgb.fns))
	for _, fn := range pjgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pjgb.flds)+len(pjgb.fns))
		for _, f := range *pjgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pjgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pjgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ProwJobsSelect is the builder for selecting fields of ProwJobs entities.
type ProwJobsSelect struct {
	*ProwJobsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pjs *ProwJobsSelect) Aggregate(fns ...AggregateFunc) *ProwJobsSelect {
	pjs.fns = append(pjs.fns, fns...)
	return pjs
}

// Scan applies the selector query and scans the result into the given value.
func (pjs *ProwJobsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pjs.ctx, ent.OpQuerySelect)
	if err := pjs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProwJobsQuery, *ProwJobsSelect](ctx, pjs.ProwJobsQuery, pjs, pjs.inters, v)
}

func (pjs *ProwJobsSelect) sqlScan(ctx context.Context, root *ProwJobsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pjs.fns))
	for _, fn := range pjs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pjs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pjs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
