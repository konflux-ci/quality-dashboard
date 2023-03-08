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
