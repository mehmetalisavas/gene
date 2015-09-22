package main

import (
	"github.com/cihangir/gene/example/twitter/models"
	"golang.org/x/net/context"
)

type TweetService interface {
	Create(ctx context.Context, req *models.Account) (res *models.Account, err error)
	Delete(ctx context.Context, req *models.Account) (res *models.Account, err error)
	One(ctx context.Context, req *models.Account) (res *models.Account, err error)
}
