package schema

import (
	"goserv/internal/static/enum"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.Enum("media_type").
			Values(enum.MediaType("").Values()...).
			SchemaType(map[string]string{
				dialect.Postgres: "media_type",
			}),
		field.String("filename").Unique(),
		field.String("file_ext"),
		field.Int("user_owns").Optional(),
	}
}

func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("owns").Unique().Field("user_owns"),
		edge.From("favourited_by", User.Type).Ref("favourites"),
		edge.To("tags", Tag.Type),
	}
}
