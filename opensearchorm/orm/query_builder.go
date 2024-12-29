package orm

// QueryBuilder helps dynamically create queries for OpenSearch
type QueryBuilder struct {
	Filters      map[string]interface{}
	Sorts        []map[string]interface{}
	Aggregations map[string]interface{}
}

// NewQueryBuilder initializes a new QueryBuilder instance
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		Filters:      make(map[string]interface{}),
		Sorts:        []map[string]interface{}{},
		Aggregations: make(map[string]interface{}),
	}
}

// Filter adds a filter condition to the query
func (qb *QueryBuilder) Filter(field string, value interface{}) *QueryBuilder {
	qb.Filters[field] = value
	return qb
}

// OrderBy adds an order condition to the query
func (qb *QueryBuilder) OrderBy(field string, direction string) *QueryBuilder {
	if direction != "asc" && direction != "desc" {
		direction = "asc" // Default to ascending if invalid direction is provided
	}
	qb.Sorts = append(qb.Sorts, map[string]interface{}{
		field: map[string]interface{}{
			"order": direction,
		},
	})
	return qb
}

// AddFacet adds an aggregation (facet) to the query
func (qb *QueryBuilder) AddFacet(name string, field string) *QueryBuilder {
	qb.Aggregations[name] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": field,
		},
	}
	return qb
}

// Build constructs the OpenSearch-compatible query
func (qb *QueryBuilder) Build() map[string]interface{} {
	mustClauses := []map[string]interface{}{}
	for field, value := range qb.Filters {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				field: value,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
	}

	if len(qb.Sorts) > 0 {
		query["sort"] = qb.Sorts
	}

	if len(qb.Aggregations) > 0 {
		query["aggs"] = qb.Aggregations
	}

	return query
}
