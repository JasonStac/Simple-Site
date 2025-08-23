package server

import (
	postHandler "goserv/internal/domain/posts/handler"
	sessionHandler "goserv/internal/domain/sessions/handler"
	tagHandler "goserv/internal/domain/tags/handler"
	userHandler "goserv/internal/domain/users/handler"
	"goserv/internal/middleware"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *Server) initRoutes(
	tagHandler *tagHandler.TagHandler,
	postHandler *postHandler.PostHandler,
	userHandler *userHandler.UserHandler,
	sessionHandler *sessionHandler.SessionHandler) {

	authMiddleware := middleware.AuthRestrictMiddleware(s.session)
	checkMiddleware := middleware.AuthCheckMiddleware(s.session)
	deleteMiddleware := middleware.DeleteMiddleware(s.user, s.post)

	s.router.With(checkMiddleware).Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			isUser := false
			userID, ok := middleware.GetUserID(r)
			if ok && userID != 0 {
				isUser = true
			}

			err := s.tmplCache.ExecuteTemplate(w, "home.html", struct{ IsUser bool }{
				IsUser: isUser,
			})
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
		},
	)

	s.router.With(checkMiddleware).Get("/register", userHandler.DisplayRegister)
	s.router.With(checkMiddleware).Post("/register", userHandler.Register)
	s.router.With(checkMiddleware).Get("/login", sessionHandler.DisplayLogin)
	s.router.With(checkMiddleware).Post("/login", sessionHandler.Login)
	s.router.With(authMiddleware).Get("/logout", sessionHandler.DisplayLogout)
	s.router.With(authMiddleware).Post("/logout", sessionHandler.Logout)

	s.router.With(checkMiddleware).Route("/view", func(r chi.Router) {
		r.Get("/posts", postHandler.ListPosts)
		r.Mount("/posts/", routeSinglePosts(postHandler))
		r.Get("/tags", tagHandler.ListGeneralTags)
		r.Get("/people", tagHandler.ListPeopleTags)
	})

	s.router.With(authMiddleware).Route("/profile", func(r chi.Router) {
		r.Get("/", userHandler.Profile)
		r.Get("/create", postHandler.ViewAddPost)
		r.Post("/create", postHandler.AddPost)
		r.Get("/uploads", postHandler.ListUserPosts)
		//r.Mount("/uploads/", routeSingleUploads(postHandler))
		r.Get("/favourites", postHandler.ListUserFavs)
	})

	s.router.With(authMiddleware, deleteMiddleware).Post("/delete", postHandler.DeletePost)
	s.router.With(authMiddleware).Post("/favourite", postHandler.FavouritePost)
	s.router.With(authMiddleware).Post("/unfavourite", postHandler.UnfavouritePost)

	s.router.Mount("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	//s.router.Mount("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("content"))))
	s.router.Mount("/assets/images/", routeFileServe())

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})
}

func routeSinglePosts(postHandler *postHandler.PostHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postHandler.ViewPost(w, r)
		default:
			http.Error(w, "Unsupported status method", http.StatusMethodNotAllowed)
		}
	}
}

func routeFileServe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			filename, ok := strings.CutPrefix(r.URL.Path, "/assets/images/")
			if !ok {
				http.NotFound(w, r)
				return
			}

			if len(filename) < 74 {
				http.NotFound(w, r)
				return
			}

			title := filename[70:]
			w.Header().Set("Content-Disposition", "attachment; filename="+title)
			http.ServeFile(w, r, filepath.Join("content", filename))
		default:
			http.Error(w, "Unsupported status method", http.StatusMethodNotAllowed)
		}
	}
}

// func routeSingleUploads(postHandler *postHandler.PostHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodPost:
// 			postHandler.DeletePost(w, r)
// 		default:
// 			http.Error(w, "Unsupported status method", http.StatusMethodNotAllowed)
// 		}
// 	}
// }
