package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProwSuites holds the schema definition for the ProwSuites entity.
type ProwSuites struct {
	ent.Schema
}

// Fields of the Password.
func (ProwSuites) Fields() []ent.Field {
	return []ent.Field{
		field.String("job_id").
			SchemaType(textSchema),
		field.Text("name").
			SchemaType(textSchema),
		field.Text("status").
			SchemaType(textSchema),
		field.Float("time").
			SchemaType(textSchema),
	}
}

// Edges of the Password.
func (ProwSuites) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("prow_suites", Repository.Type).
			Ref("prow_suites").
			Unique(),
	}
}
