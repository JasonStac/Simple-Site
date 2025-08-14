package server

import (
	artistHandler "goserv/internal/domain/artists/handler"
	postHandler "goserv/internal/domain/posts/handler"
	sessionHandler "goserv/internal/domain/sessions/handler"
	tagHandler "goserv/internal/domain/tags/handler"
	userHandler "goserv/internal/domain/users/handler"
	"goserv/internal/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) initRoutes(
	artistHandler *artistHandler.ArtistHandler,
	tagHandler *tagHandler.TagHandler,
	postHandler *postHandler.PostHandler,
	userHandler *userHandler.UserHandler,
	sessionHandler *sessionHandler.SessionHandler) {

	authMiddleware := middleware.AuthRestrictMiddleware(s.session)
	checkMiddleware := middleware.AuthCheckMiddleware(s.session)

	s.router.With(checkMiddleware).Get("/",
		func(w http.ResponseWriter, r *http.Request) {
			isUser := false
			userID, ok := middleware.GetUserID(r)
			if ok && userID != -1 {
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
		r.Mount("/posts/", http.HandlerFunc(postHandler.ViewPost))
		r.Get("/tags", tagHandler.ListTags)
		r.Get("/artists", artistHandler.ListArtists)
	})

	s.router.With(authMiddleware).Route("/profile", func(r chi.Router) {
		r.Get("/", userHandler.Profile)
		r.Get("/create", postHandler.ViewAddPost)
		r.Post("/create", postHandler.AddPost)
		r.Get("/uploads", postHandler.ListUserPosts)
		r.Get("/favourites", postHandler.ListUserFavs)
	})

	s.router.Mount("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	s.router.Mount("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("content"))))

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})
}
