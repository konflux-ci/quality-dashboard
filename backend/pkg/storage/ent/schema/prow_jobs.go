package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
			Optional().
			Nillable().
			SchemaType(textSchema),
		field.Int64("failed_count").
			Optional().
			Nillable().
			SchemaType(textSchema),
		field.Int64("skipped_count").
			Optional().
			Nillable().
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
		field.Bool("external_services_impact").
			Optional().
			Default(false).
			Nillable(),
		field.Text("e2e_failed_test_messages").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("suites_xml_url").
			SchemaType(textSchema).
			Optional().
			Nillable(),
		field.Text("build_error_logs").
			SchemaType(textSchema).
			Optional().
			Nillable(),
	}
}

// Edges of the Password.
func (ProwJobs) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("prow_jobs").
			Unique(),
		edge.To("tekton_tasks", TektonTasks.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
