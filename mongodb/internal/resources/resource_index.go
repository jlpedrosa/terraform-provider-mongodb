package resources
import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceIndex defines the terraform resource for a mongodb collection
func ResourceIndex() *schema.Resource {

	return &schema.Resource{
		Create: createIndex,
		Read:   readIndex,
		Delete: deleteIndex,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"collection": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},			
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},			
			"key": {
				Type:	schema.TypeMap,
				Required: true,
				ForceNew: true,				
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					ValidateFunc: validateIndexDirection,
				  },			
			},
		},
	}
}

func createIndex(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)

	collID := d.Get("collection").(string)
	idxName := d.Get("name").(string)
	keys := d.Get("key").(map[string]interface{})

	
	databaseName, collectionName := SplitCollectionID(collID)
	coll := getCollectionFromID(client, collID)

	idxModel := mongo.IndexModel{
		Keys: keys,

		Options: &options.IndexOptions{
			Name: &idxName,
		},		
	}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)	
	idxNameRes, err := coll.Indexes().CreateOne(ctx, idxModel)

	if err != nil {
		return err
	}
	
	d.SetId(fmt.Sprintf("%s.%s.%s", databaseName, collectionName, idxNameRes))
	return readIndex(d, meta)
}


func readIndex(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)	
	idx, err := getIndexByID(client, d.Id())

	if err != nil {
		return err
	}

	idxMap:=*idx
	
	d.Set("collection", idxMap["ns"])
	d.Set("name", idxMap["name"])
	d.Set("key", idxMap["key"])
	return nil
}

func deleteIndex(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongo.Client)
	
	idxID := d.Id()
	coll := getCollectionFromID(client, idxID)
	_,_, idxName := SplitIndexID(idxID)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	res, err := coll.Indexes().DropOne(ctx, idxName)
	if err != nil {
		return fmt.Errorf("Unable to delete index %s, %s", idxID, err)
	}

	res.Index(0)
	d.SetId("")
	return nil
}

func getIndexByID(c *mongo.Client, idx string) (*bson.M, error)  {
	_, _, idxName := SplitIndexID(idx)
	coll := getCollectionFromID(c, idx)
	return getIndexByNameFromCollection(coll, idxName)	
}

//SplitIndexID returns the DbName, Collectionname and index name of a MongoBD Index
func SplitIndexID(id string) (string, string, string) {
	slice := strings.Split(id, ".")
	return slice[0],slice[1],slice[2]
}

func getIndexByNameFromCollection(c *mongo.Collection, idx string) (*bson.M, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)	
	idxs, err := c.Indexes().List(ctx)

	
	if err != nil {
		return nil, err
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)

	idxRes := make([]bson.D, 5)
	err = idxs.All(ctx, &idxRes) 

	if err != nil {
		return nil, err
	}

	for _, idxIf := range idxRes {		
		idxMap := idxIf.Map()
		idxName := idxMap["name"].(string)
		if idxName == idx {
			return &idxMap, nil
		}
	}	
	return nil, fmt.Errorf("Unable to find index: %s", idx)
}

func validateIndexDirection(i interface{}, n string) ([]string, []error) {

	return nil, nil
}