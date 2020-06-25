provider "mongodb" {
  url = "mongodb://localhost:27017"
}

resource "mongodb_collection" "sample" {
  database = "terradb"
  collection = "terracoll"
}

resource "mongodb_index" "sample" {
  collection = mongodb_collection.sample.id
  name = "docs"
  key = {
    "type" = 1
    "document" = -1
    "otherfield" = 1
  }
}


