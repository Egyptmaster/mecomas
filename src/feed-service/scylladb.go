package main

import (
	"github.com/gocql/gocql"
)

func test() {
	conf := gocql.NewCluster("localhost:9041", "localhost:9042")

	session, err := conf.CreateSession()
	if err != nil {

	}

	session.Close()
}
