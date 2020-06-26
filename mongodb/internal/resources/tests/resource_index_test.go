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

func TestAccRIndexCanCreate(t *testing.T) {
	
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  acceptance.GetResourceProvider(),
		Steps: []resource.TestStep {
			{
				Config: simpleIndex(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_index.testindex", "name", "sample_idx"),
					resource.TestCheckResourceAttr("mongodb_index.testindex", "collection", "indextstdb.indextstcoll"),
					resource.TestCheckResourceAttr("mongodb_index.testindex", "id", "indextstdb.indextstcoll.sample_idx"),
					resource.TestCheckResourceAttr("mongodb_index.testindex", "key.region", "1"),
					resource.TestCheckResourceAttr("mongodb_index.testindex", "key.state", "-1"),
					resource.TestCheckResourceAttr("mongodb_index.testindex", "key.age", "1"),					
				),
			},
		},
		CheckDestroy: testIndexDestroy,
	})
}

func simpleIndex(t *testing.T) string {
	return fmt.Sprintf(`
		resource "mongodb_collection" "testdb" {
			database = "indextstdb"
			collection = "indextstcoll"
		}

		resource "mongodb_index" "testindex" {
			collection = mongodb_collection.testdb.id
			name = "sample_idx"
			key = {
				region = 1
				state = -1
				age = 1
			}
		}
	`)
}

func testIndexDestroy(s *terraform.State) error {
	
	client := acceptance.GetMongoDBClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodb_index" {
			continue
		}
		dbName, collName, idxName := resources.SplitIndexID(rs.Primary.ID)
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

		idxs, err := client.Database(dbName).Collection(collName).Indexes().List(ctx)
		if err != nil {
		 	return fmt.Errorf("Unable to check if the index was destroyed, error listing collections, %s",err)
		}		
		
		idxRes := make([]bson.D, 5)
		err = idxs.All(ctx, &idxRes) 

		if err != nil {
			return err
		}
		
		for _, idxIf := range idxRes {		
			idxMap := idxIf.Map()
			idxInDb := idxMap["name"].(string)
			if idxName == idxInDb {
				return fmt.Errorf("Index %s was not deleted correctly", idxName)
			}
		}
	}
	return nil
}