package mongodbapi

import (
	"time"
)

//
// User data
//

// UserData -
type UserData struct {
	ID                 int       `bson:"_id,omitempty"`
	PremiumExpiration  time.Time `bson:"premium_expiration,omitempty"`
	HasPremiumSub      bool      `bson:"has_premium_sub,omitempty"`
	ExcessPremiumMin   int       `bson:"excess_premium_min,omitempty"`
	VerifiedExpiration time.Time `bson:"verified_expiration,omitempty"`
	VerifiedID         int       `bson:"verified_id,omitempty"`
	DefaultPID         int       `bson:"default_player_id,omitempty"`
	CustomBgURL        string    `bson:"custom_bg,omitempty"`
}

// BanData -
type BanData struct {
	UserID     int       `bson:"user_id"`
	Reason     string    `bson:"reason"`
	Notified   bool      `bson:"notified"`
	Timestamp  time.Time `bson:"timestamp"`
	Expiration time.Time `bson:"expiration"`
}

// UserDataIntent - Intent to edit User data
type UserDataIntent struct {
	IntentID  string    `bson:"_id"`
	Timestamp time.Time `bson:"timestamp"`
	Data      UserData  `bson:"data"`
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
