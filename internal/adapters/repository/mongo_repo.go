package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"newsfeed/internal/config"
	"newsfeed/internal/domain/model"
	"newsfeed/pkg/mongoDB"
)

type Repo struct {
	*mongoDB.DB
	collection *mongo.Collection
}

func New(cfg *config.Config, db *mongoDB.DB) (*Repo, error) {
	postCol := db.Client.Database(cfg.Database.Name).Collection(postsCollection)
	return &Repo{db, postCol}, nil
}

func (r *Repo) GetFeed(ctx context.Context, limit int64, offset int) ([]*model.Post, error) {
	posts := make([]*model.Post, 0, limit)

	filter := bson.M{
		"order": bson.M{
			"$gte": offset,
		},
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"order", 1}})
	opts.SetLimit(limit)

	cur, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
func (r *Repo) AddPost(ctx context.Context, typePost model.PostType, authorId string, attach *model.PostInput) (*model.Post, error) {
	countPosts, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, nil
	}

	order := int(countPosts)

	var attachment model.Attachment
	switch typePost {
	case model.PostTypeText:
		attachment = model.TextPost{
			Text: attach.TextPost.Text,
		}
	case model.PostTypeImage:
		attachment = model.ImagePost{
			Location: attach.ImagePost.Location,
		}
	case model.PostTypeVideo:
		attachment = model.VideoPost{
			Location: attach.VideoPost.Location,
		}
	default:
		return nil, errors.Errorf("Unknown type post %s", typePost)
	}

	post := model.Post{
		Order: &order,
		AuthorID: authorId,
		Type: typePost,
		Attach: attachment,
	}

	_, err = r.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}