package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Bugs holds the schema definition for the Bugs entity.
type Bugs struct {
	ent.Schema
}

// Fields of the Password.
func (Bugs) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),
		field.Text("jira_key").
			SchemaType(textSchema).
			NotEmpty().
			Unique(),
		field.Time("created_at").
			SchemaType(timeSchema),
		field.Time("updated_at").
			SchemaType(timeSchema),
		field.Time("resolved_at").
			SchemaType(timeSchema).
			Nillable(),
		field.Bool("resolved").
			Default(false),
		field.Text("priority").
			SchemaType(textSchema),
		field.Float("resolution_time").
			Default(0),
		field.Text("status").
			SchemaType(textSchema),
		field.Text("summary").
			SchemaType(textSchema),
		field.Text("url").
			SchemaType(textSchema),
		field.Text("project_key").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("assignment_time").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("prioritization_time").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("days_without_assignee").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("days_without_priority").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Float("days_without_resolution").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("labels").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("component").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("assignee").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("age").
			SchemaType(textSchema).
			Optional().
			Nillable(),
	}
}

// Edges of the Password.
func (Bugs) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bugs", Teams.Type).
			Ref("bugs").
			Unique(),
	}
}
