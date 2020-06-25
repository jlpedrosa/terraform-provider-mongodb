package acceptance

import (
	"context"
	"fmt"
	"time"
	"sync"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb/internal/provider"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBProvider is the provider instace for the acceptance tests
var MongoDBProvider *schema.Provider

var once sync.Once

// GetMongoDBClient returns the Client instance from the provider
func GetMongoDBClient() *mongo.Client {
	return MongoDBProvider.Meta().(*mongo.Client)
}

//GetResourceProvider generates a resource provider for accceptance
func GetResourceProvider() map[string]terraform.ResourceProvider {
	once.Do(func() {
		MongoDBProvider = provider.Provider()
	})	
	
	return map[string]terraform.ResourceProvider{
		"mongodb": MongoDBProvider,
	}
}

//EnsureMongoIsUp checks if mongodb for acceptance is up and running
func EnsureMongoIsUp() {
	client := GetMongoDBClient()
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Errorf("Can't run tests, mongodb is not up"))
	}
}
