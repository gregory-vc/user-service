package main

import (
	"os"

	"gopkg.in/couchbase/gocb.v1"
)

var couchbaseBucket string

func CreateGlobalBucket() (*gocb.Bucket, error) {

	couchbaseUri := os.Getenv("COUCHBASE")
	couchbaseBucket = os.Getenv("COUCHBASE_BUCKET")
	couchbaseUser := os.Getenv("COUCHBASE_USER")
	couchbasePassword := os.Getenv("COUCHBASE_PASSWORD")

	cluster, err := gocb.Connect(couchbaseUri)
	if err != nil {
		return nil, err
	}

	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: couchbaseUser,
		Password: couchbasePassword,
	})

	return cluster.OpenBucket(couchbaseBucket, "")
}
