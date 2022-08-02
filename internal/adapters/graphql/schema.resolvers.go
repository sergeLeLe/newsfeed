package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"newsfeed/internal/domain/model"
	"newsfeed/pkg/graph/generated"
)

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, typeArg model.PostType, authorID string, attach *model.PostInput) (*model.Post, error) {
	post, err := r.repo.AddPost(ctx, typeArg, authorID, attach)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Feed is the resolver for the feed field.
func (r *queryResolver) Feed(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	posts, err := r.repo.GetFeed(ctx, int64(*limit), *offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
