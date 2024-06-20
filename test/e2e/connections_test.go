package e2e

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongoDb(t *testing.T) {
	uri := "mongodb://localhost:27017/?replicaSet=rs0&tls=false&connect=direct&retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		t.Fatalf("could not connect to mongo: %s", err)
	}

	db := client.Database("test")
	col := db.Collection("test")

	t.Run("Can insert record", func(t *testing.T) {
		if _, err := col.InsertOne(context.Background(), bson.M{
			"testKey": "testValue",
		}); err != nil {
			t.Fatalf("failed to insert record: %s", err)
		}
	})

	t.Run("Can insert record in transaction", func(t *testing.T) {
		sess, err := client.StartSession()
		if err != nil {
			t.Fatalf("failed to create session: %s", err)
		}

		if _, err := sess.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
			if _, err := col.InsertOne(context.Background(), bson.M{
				"testKey": "testValue2",
			}); err != nil {
				return nil, err
			}
			return nil, nil
		}); err != nil {
			t.Fatalf("failed to insert record in transaction: %s", err)
		}
	})
}
