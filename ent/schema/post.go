package schema

import (
	"goserv/internal/models"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.Enum("media_type").Values(models.MediaType("").Values()...),
		field.String("filename").Unique(),
		field.Int("user_owns"),
	}
}

func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("owns").Unique().Field("user_owns").Required(),
		edge.From("favourited_by", User.Type).Ref("favourites"),
		edge.To("artists", Artist.Type),
		edge.To("tags", Tag.Type),
	}
}
