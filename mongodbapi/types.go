package mongodbapi

import (
	"time"
)

//
// User data
//

// UserData -
type UserData struct {
	ID                 int       `bson:"_id"`
	Locale             string    `bson:"locale"`
	PremiumExpiration  time.Time `bson:"premium_expiration"`
	HasPremiumSub      bool      `bson:"has_premium_sub"`
	ExcessPremiumMin   int       `bson:"excess_premium_min"`
	VerifiedExpiration time.Time `bson:"verified_expiration"`
	VerifiedID         int       `bson:"verified_id"`
	DefaultPID         int       `bson:"default_player_id"`
	AccessToken        string    `bson:"access_token"`
	CustomBgURL        string    `bson:"custom_bg"`

	ShadowBanned    bool   `json:"shadow_banned"`
	ShadowBanReason string `json:"shadow_ban_reason"`
}

// BanData -
type BanData struct {
	UserID     int       `bson:"user_id" json:"user_id"`
	Reason     string    `bson:"reason" json:"reason"`
	Notified   bool      `bson:"notified" json:"notified"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Expiration time.Time `bson:"expiration" json:"expiration"`
}

// UserDataIntent - Intent to edit User data
type UserDataIntent struct {
	IntentID  string    `bson:"_id"`
	Timestamp time.Time `bson:"timestamp"`
	Data      UserData  `bson:"data"`
	Realm     string    `bson:"realm"`
}

//
// Logins
//

// LoginIntent -
type LoginIntent struct {
	IntentID  string    `bson:"_id" json:"-"`
	Timestamp time.Time `bson:"timestamp" json:"-"`
	LoginData
}

// LoginData -
type LoginData struct {
	DiscordID int    `bson:"discord_user_id" json:"discord_user_id"`
	Realm     string `bson:"realm" json:"realm"`
}

//
// Payments
//

// PayPalPaymentIntentData - Data for a payment intent
type PayPalPaymentIntentData struct {
	UserID    int    `bson:"user_id"`
	SubID     string `bson:"sub_id"`
	PlanID    string `bson:"plan_id"`
	PatchLink string `bson:"patch_link"`
	Status    string `bson:"status"`
}

// PayPalPaymentIntent - Intent for a paypal payment
type PayPalPaymentIntent struct {
	IntentID   string                  `bson:"_id"`
	Timestamp  time.Time               `bson:"timestamp"`
	LastUpdate time.Time               `bson:"last_update"`
	Data       PayPalPaymentIntentData `bson:"data"`
}

//
// Referrals
//

// ReferralData - Referal link data
type ReferralData struct {
	ID          string          `bson:"_id"`
	URL         string          `bson:"url"`
	Description string          `bson:"description"`
	Title       string          `bson:"title"`
	Clicks      []ReferralClick `bson:"clicks"`
}

// ReferralClick - Referral link click
type ReferralClick struct {
	URL      string `bson:"url"`
	UserID   int    `bson:"user_id"`
	MetaJSON string `bson:"meta_json"`
}

//
// Players
//

// DBPlayerPofile - Player data db entry struct
type DBPlayerPofile struct {
	ID         int       `json:"player_id" bson:"_id,omitempty"`
	ClanID     int       `json:"clan_id" bson:"clan_id,omitempty"`
	ClanName   string    `json:"clan_name" bson:"clan_name,omitempty"`
	ClanTag    string    `json:"clan_tag" bson:"clan_tag,omitempty"`
	LastBattle time.Time `json:"last_battle_time" bson:"last_battle_time,omitempty"`
	Nickname   string    `json:"nickname" bson:"nickname,omitempty"`
	Realm      string    `json:"realm" bson:"realm,omitempty"`
	CareerWN8  int       `json:"career_wn8" bson:"career_wn8,omitempty"`
}
