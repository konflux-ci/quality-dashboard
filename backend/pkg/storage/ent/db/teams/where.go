// Code generated by ent, DO NOT EDIT.

package teams

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Teams {
	return predicate.Teams(sql.FieldLTE(FieldID, id))
}

// TeamName applies equality check predicate on the "team_name" field. It's identical to TeamNameEQ.
func TeamName(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldTeamName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldDescription, v))
}

// JiraKeys applies equality check predicate on the "jira_keys" field. It's identical to JiraKeysEQ.
func JiraKeys(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldJiraKeys, v))
}

// TeamNameEQ applies the EQ predicate on the "team_name" field.
func TeamNameEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldTeamName, v))
}

// TeamNameNEQ applies the NEQ predicate on the "team_name" field.
func TeamNameNEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldNEQ(FieldTeamName, v))
}

// TeamNameIn applies the In predicate on the "team_name" field.
func TeamNameIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldIn(FieldTeamName, vs...))
}

// TeamNameNotIn applies the NotIn predicate on the "team_name" field.
func TeamNameNotIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldNotIn(FieldTeamName, vs...))
}

// TeamNameGT applies the GT predicate on the "team_name" field.
func TeamNameGT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGT(FieldTeamName, v))
}

// TeamNameGTE applies the GTE predicate on the "team_name" field.
func TeamNameGTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGTE(FieldTeamName, v))
}

// TeamNameLT applies the LT predicate on the "team_name" field.
func TeamNameLT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLT(FieldTeamName, v))
}

// TeamNameLTE applies the LTE predicate on the "team_name" field.
func TeamNameLTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLTE(FieldTeamName, v))
}

// TeamNameContains applies the Contains predicate on the "team_name" field.
func TeamNameContains(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContains(FieldTeamName, v))
}

// TeamNameHasPrefix applies the HasPrefix predicate on the "team_name" field.
func TeamNameHasPrefix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasPrefix(FieldTeamName, v))
}

// TeamNameHasSuffix applies the HasSuffix predicate on the "team_name" field.
func TeamNameHasSuffix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasSuffix(FieldTeamName, v))
}

// TeamNameEqualFold applies the EqualFold predicate on the "team_name" field.
func TeamNameEqualFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEqualFold(FieldTeamName, v))
}

// TeamNameContainsFold applies the ContainsFold predicate on the "team_name" field.
func TeamNameContainsFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContainsFold(FieldTeamName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContainsFold(FieldDescription, v))
}

// JiraKeysEQ applies the EQ predicate on the "jira_keys" field.
func JiraKeysEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEQ(FieldJiraKeys, v))
}

// JiraKeysNEQ applies the NEQ predicate on the "jira_keys" field.
func JiraKeysNEQ(v string) predicate.Teams {
	return predicate.Teams(sql.FieldNEQ(FieldJiraKeys, v))
}

// JiraKeysIn applies the In predicate on the "jira_keys" field.
func JiraKeysIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldIn(FieldJiraKeys, vs...))
}

// JiraKeysNotIn applies the NotIn predicate on the "jira_keys" field.
func JiraKeysNotIn(vs ...string) predicate.Teams {
	return predicate.Teams(sql.FieldNotIn(FieldJiraKeys, vs...))
}

// JiraKeysGT applies the GT predicate on the "jira_keys" field.
func JiraKeysGT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGT(FieldJiraKeys, v))
}

// JiraKeysGTE applies the GTE predicate on the "jira_keys" field.
func JiraKeysGTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldGTE(FieldJiraKeys, v))
}

// JiraKeysLT applies the LT predicate on the "jira_keys" field.
func JiraKeysLT(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLT(FieldJiraKeys, v))
}

// JiraKeysLTE applies the LTE predicate on the "jira_keys" field.
func JiraKeysLTE(v string) predicate.Teams {
	return predicate.Teams(sql.FieldLTE(FieldJiraKeys, v))
}

// JiraKeysContains applies the Contains predicate on the "jira_keys" field.
func JiraKeysContains(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContains(FieldJiraKeys, v))
}

// JiraKeysHasPrefix applies the HasPrefix predicate on the "jira_keys" field.
func JiraKeysHasPrefix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasPrefix(FieldJiraKeys, v))
}

// JiraKeysHasSuffix applies the HasSuffix predicate on the "jira_keys" field.
func JiraKeysHasSuffix(v string) predicate.Teams {
	return predicate.Teams(sql.FieldHasSuffix(FieldJiraKeys, v))
}

// JiraKeysEqualFold applies the EqualFold predicate on the "jira_keys" field.
func JiraKeysEqualFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldEqualFold(FieldJiraKeys, v))
}

// JiraKeysContainsFold applies the ContainsFold predicate on the "jira_keys" field.
func JiraKeysContainsFold(v string) predicate.Teams {
	return predicate.Teams(sql.FieldContainsFold(FieldJiraKeys, v))
}

// HasRepositories applies the HasEdge predicate on the "repositories" edge.
func HasRepositories() predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RepositoriesTable, RepositoriesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRepositoriesWith applies the HasEdge predicate on the "repositories" edge with a given conditions (other predicates).
func HasRepositoriesWith(preds ...predicate.Repository) predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RepositoriesInverseTable, RepositoryFieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RepositoriesTable, RepositoriesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBugs applies the HasEdge predicate on the "bugs" edge.
func HasBugs() predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BugsTable, BugsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBugsWith applies the HasEdge predicate on the "bugs" edge with a given conditions (other predicates).
func HasBugsWith(preds ...predicate.Bugs) predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BugsInverseTable, BugsFieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BugsTable, BugsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Teams) predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Teams) predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Teams) predicate.Teams {
	return predicate.Teams(func(s *sql.Selector) {
		p(s.Not())
	})
}
