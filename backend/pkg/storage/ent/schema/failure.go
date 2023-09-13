package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Failures holds the schema definition for the Failures entity.
type Failure struct {
	ent.Schema
}

// Fields of the Failures.
func (Failure) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),
		field.Text("jira_key").
			SchemaType(textSchema).
			NotEmpty().
			Unique(),
		field.Text("jira_status").
			SchemaType(textSchema),
		field.Text("error_message").
			SchemaType(textSchema),
	}
}

// Edges of the Failures.
func (Failure) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("failures", Teams.Type).
			Ref("failures").
			Unique(),
	}
}
