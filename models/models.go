package models

import (
	"context"
	"encoding/json"
	"golearning/opensearchorm/orm"
)

type User struct {
	orm.BaseModel
	IDField string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

func (u *User) IndexName() string {
	return "users"
}

func (u *User) ID() string {
	return u.IDField
}

func (u *User) ResultToModel(data map[string]interface{}) (*User, error) {
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var model User
	// Map result data to model
	// (This would depend on the structure of your result and model)
	err = json.Unmarshal(jsonResult, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Class-level methods for querying
func (u *User) Filter(ctx context.Context, qb *orm.QueryBuilder) ([]*User, error) {
	results, err := u.BaseModel.Filter(ctx, u, qb)
	if err != nil {
		return nil, err
	}

	var models []*User
	// Parse the results into an array of models (you might need to adjust this based on how the results are returned)
	for _, result := range results {
		model, err := u.ResultToModel(result)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

func (u *User) First(ctx context.Context, qb *orm.QueryBuilder) (*User, error) {
	result, err := u.BaseModel.First(ctx, u, qb)
	if err != nil {
		return nil, err
	}
	return u.ResultToModel(result)
}

func (u *User) Count(ctx context.Context, qb *orm.QueryBuilder) (int, error) {
	return u.BaseModel.Count(ctx, u, qb)
}

func (u *User) Save(ctx context.Context) error {
	return u.BaseModel.Save(ctx, u)
}

func (u *User) Delete(ctx context.Context) error {
	return u.BaseModel.Delete(ctx, u)
}
