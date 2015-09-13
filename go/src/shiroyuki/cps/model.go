package cps

import "time"
import "gopkg.in/mgo.v2/bson"

type Metadata struct {
    Url  string
    Type string
    Size uint64
}

type DBMetadata struct {
    Metadata // extended from Metadata
    Id        bson.ObjectId `bson:"_id"`
    CreatedAt time.Time     `bson:"created_at"`
}
