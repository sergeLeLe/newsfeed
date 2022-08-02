package ports

import (
	"context"
	"newsfeed/internal/domain/model"
)

type FeedStorage interface {
	GetFeed(ctx context.Context, limit int64, offset int) ([]*model.Post, error)
	AddPost(ctx context.Context, typePost model.PostType, authorId string, attach *model.PostInput) (*model.Post, error)
}
