package resources

import (
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
)


// SplitCollectionID gets the Id of a collection in db.collection and returns the db and collection
func SplitCollectionID(id string) (string, string) {
	s := strings.Split(id,".")
	return s[0], s[1]
}

func toPrt(v bool) (*bool) {
	return &v
}

func getCollectionFromID(c *mongo.Client, id string) (*mongo.Collection) {

	dbName, collName :=  SplitCollectionID(id);
	coll := c.Database(dbName).Collection(collName)
	return coll
} 