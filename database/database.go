package database

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/JRagone/mongo-data-gen/collections"
	"github.com/JRagone/mongo-data-gen/generators"
	"github.com/MichaelTJones/pcg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PopulateDatabase(client *mongo.Client, ctx context.Context) {
	// Get database
	testDB := client.Database("testing")
	// Drop old collections
	collectionNames, _ := testDB.ListCollectionNames(ctx, bson.D{})
	for _, collectionName := range collectionNames {
		if err := testDB.Collection(collectionName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}
	// Object holding seeded random generator
	seed := uint64(64)
	base := generators.Base{
		Seed:  seed,
		Pcg32: pcg.NewPCG32().Seed(seed, seed),
	}
	rand.Seed(int64(base.Seed))

	start := time.Now()
	collections.PopulateOrgs(testDB, ctx, base, 4)
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
