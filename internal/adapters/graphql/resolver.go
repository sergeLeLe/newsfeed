package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"newsfeed/internal/adapters/repository"
	"newsfeed/internal/domain/model"
	"newsfeed/pkg/graph/generated"
)

type Resolver struct{
	repo *repository.Repo
}

func NewRootResolvers(repo *repository.Repo) generated.Config {
	c := generated.Config{
		Resolvers: &Resolver{
			repo: repo,
		},
	}

	c.Directives.OneOf = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		input := obj.(map[string]interface{})

		if len(input) > 1{
			return nil, errors.New("One only input")
		}

		return next(ctx)
	}

	c.Directives.ValidatePostType = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		input := obj.(map[string]interface{})

		var attachType string
		typePost := input["type"].(string)
		attach := input["attach"].(map[string]interface{})
		for key, _ := range attach {
			attachType = key
			break
		}

		if typePost == model.PostTypeText.String() && attachType != "textPost" {
			return nil, errors.New("Invalid args")
		} else if typePost == model.PostTypeImage.String() && attachType != "imagePost" {
			return nil, errors.New("Invalid args")
		} else if typePost == model.PostTypeVideo.String() && attachType != "videoPost" {
			return nil, errors.New("Invalid args")
		}

		return next(ctx)
	}
	return c
}
