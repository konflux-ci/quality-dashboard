package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TektonTasks holds the schema definition for the TektonTasks entity.
// Each TektonTask record corresponds to a single TaskRun from a PipelineRun.
type TektonTasks struct {
	ent.Schema
}

// Fields of the TektonTasks entity.
func (TektonTasks) Fields() []ent.Field {
	return []ent.Field{
		// Unique identifier for the task record
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),

		// Name of the Tekton TaskRun (e.g., "fetch-source", "build-image")
		field.String("task_name").
			SchemaType(textSchema),

		// Duration of the TaskRun in seconds (as string)
		field.String("duration_seconds").
			SchemaType(textSchema),

		// indicates when a tekton task was inserted.
		field.Time("created_at").
			Default(time.Now).
			Nillable().
			Comment("Indicates when a tekton task was inserted."),
		// Final status of the TaskRun (e.g., "Succeeded", "Failed")
		field.String("status").
			SchemaType(textSchema),
	}
}

// Edges of the TektonTasks entity.
func (TektonTasks) Edges() []ent.Edge {
	return []ent.Edge{
		// Reference to the ProwJob this task is associated with
		edge.From("tekton_tasks", ProwJobs.Type).
			Ref("tekton_tasks").
			Unique(),
	}
}
