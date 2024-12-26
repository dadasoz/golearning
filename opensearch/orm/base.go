package orm

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "strings"

    "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type BaseModel struct {
    CreatedDate string `json:"created_date"`
    UpdatedDate string `json:"updated_date"`
    Version     int    `json:"version"`
}

type Model interface {
    IndexName() string // Each model must specify its OpenSearch index name
    ID() string        // Each model must provide a unique ID
}

// PrepareForSave updates the default fields before saving
func (b *BaseModel) PrepareForSave(isNew bool) {
    now := time.Now().Format(time.RFC3339)
    if isNew {
        b.CreatedDate = now
        b.Version = 1
    } else {
        b.Version++
    }
    b.UpdatedDate = now
}

// Save the current model instance to the database
func (b *BaseModel) Save(ctx context.Context, model Model) error {
    isNew := model.ID() == "" || !b.documentExists(ctx, model.IndexName(), model.ID())
    b.PrepareForSave(isNew)

    body, err := json.Marshal(model)
    if err != nil {
        return err
    }

    req := opensearchapi.IndexRequest{
        Index:      model.IndexName(),
        DocumentID: model.ID(),
        Body:       strings.NewReader(string(body)),
        Refresh:    "true",
    }

    res, err := req.Do(ctx, GetClient())
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("error saving document: %s", res.String())
    }

    return nil
}

// Delete the current model instance from the database
func (b *BaseModel) Delete(ctx context.Context, model Model) error {
    req := opensearchapi.DeleteRequest{
        Index:      model.IndexName(),
        DocumentID: model.ID(),
        Refresh:    "true",
    }

    res, err := req.Do(ctx, GetClient())
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("error deleting document: %s", res.String())
    }

    return nil
}

// Filter documents using QueryBuilder
func (b *BaseModel) Filter(ctx context.Context, model Model, qb *QueryBuilder) ([]map[string]interface{}, error) {
    queryBody, err := json.Marshal(qb.Build())
    if err != nil {
        return nil, err
    }

    req := opensearchapi.SearchRequest{
        Index: []string{model.IndexName()},
        Body:  strings.NewReader(string(queryBody)),
    }

    res, err := req.Do(ctx, GetClient())
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.IsError() {
        return nil, fmt.Errorf("error filtering documents: %s", res.String())
    }

    var searchResults map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
        return nil, err
    }

    hits := searchResults["hits"].(map[string]interface{})["hits"].([]interface{})
    results := make([]map[string]interface{}, len(hits))

    for i, hit := range hits {
        results[i] = hit.(map[string]interface{})["_source"].(map[string]interface{})
    }

    return results, nil
}

// First retrieves the first document using QueryBuilder
func (b *BaseModel) First(ctx context.Context, model Model, qb *QueryBuilder) (map[string]interface{}, error) {
    qb.Filters["size"] = 1
    results, err := b.Filter(ctx, model, qb)
    if err != nil {
        return nil, err
    }
    if len(results) == 0 {
        return nil, errors.New("no documents found")
    }
    return results[0], nil
}

// Count returns the number of documents matching the QueryBuilder
func (b *BaseModel) Count(ctx context.Context, model Model, qb *QueryBuilder) (int, error) {
    qb.Filters["size"] = 0
    queryBody, err := json.Marshal(qb.Build())
    if err != nil {
        return 0, err
    }

    req := opensearchapi.SearchRequest{
        Index: []string{model.IndexName()},
        Body:  strings.NewReader(string(queryBody)),
    }

    res, err := req.Do(ctx, GetClient())
    if err != nil {
        return 0, err
    }
    defer res.Body.Close()

    if res.IsError() {
        return 0, fmt.Errorf("error counting documents: %s", res.String())
    }

    var searchResults map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
        return 0, err
    }

    count := int(searchResults["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
    return count, nil
}
