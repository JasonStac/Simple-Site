package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().Unique().Immutable(),
		field.String("pass_hash").NotEmpty(),
		field.Bool("is_admin").Default(false),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owns", Post.Type),
		edge.To("favourites", Post.Type),
		edge.To("session", Session.Type).Unique(),
	}
}
