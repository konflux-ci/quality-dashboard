package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// PullRequests holds the schema definition for the PullRequests entity.
type PullRequests struct {
	ent.Schema
}

// Fields of the PullRequests.
func (PullRequests) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("pr_id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.Text("repository_name").
			SchemaType(textSchema),
		field.Text("repository_organization").
			SchemaType(textSchema),
		field.Int("number").
			SchemaType(intSchema),
		field.Time("created_at").
			SchemaType(timeSchema),
		field.Time("closed_at").
			SchemaType(timeSchema),
		field.Time("merged_at").
			SchemaType(timeSchema),
		field.Text("state").
			SchemaType(textSchema),
		field.Text("author").
			SchemaType(textSchema),
		field.Text("title").
			SchemaType(textSchema),
		field.Text("merge_commit").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("retest_before_merge_count").
			SchemaType(textSchema).
			Optional().
			Nillable(),
	}
}

// Edges of the PullRequests.
func (PullRequests) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("prs", Repository.Type).
			Ref("prs").
			Unique(),
	}
}
