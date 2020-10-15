package mongodbapi

import "time"

// UserData -
type UserData struct {
	ID                 int       `bson:"_id"`
	Premium            bool      `bson:"premium"`
	PremiumExpiration  time.Time `bson:"premium_expiration"`
	Verified           bool      `bson:"verified"`
	VerifiedExpiration time.Time `bson:"verified_expiration"`
	VerifiedID         int       `bson:"verified_id"`
	DefaultPID         int       `bson:"default_player_id"`
	CustomBgURL        string    `bson:"custom_bg"`
}

// UserDataIntent - Intent to edit User data
type UserDataIntent struct {
	IntentID  string    `bson:"_id"`
	Timestamp time.Time `bson:"timestamp"`
	Data      UserData  `bson:"data"`
}
