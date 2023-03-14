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
			Unique().
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
		field.Text("job_name").
			SchemaType(textSchema),
		field.Text("job_type").
			SchemaType(textSchema),
		field.Text("state").
			SchemaType(textSchema),
		field.Text("job_url").
			SchemaType(textSchema),
		field.Int16("ci_failed").
			SchemaType(intSchema),
		field.Text("e2e_failed_test_messages").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("suites_xml_url").
			SchemaType(textSchema).
			Optional().
			Nillable(),
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
