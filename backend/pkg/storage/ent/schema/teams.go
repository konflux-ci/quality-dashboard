package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ProwSuites holds the schema definition for the ProwSuites entity.
type Teams struct {
	ent.Schema
}

// Fields of the Password.
func (Teams) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("team_id").
			Unique(),
		field.String("team_name").
			Unique(),
		field.String("description").
			Unique(),
		field.String("jira_keys"),
	}
}

// Edges of the Password.
func (Teams) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.To("repositories", Repository.Type),
		edge.To("bugs", Bugs.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("failures", Failure.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("plugins", Plugins.Type).
			Ref("teams"),
	}
}
