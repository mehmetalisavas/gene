package tweet

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

// Tweet represents a post a user created
type TweetService interface {
	Create(ctx context.Context, req *models.Account) (res *models.Account, err error)

	Delete(ctx context.Context, req *models.Account) (res *models.Account, err error)

	One(ctx context.Context, req *models.Account) (res *models.Account, err error)
}
