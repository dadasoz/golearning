package models

import (
	"context"
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

// Class-level methods for querying
func (u *User) Filter(ctx context.Context, qb *orm.QueryBuilder) ([]map[string]interface{}, error) {
	return u.BaseModel.Filter(ctx, u, qb)
}

func (u *User) First(ctx context.Context, qb *orm.QueryBuilder) (map[string]interface{}, error) {
	return u.BaseModel.First(ctx, u, qb)
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
