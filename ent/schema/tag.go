package schema

import (
	"goserv/internal/models"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.Enum("tag_type").
			Values(models.TagType("").Values()...).
			SchemaType(map[string]string{
				dialect.Postgres: "tag_type",
			}),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("posts", Post.Type).Ref("tags"),
	}
}
