package intents

import (
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/lithammer/shortuuid/v3"
)

// NewIntentID -
func NewIntentID() string {
	return shortuuid.New()
}

// CreateUserIntent -
func CreateUserIntent(data db.UserData, realm string) (intent db.UserDataIntent, err error) {
	intent.IntentID = shortuuid.New()
	intent.Timestamp = time.Now()
	intent.Realm = realm
	intent.Data = data
	return intent, db.NewUserIntent(intent)
}

// CreateLoginIntent -
func CreateLoginIntent(data db.LoginData) (intentID string, err error) {
	var intent db.LoginIntent
	intent.IntentID = shortuuid.New()
	intent.Timestamp = time.Now()
	intent.LoginData = data
	return intent.IntentID, db.NewLoginIntent(intent)
}

// CreatePaymentIntent -
func CreatePaymentIntent(data db.PayPalPaymentIntentData) (intent db.PayPalPaymentIntent, err error) {
	intent.IntentID = shortuuid.New()
	intent.Timestamp = time.Now()
	intent.Data = data
	return intent, db.NewPaymentIntent(intent)
}
