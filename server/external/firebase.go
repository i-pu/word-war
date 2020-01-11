package external

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("firebase init: %v", err)
	}

	FirebaseApp = app
}

func GetFirestore() *firestore.Client {
	ctx := context.Background()
	store, err := FirebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("firestore init: %v", err)
	}
	return store
}
