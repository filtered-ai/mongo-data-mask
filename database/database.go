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

	start := time.Now()
	// Perperatory data generation pass
	// Complete data generation pass
	collections.PopulateUsers(db, ctx, base, 4)
	collections.PopulateSubscriptions(db, ctx, base, 4)
	collections.PopulateOrgs(db, ctx, base, 4)
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
