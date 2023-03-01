package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Password holds the schema definition for the Password entity.
type Repository struct {
	ent.Schema
}

// Fields of the Password.
func (Repository) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("repo_id").
			Unique(),
		field.Text("repository_name").
			SchemaType(textSchema).
			NotEmpty().
			Unique(),
		field.Text("git_organization").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("description").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("git_url").
			SchemaType(textSchema).
			NotEmpty(),
	}
}

// Edges of the Password.
func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repositories", Teams.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			Ref("repositories").
			Unique(),
		edge.To("workflows", Workflows.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("codecov", CodeCov.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("prow_suites", ProwSuites.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("prow_jobs", ProwJobs.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
