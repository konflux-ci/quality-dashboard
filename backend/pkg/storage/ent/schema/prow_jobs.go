package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProwJobs holds the schema definition for the ProwJobs entity.
type ProwJobs struct {
	ent.Schema
}

// Fields of the Password.
func (ProwJobs) Fields() []ent.Field {
	return []ent.Field{
		field.Text("job_id").
			SchemaType(textSchema),
		field.Time("created_at").
			SchemaType(timeSchema),
		field.Float("duration").
			SchemaType(textSchema),
		field.Int64("tests_count").
			SchemaType(textSchema),
		field.Int64("failed_count").
			SchemaType(textSchema),
		field.Int64("skipped_count").
			SchemaType(textSchema),
	}
}

// Edges of the Password.
func (ProwJobs) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("prow_jobs", Repository.Type).
			Ref("prow_jobs").
			Unique(),
	}
}
