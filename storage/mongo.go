package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/ympyst/uvindex-tgbot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type Mongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongo() (*Mongo, error) {
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Mongo{
		client:     client,
		collection: client.Database("uvindex").Collection("usersState"),
	}, nil
}

func (s *Mongo) GetUserStateOrCreate(ctx context.Context, userId int64) (model.UserState, error) {
	var u model.UserState
	err := s.collection.FindOne(ctx, bson.D{{"id", userId}}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		u = model.UserState{
			Id:          userId,
			UVThreshold: model.DefaultUVThreshold,
			State:       model.Start,
		}
		_, err = s.collection.InsertOne(ctx, u)
		if err != nil {
			return model.UserState{}, err
		}
	} else if err != nil {
		return model.UserState{}, err
	}

	return u, nil
}

func (s *Mongo) SaveState(ctx context.Context, state *model.UserState) error {
	update := bson.M{
		"$set": bson.M{
			"is_group":     state.IsGroup,
			"username":     state.Username,
			"state":        state.State,
			"location":     state.Location,
			"uv_threshold": state.UVThreshold,
			//todo alert_schedule
		},
	}
	res, err := s.collection.UpdateOne(ctx, bson.D{{"id", state.Id}}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New(fmt.Sprintf("user %d not found", state.Id))
	}
	return nil
}
