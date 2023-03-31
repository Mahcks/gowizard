package mariadb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	ctx    context.Context
	Client *mongo.Client
}

func New(gCtx context.Context, uri string) (*MongoDB, error) {
	client, err := mongo.Connect(gCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping to see if connection was successful
	err = client.Ping(gCtx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		Client: client,
	}, nil
}

func (m *MongoDB) Close(gCtx context.Context) {
	if m.Client != nil {
		m.Client.Disconnect(gCtx)
	}
}
