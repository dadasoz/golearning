package main

import (
	"context"
	"fmt"
	"golearning/models"
	"golearning/opensearchorm/orm"
	"log"
)

func main() {
	ctx := context.Background()

	// Initialize the OpenSearch client
	if err := orm.InitializeClient("http://localhost:9200"); err != nil {
		log.Fatalf("Failed to initialize OpenSearch client: %v", err)
	}

	user := &models.User{}
	user.Name = "dadaso"
	user.IDField = "1244"
	user.Email = "test"
	err := user.Save(ctx)
	if err != nil {
		fmt.Print(err)
	}

	qb := orm.NewQueryBuilder()
	qb.Filter("name", "dadaso")

	users, err := user.Filter(ctx, qb)
	if err != nil {
		fmt.Print(err)
	}

	for _, u := range users {
		fmt.Println(u)
	}

}
