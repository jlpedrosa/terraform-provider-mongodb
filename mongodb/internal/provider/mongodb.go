package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb/internal/resources"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Provider defines the terraform mongodb-provider 
func Provider() *schema.Provider {
	return &schema.Provider{

		ResourcesMap: map[string]*schema.Resource{
			"mongodb_collection": resources.ResourceCollection(),
			"mongodb_index" : resources.ResourceIndex(),
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Default:  "mongodb://localhost:27017",
				Optional: true,
			},
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	url := d.Get("url").(string)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to mongodb server did not response on time: %v", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, fmt.Errorf("Mongodb server did not response on time: %v", err)
	}

	d.SetId(url)
	return client, nil
}
