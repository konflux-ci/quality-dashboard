package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CodeCov holds the schema definition for the CodeCov entity.
type CodeCov struct {
	ent.Schema
}

// Fields of the Password.
func (CodeCov) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),
		field.Text("repository_name").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("git_organization").
			SchemaType(textSchema),
		field.Float("coverage_percentage").
			SchemaType(intSchema),
		field.Float("average_retests").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		// average_retests_to_merge deprecated
		// this info is now saved on each merged pr
		field.Float("average_retests_to_merge").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("coverage_trend").
			SchemaType(textSchema).
			Optional().
			Nillable(),
	}
}

// Edges of the Password.
func (CodeCov) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("codecov", Repository.Type).
			Ref("codecov").
			Unique(),
	}
}
