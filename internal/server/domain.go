package server

import (
	artistHandler "goserv/internal/domain/artists/handler"
	artistRepo "goserv/internal/domain/artists/repository"
	artistService "goserv/internal/domain/artists/service"
	postHandler "goserv/internal/domain/posts/handler"
	postRepo "goserv/internal/domain/posts/repository"
	postService "goserv/internal/domain/posts/service"
	sessionHandler "goserv/internal/domain/sessions/handler"
	sessionRepo "goserv/internal/domain/sessions/repository"
	sessionService "goserv/internal/domain/sessions/service"
	tagHandler "goserv/internal/domain/tags/handler"
	tagRepo "goserv/internal/domain/tags/repository"
	tagService "goserv/internal/domain/tags/service"
	userHandler "goserv/internal/domain/users/handler"
	userRepo "goserv/internal/domain/users/repository"
	userService "goserv/internal/domain/users/service"
)

func (s *Server) initDomain() {
	postHandler, tagHandler, artistHandler := s.initContent()
	userHandler, sessionHandler := s.initAuth()

	s.initRoutes(artistHandler, tagHandler, postHandler, userHandler, sessionHandler)
}

func (s *Server) initContent() (*postHandler.PostHandler, *tagHandler.TagHandler, *artistHandler.ArtistHandler) {
	tRepo := tagRepo.NewTagRepository(s.ent)
	tService := tagService.NewTagService(tRepo)
	tHandler := tagHandler.NewTagHandler(tService, s.tmplCache)

	aRepo := artistRepo.NewArtistRepository(s.ent)
	aService := artistService.NewArtistService(aRepo)
	aHandler := artistHandler.NewArtistHandler(aService, s.tmplCache)

	pRepo := postRepo.NewPostRepository(s.ent)
	pService := postService.NewPostService(pRepo)
	pHandler := postHandler.NewPostHandler(pService, tService, aService, s.tmplCache)

	return pHandler, tHandler, aHandler
}

func (s *Server) initAuth() (*userHandler.UserHandler, *sessionHandler.SessionHandler) {
	userRepo := userRepo.NewUserRepository(s.ent)
	userService := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userService, s.tmplCache)

	sessionRepo := sessionRepo.NewSessionRepository(s.ent)
	sessionService := sessionService.NewSessionService(sessionRepo, userService)
	sessionHandler := sessionHandler.NewSessionHandler(sessionService, s.tmplCache)
	s.session = sessionRepo

	return userHandler, sessionHandler
}
