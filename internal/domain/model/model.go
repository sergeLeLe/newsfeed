package model

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	AuthorID string     `json:"author_id" bson:"author_id"`
	Type     PostType   `json:"type" bson:"type"`
	Order    *int       `json:"order" bson:"order"`
	Attach   Attachment `json:"attach" bson:"attach"`
}

func (ImagePost) IsAttachment() {}
func (VideoPost) IsAttachment() {}
func (TextPost) IsAttachment()  {}

type User struct {
	ID   string  `json:"id" bson:"id"`
	Name *string `json:"name" bson:"name"`
}

type VideoPost struct {
	Location *string `json:"location" bson:"location"`
}

type TextPost struct {
	Text *string `json:"text" bson:"text"`
}

type ImagePost struct {
	Location *string `json:"location" bson:"location"`
}

func (p *Post) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var rawData bson.Raw
	err := bson.Unmarshal(data, &rawData)
	if err != nil {
		return err
	}

	type PostUnmarshal struct {
		AuthorID string      `bson:"author_id"`
		Type     PostType    `bson:"type"`
		Order    *int        `bson:"order"`
		Attach   interface{} `bson:"attach"`
	}

	pu := PostUnmarshal{}
	err = rawData.Unmarshal(&pu)
	if err != nil {
		return err
	}

	p.Type = pu.Type
	p.AuthorID = pu.AuthorID
	p.Order = pu.Order

	var attach struct {
		Attach bson.Raw
	}

	err = rawData.Unmarshal(&attach)
	if err != nil {
		return err
	}

	switch p.Type {
	case PostTypeText:
		attachStructure := TextPost{}
		err = attach.Attach.Unmarshal(&attachStructure)
		p.Attach = attachStructure
	case PostTypeImage:
		attachStructure := ImagePost{}
		err = attach.Attach.Unmarshal(&attachStructure)
		p.Attach = attachStructure
	case PostTypeVideo:
		attachStructure := VideoPost{}
		err = attach.Attach.Unmarshal(&attachStructure)
		p.Attach = attachStructure
	default:
		return errors.Errorf("Unknown game name %s", p.Type)
	}

	return err
}
