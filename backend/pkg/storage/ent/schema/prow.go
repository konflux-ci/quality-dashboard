package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
type Prow struct {
	ent.Schema
}

// Fields of the Password.
func (Prow) Fields() []ent.Field {
	return []ent.Field{
		field.String("job_id").
			SchemaType(textSchema),
		field.Text("Name").
			SchemaType(textSchema),
		field.Text("Status").
			SchemaType(textSchema),
		field.Float("time").
			SchemaType(textSchema),
	}
}

// Edges of the Password.
func (Prow) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge called "owner" of type `User`
		// and reference it to the "cars" edge (in User schema)
		// explicitly using the `Ref` method.
		edge.From("prow", Repository.Type).
			Ref("prow").
			Unique(),
	}
}
