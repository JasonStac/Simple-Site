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
	artistHandler := s.initArtist()
	tagHandler := s.initTag()
	postHandler := s.initPost()
	userHandler, sessionHandler := s.initAuth()

	s.initRoutes(artistHandler, tagHandler, postHandler, userHandler, sessionHandler)
}

func (s *Server) initArtist() *artistHandler.ArtistHandler {
	repo := artistRepo.NewArtistRepository(s.ent)
	service := artistService.NewArtistService(repo)
	handler := artistHandler.NewArtistHandler(service, s.tmplCache)
	return handler
}

func (s *Server) initTag() *tagHandler.TagHandler {
	repo := tagRepo.NewTagRepository(s.ent)
	service := tagService.NewTagService(repo)
	handler := tagHandler.NewTagHandler(service, s.tmplCache)
	return handler
}

func (s *Server) initPost() *postHandler.PostHandler {
	repo := postRepo.NewPostRepository(s.ent)
	service := postService.NewPostService(repo)
	handler := postHandler.NewPostHandler(service, s.tmplCache)
	return handler
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
