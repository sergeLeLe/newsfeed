package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"newsfeed/internal/adapters/graphql"
	"newsfeed/internal/adapters/repository"
	"newsfeed/pkg/graph/generated"
)

type Handler struct {
	repo *repository.Repo
}

func NewHandler(repo *repository.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetHandler() *chi.Mux {
	rootHandler:= handler.NewDefaultServer(generated.NewExecutableSchema(graphql.NewRootResolvers(h.repo)))

	r := chi.NewRouter()

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", rootHandler)

	return r
}

