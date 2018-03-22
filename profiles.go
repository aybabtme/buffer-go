package buffer

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

var (
	profilesListRef   = mustURL("1/profiles.json")
	profilesGetRawRef = "1/profiles/%s.json"
)

type ProfilesService struct {
	client *Client
}

func (us *ProfilesService) List(ctx context.Context) ([]Profile, error) {
	list := make([]Profile, 0)
	return list, us.client.get(ctx, profilesListRef, &list)
}

func (us *ProfilesService) Get(ctx context.Context, id string) (*Profile, error) {
	ref, err := url.Parse(fmt.Sprintf(profilesGetRawRef, id))
	if err != nil {
		return nil, errors.Wrap(err, "preparing URL")
	}
	profile := new(Profile)
	return list, us.client.get(ctx, ref, &profile)
}

func (us *ProfilesService) Deauthorize(ctx context.Context) (bool, error) {
	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	return resp.Success, us.client.post(ctx, userDeauthorizeRef, nil, &resp)
}

type Profile struct {
	UnderscoreID         string `json:"_id"`
	Avatar               string `json:"avatar"`
	AvatarHTTPS          string `json:"avatar_https"`
	CanSeeContentLibrary bool   `json:"can_see_content_library"`
	Counts               struct {
		DailySuggestions int `json:"daily_suggestions"`
		Drafts           int `json:"drafts"`
		Reminders        int `json:"reminders"`
		Sent             int `json:"sent"`
		Pending          int `json:"pending"`
	} `json:"counts"`

	CreatedAt         UnixTime   `json:"created_at"`
	Default           bool       `json:"default"`
	Disabled          bool       `json:"disabled"`
	Disconnected      bool       `json:"disconnected"`
	FormattedService  string     `json:"formatted_service"`
	FormattedUsername string     `json:"formatted_username"`
	ID                string     `json:"id"`
	IsOnBusinessV2    bool       `json:"is_on_business_v2"`
	Locked            bool       `json:"locked"`
	Paused            bool       `json:"paused"`
	PausedSchedules   []Schedule `json:"paused_schedules"`
	Preferences       struct {
	} `json:"preferences"`
	Schedules       []Schedule `json:"schedules"`
	Service         string     `json:"service"`
	ServiceID       string     `json:"service_id"`
	ServiceType     string     `json:"service_type"`
	ServiceUsername string     `json:"service_username"`
	Shortener       struct {
		Domain string `json:"domain"`
	} `json:"shortener"`
	Statistics   *Statistics `json:"statistics"`
	Timezone     Timezone    `json:"timezone"`
	TimezoneCity string      `json:"timezone_city"`
	UserID       string      `json:"user_id"`
	UtmTracking  string      `json:"utm_tracking"`
	Verb         string      `json:"verb"`

	// Fields that show up in responses but are not documented/not stable.

	UnspecifiedCoverPhoto       interface{}   `json:"cover_photo"`
	UnspecifiedDisabledFeatures []interface{} `json:"disabled_features"`
	UnspecifiedReportsLogo      interface{}   `json:"reports_logo"`
}

type Schedule struct {
	Days  []Weekday `json:"days"`
	Times []Daytime `json:"times"`
}

type Statistics struct {
	Followers int `json:"followers"`
}
