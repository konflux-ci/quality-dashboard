package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Fields of the Configuration.
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),
		field.String("team_name").
			SchemaType(textSchema).
			Unique(),
		field.Text("jira_config").
			SchemaType(textSchema),
		field.Text("bug_slos_config").
			SchemaType(textSchema),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("configuration", Teams.Type).
			Ref("configuration").
			Unique(),
	}
}
