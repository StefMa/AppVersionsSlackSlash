package database

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"os"
)

type Database struct {
	client  *firestore.Client
	context context.Context
}

func CreateDatabase() (Database, error) {
	ctx := context.Background()
	serviceAccount := []byte(os.Getenv("FIREBASE_SERVICE_ACCOUNT"))
	sa := option.WithCredentialsJSON(serviceAccount)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return Database{}, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return Database{}, err
	}

	return Database{
		client:  client,
		context: ctx, // TODO: Context should not be saved in structs...
	}, nil
}

func (db Database) Close() {
	db.client.Close()
}
