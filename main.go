package main

import (
    "context"
    "log"
    "golearning/models"
)

func main() {
    ctx := context.Background()

    // Initialize the OpenSearch client
    if err := models.InitializeClient("http://localhost:9200"); err != nil {
        log.Fatalf("Failed to initialize OpenSearch client: %v", err)
    }

    // Build a query using QueryBuilder with facets
    qb := models.NewQueryBuilder().
        Filter("name", "John Doe").
        AddFacet("email_domains", "email.keyword")

    user := &models.User{}

    // Get facets
    facets, err := user.GetFacets(ctx, qb)
    if err != nil {
        log.Fatalf("Failed to get facets: %v", err)
    }

    log.Printf("Facets: %+v", facets)
}
