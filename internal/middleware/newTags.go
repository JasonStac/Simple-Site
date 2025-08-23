package middleware

import (
	"context"
	"encoding/json"
	"goserv/internal/domain/tags"
	"goserv/internal/domain/tags/repository"
	"goserv/internal/models"
	"net/http"
)

func AddNewTags(repo repository.Tag) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jsonGeneralTags := r.FormValue("tags")
			var generalTags []tags.Tag
			if jsonGeneralTags != "" {
				if err := json.Unmarshal([]byte(jsonGeneralTags), &generalTags); err != nil {
					http.Error(w, "Failed to read tags", http.StatusBadRequest)
					return
				}
			}

			jsonPeopleTags := r.FormValue("people")
			var peopleTags []tags.Tag
			if jsonPeopleTags != "" {
				if err := json.Unmarshal([]byte(jsonPeopleTags), &peopleTags); err != nil {
					http.Error(w, "Failed to read people", http.StatusBadRequest)
					return
				}
			}

			allTags := append(generalTags, peopleTags...)
			for i := range allTags {
				if allTags[i].ID == 0 {
					id, err := repo.AddTag(r.Context(), allTags[i].Name, models.TagType(allTags[i].Type))
					if err != nil {
						http.Error(w, "Failed to add tag", http.StatusInternalServerError)
						return
					}
					allTags[i].ID = id
				}
			}

			ctx := context.WithValue(r.Context(), tagKey, allTags)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
