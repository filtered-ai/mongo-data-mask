package database

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/JRagone/mongo-data-gen/collections"
	"github.com/JRagone/mongo-data-gen/generators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func PopulateDatabase(client *mongo.Client, ctx context.Context) {
	// Get database
	db := client.Database("testing")
	// Drop old collections
	collectionNames, _ := db.ListCollectionNames(ctx, bson.M{})
	for _, collectionName := range collectionNames {
		if err := db.Collection(collectionName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}
	// Object holding seeded random generator
	seed := uint64(64)
	base := generators.Base{
		Seed: seed,
	}
	rand.Seed(int64(base.Seed))

	start := time.Now()
	orgCount := uint(1000)
	userCount := uint(4)
	subCount := uint(4)
	preOrgs := make(map[int32]collections.Organization)
	users := make(map[int32]collections.User)
	subscriptions := make(map[primitive.ObjectID]collections.Subscription)
	// Preparatory data generation pass
	collections.PrepopulateOrgs(preOrgs, db, ctx, orgCount)
	collections.PrepopulateUsers(users, db, ctx, userCount)
	collections.PrepopulateSubscriptions(subscriptions, db, ctx, subCount)
	// Complete data generation pass
	collections.PopulateUsers(users, db, ctx, userCount)
	collections.PopulateSubscriptions(subscriptions, db, ctx, subCount)
	collections.PopulateOrgs(preOrgs, users, subscriptions, db, ctx, orgCount)
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
