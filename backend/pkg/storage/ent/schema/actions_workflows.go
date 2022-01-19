package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

/* Original SQL table:
create table password
(
    email    text not null  primary key,
    hash     blob not null,
    username text not null,
    user_id  text not null
);
*/

// Password holds the schema definition for the Password entity.
type Workflows struct {
	ent.Schema
}

// Fields of the Password.
func (Workflows) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("workflow_id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("workflow_id").
			Unique(),
		field.Text("workflow_name").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("badge_url").
			SchemaType(textSchema),
		field.Text("html_url").
			SchemaType(textSchema),
		field.Text("job_url").
			SchemaType(textSchema),
		field.Text("state").
			SchemaType(textSchema),
	}
}

// Edges of the Password.
func (Workflows) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workflows", Repository.Type).
			Ref("workflows").
			Unique(),
	}
}
