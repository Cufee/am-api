package intents

import (
	"time"

	db "github.com/cufee/am-api/mongodbapi"
	"github.com/lithammer/shortuuid/v3"
)

// CreateUserIntent -
func CreateUserIntent(data db.UserData) (intentID string, err error) {
	var intent db.UserDataIntent
	intent.IntentID = shortuuid.New()
	intent.Timestamp = time.Now()
	intent.Data = data
	return intent.IntentID, db.NewIntent(intent)
}

// GetUserIntent -
func GetUserIntent(intentID string) (intent db.UserDataIntent, err error) {
	return db.GetIntent(intentID)
}

// CommitUserIntent -
func CommitUserIntent(intentID string) error {
	intent, err := db.GetIntent(intentID)
	if err != nil {
		return err
	}
	upsert := true
	return db.UpdateUser(intent.Data, upsert)
}
