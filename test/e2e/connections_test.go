package e2e

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func TestNats(t *testing.T) {
	uri := "nats://localhost:4222"
	nc, err := nats.Connect(uri)
	if err != nil {
		t.Fatalf("could not connect to nats: %s", err)
	}

	t.Run("Can create Jetstream Client", func(t *testing.T) {
		js, err := jetstream.New(nc)
		if err != nil {
			t.Fatalf("could not create JetStream client: %s", err)
		}
		sif := js.ListStreams(context.Background())
		if err := sif.Err(); err != nil {
			t.Fatalf("could not complete request to get streams: %s", err)
		}
	})

	t.Run("Can publish and consume a message", func(t *testing.T) {
		js, err := jetstream.New(nc)
		if err != nil {
			t.Fatalf("could not create JetStream client: %s", err)
		}

		s, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
			Name:     "TEST_STREAM",
			Subjects: []string{"FOO.*"},
		})
		if err != nil {
			t.Fatalf("could not create stream: %s", err)
		}

		cons, err := s.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
			Durable:   "TestConsumerConsume",
			AckPolicy: jetstream.AckExplicitPolicy,
		})
		if err != nil {
			t.Fatalf("could not create consumer: %s", err)
		}

		if _, err = js.Publish(context.Background(), "FOO.TEST1", []byte("msg")); err != nil {
			fmt.Println("pub error: ", err)
		}

		var message jetstream.Msg
		var wg sync.WaitGroup
		wg.Add(1)
		cc, err := cons.Consume(func(msg jetstream.Msg) {
			message = msg
			if err2 := msg.Ack(); err != nil {
				t.Fatalf("could not ack message: %s", err2)
			}
			wg.Done()
		}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			fmt.Println(err)
		}))
		if err != nil {
			t.Fatalf("could not start consumer: %s", err)
		}
		defer cc.Stop()

		wg.Wait()

		require.NotNil(t, message)
		assert.Equal(t, "msg", string(message.Data()))
	})
}
