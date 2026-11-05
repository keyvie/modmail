package database

import (
	"errors"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Data struct {
	Conn     *mongo.Client
	Database *mongo.Database
	Coll     *Collections
}

type Collections struct {
	Config *mongo.Collection
}

func NewDatabase(uri string) (*Data, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	name := strings.TrimPrefix(parsed.Path, "/")
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("database name is missing in the URI")
	}

	conn, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	database := conn.Database(name)

	coll := &Collections{
		Config: database.Collection("config"),
	}

	return &Data{
		Conn:     conn,
		Database: database,
		Coll:     coll,
	}, nil
}
