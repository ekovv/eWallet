package mongodb

import (
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	conn *sql.DB
}

func NewRepository(database *mongo.Database) *MongoStorage {
	return &MongoStorage{
		conn: database,
	}
}
