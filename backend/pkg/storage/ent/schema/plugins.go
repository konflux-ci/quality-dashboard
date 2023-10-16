package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Plugins holds the schema definition for the Plugins entity.
type Plugins struct {
	ent.Schema
}

// Fields of the Password.
func (Plugins) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id").
			Unique(),
		field.Text("name").
			SchemaType(textSchema).
			NotEmpty().
			Unique(),
		field.Text("category").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("logo").
			SchemaType(textSchema).
			NotEmpty(),
		field.Text("description").
			SchemaType(textSchema),
		field.Text("status").
			SchemaType(textSchema).
			NotEmpty(),
	}
}

// Edges of the Password.
func (Plugins) Edges() []ent.Edge {
	return []ent.Edge{}
}
