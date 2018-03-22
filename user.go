package buffer

import (
	"context"
)

var (
	userGetRef         = mustURL("1/user.json")
	userDeauthorizeRef = mustURL("1/user/deauthorize.json")
)

type UserService struct {
	client *Client
}

func (us *UserService) Get(ctx context.Context) (*User, error) {
	user := new(User)
	return user, us.client.get(ctx, userGetRef, user)
}

func (us *UserService) Deauthorize(ctx context.Context) (bool, error) {
	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	return resp.Success, us.client.post(ctx, userDeauthorizeRef, nil, &resp)
}

type User struct {
	UnderscoreID       string   `json:"_id"`
	ActivityAt         UnixTime `json:"activity_at"`
	CreatedAt          UnixTime `json:"created_at"`
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Plan               string   `json:"plan"`
	Timezone           Timezone `json:"timezone"`
	TwentyfourHourTime bool     `json:"twentyfour_hour_time"`
	WeekStartsMonday   bool     `json:"week_starts_monday"`

	// Fields that show up in responses but are not documented/not stable.

	UnspecifiedPabloPreferences []interface{} `json:"pablo_preferences"`
	UnspecifiedProfileGroups    []interface{} `json:"profile_groups"`
}
