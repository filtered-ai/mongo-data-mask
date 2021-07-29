package database

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/JRagone/mongo-data-gen/collections"
	"github.com/JRagone/mongo-data-gen/generators"
	"go.mongodb.org/mongo-driver/bson"
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

	// This and many functions can be refactored when generics are introduced
	// to Go, hopefully by or in 2022
	start := time.Now()
	var generated collections.Collections
	generated.Organizations.Count = int32(1000)
	generated.Organizations.Data = make(collections.OrganizationData)
	generated.Users.Count = int32(400)
	generated.Users.Data = make(collections.UserData)
	generated.Subscriptions.Count = int32(400)
	generated.Subscriptions.Data = make(collections.SubscriptionData)
	// Preparatory data generation pass
	collections.PrepopulateOrgs(generated, db, ctx)
	collections.PrepopulateUsers(generated, db, ctx)
	collections.PrepopulateSubscriptions(generated, db, ctx)
	// Complete data generation pass
	collections.PopulateUsers(generated, db, ctx)
	collections.PopulateSubscriptions(generated, db, ctx)
	collections.PopulateOrgs(generated, db, ctx)
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
