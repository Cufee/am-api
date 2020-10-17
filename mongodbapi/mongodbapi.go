package mongodbapi

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserByDiscordID - Get existing user by discordID
func UserByDiscordID(did int) (user UserData, err error) {
	err = userDataCollection.FindOne(ctx, bson.M{"_id": did}).Decode(&user)
	return user, err
}

// UserByPlayerID - Get existing user by discordID
func UserByPlayerID(did int) (user UserData, err error) {
	err = userDataCollection.FindOne(ctx, bson.M{"_id": did}).Decode(&user)
	return user, err
}

// UpdateUser - Update existing user with upsert
func UpdateUser(newData UserData, upsert bool) error {
	opts := options.Update().SetUpsert(upsert)
	_, err := userDataCollection.UpdateOne(ctx, bson.M{"_id": newData.ID}, bson.M{"$set": newData}, opts)
	return err
}

// Remove user by DiscordID/WG_player_id

// DeleteIntent - Add new intent to DB
func DeleteIntent(intentID string) {
	intentsCollection.DeleteOne(ctx, bson.M{"_id": intentID})
	return
}

// NewUserIntent - Add new intent to DB
func NewUserIntent(intent UserDataIntent) error {
	_, err := intentsCollection.InsertOne(ctx, intent)
	return err
}

// GetUserIntent - Get intent
func GetUserIntent(intentID string) (intent UserDataIntent, err error) {
	err = intentsCollection.FindOne(ctx, bson.M{"_id": intentID}).Decode(&intent)
	return intent, err
}

// GetLogin - Add new intent to DB
func GetLogin(discordID int) int {
	var result UserData
	userDataCollection.FindOne(ctx, bson.M{"_id": discordID}).Decode(&result)

	if time.Now().Before(result.VerifiedExpiration) {
		return result.VerifiedID
	}
	return 0
}

// NewLoginIntent - Add new intent to DB
func NewLoginIntent(intent LoginIntent) error {
	_, err := intentsCollection.InsertOne(ctx, intent)
	return err
}

// GetLoginIntent - Get intent
func GetLoginIntent(intentID string) (intent LoginIntent, err error) {
	err = intentsCollection.FindOne(ctx, bson.M{"_id": intentID}).Decode(&intent)
	return intent, err
}
