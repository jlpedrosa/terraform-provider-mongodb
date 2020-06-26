package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb/internal/acceptance"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb/internal/resources"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAccRCollectionCanCreate(t *testing.T) {
	
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		//PreCheck: acceptance.EnsureMongoIsUp,
		Providers:  acceptance.GetResourceProvider(),
		Steps: []resource.TestStep {
			{
				Config: simpleCollection(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_collection.testdb", "database", "sampledb"),
					resource.TestCheckResourceAttr("mongodb_collection.testdb", "collection", "any"),
				),
			},
		},
		CheckDestroy: testCollectionDestroy,
	})
}

func simpleCollection(t *testing.T) string {
	return fmt.Sprintf(`
		resource "mongodb_collection" "testdb" {
			database = "sampledb"
			collection = "any"
		}
	`)
}


func testCollectionDestroy(s *terraform.State) error {
	
	client := acceptance.GetMongoDBClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodb_collection" {
			continue
		}
		dbName, collName := resources.SplitCollectionID(rs.Primary.ID)
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)	
		colls, err := client.Database(dbName).ListCollectionNames(ctx, bson.D{})

		if err != nil {
			return fmt.Errorf("Unable to check if the collection was destroyed, error listing collections, %s",err)
		}
		
		for _, coll := range colls {
			if coll == collName {
				return fmt.Errorf("The colleciton %s was not deleted", collName)
			}
		}
	}
	return nil
}