package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Users holds the schema definition for the Users entity.
type Users struct {
	ent.Schema
}

// Fields of the Users.
func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("user_id").
			Unique(),
		field.String("user_email").
			SchemaType(textSchema).
			Unique(),
		field.String("config").
			SchemaType(textSchema),
	}
}

// Edges of the Users.
func (Users) Edges() []ent.Edge {
	return []ent.Edge{}
}
