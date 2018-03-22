package buffer

import (
	"context"
)

var userRef = mustURL("1/user.json")

type UserService struct {
	client *Client
}

type User struct {
	UnderscoreID string `json:"_id"`
	ActivityAt   int    `json:"activity_at"`
	CreatedAt    int    `json:"created_at"`
	ID           string `json:"id"`
	Plan         string `json:"plan"`
	Timezone     string `json:"timezone"`
}

func (us *UserService) Get(ctx context.Context) (*User, error) {
	user := new(User)
	return user, us.client.get(ctx, userRef, user)
}
