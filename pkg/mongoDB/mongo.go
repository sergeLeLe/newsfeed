package mongoDB

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	ErrConnectToDB  = errors.New("failed to connect to the database")
	ErrCreateClient = errors.New("failed to create a client, check the correctness of the data provided")
)

type DB struct {
	Client *mongo.Client
}

func attemptPingDB(db *mongo.Client, logger logrus.FieldLogger, nAttempts int) error {
	startDelayTime := time.Second * 2
	ctx := context.TODO()

	for nAttempts > 0 {
		err := db.Ping(ctx, nil)
		if err == nil {
			return nil
		}
		logger.Warningf("failed to connect to the database, retry via %v", nAttempts)
		time.Sleep(startDelayTime)

		nAttempts--
		startDelayTime *= 2
	}
	return ErrConnectToDB
}

func New(logger logrus.FieldLogger, host, port string, nAttempts int) (*DB, error) {
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v", host, port)),
	)
	if err != nil {
		return nil, ErrCreateClient
	}

	if err := attemptPingDB(client, logger, nAttempts); err != nil {
		return nil, fmt.Errorf("%w; Using host: %v:%v", ErrConnectToDB, host, port)
	}
	return &DB{Client: client}, nil
}

func (d *DB) Shutdown(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

