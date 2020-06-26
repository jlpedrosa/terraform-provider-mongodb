package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

// ResourceCollection defines the terraform resource for a mongodb collection
func ResourceCollection() *schema.Resource {

	return &schema.Resource{
		Create: createCollection,
		Read:   readCollection,
		Delete: deleteCollection,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				
			},
			"collection": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createCollection(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)
	databaseName := d.Get("database").(string)

	db := client.Database(databaseName, nil)

	collectionName := d.Get("collection").(string)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	err := db.CreateCollection(ctx, collectionName)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s.%s", databaseName, collectionName))
	return readCollection(d, meta)
}


func readCollection(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)

	ids := strings.Split(d.Id(), ".")
	databaseName := ids[0]
	collectionName := ids[1]

	coll := client.Database(databaseName).Collection(collectionName)
	databaseName = coll.Database().Name()
	collectionName = coll.Name()
	d.Set("database", databaseName)
	d.Set("collection",collectionName)
	return nil
}

func deleteCollection(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)
	
	databaseName := d.Get("database").(string)
	database := client.Database(databaseName)

	collectionName := d.Get("collection").(string)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := database.Collection(collectionName).Drop(ctx)
	if err != nil {
		return fmt.Errorf("Unable to delete collection %s, %s", collectionName, err)
	}

	d.SetId("")
	return nil
}
