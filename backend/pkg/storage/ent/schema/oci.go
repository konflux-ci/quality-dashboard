package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// OCI holds the schema definition for the OCI entity.
// It represents an OCI artifact's metadata.
type OCI struct {
	ent.Schema
}

// Fields of the OCI.
func (OCI) Fields() []ent.Field {
	return []ent.Field{
		// The primary key, a UUID.
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		// A URL to access the artifact manifest or blob. Can be optional.
		field.String("artifact_url").
			Optional().
			Comment("A URL to access the artifact manifest or blob"),

		// Timestamps for creation and updates.
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("The time the artifact record was created"),

		field.Time("updated_at").
			Default(time.Now).
			Nillable().
			Comment("The time the artifact record was last updated"),
	}
}

// Edges of the OCI.
func (OCI) Edges() []ent.Edge {
	return []ent.Edge{
		// Defines a relationship to a parent Repository entity.
		// This creates a many-to-one relationship: many OCIs belong to one Repository.
		edge.From("oci", Repository.Type).
			Ref("oci").
			Unique(),
	}
}
