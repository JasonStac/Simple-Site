package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Session struct {
	ent.Schema
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.Int("user_id").Immutable(),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("sessions").Unique().Field("user_id").Required().Immutable(),
	}
}

func (Session) ID() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
	}
}
