package mongodbapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/cufee/am-api/config"
	"github.com/lithammer/shortuuid/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RecordReferalClick - Record referal link usage
func RecordReferalClick(refID string, clickData ReferralClick) error {
	var refData ReferralData

	// Find referral data
	err := referralsCollection.FindOne(ctx, bson.M{"_id": refID}).Decode(&refData)
	if err != nil {
		return err
	}

	// Record click
	refData.Clicks = append(refData.Clicks, clickData)
	_, err = referralsCollection.UpdateOne(ctx, bson.M{"_id": refID}, bson.M{"$set": refData})
	return err
}

// GenerateNewReferalCode - Generate new referal code
func GenerateNewReferalCode(title string, description string) (refData ReferralData, err error) {
	// Generate referral data
	refData.ID = shortuuid.New()
	refData.Description = description
	refData.Title = title
	refData.URL = config.ReferralLinkBase + "/" + refData.ID

	// Add referral to DB
	_, err = referralsCollection.InsertOne(ctx, refData)
	return refData, err
}

// PlayerIDbyNameAndRealm - Get playerId from nickname
func PlayerIDbyNameAndRealm(name string, realm string) (pidInt int, err error) {
	name = fmt.Sprintf("^%v$", name)
	realm = strings.ToUpper(realm)

	var player DBPlayerPofile
	err = playersCollection.FindOne(ctx, bson.M{"nickname": bson.M{"$regex": name, "$options": "i"}, "realm": realm}).Decode(&player)
	if err != nil {
		return pidInt, fmt.Errorf("player not found, make sure this account is tracked by Aftermath")
	}
	return player.ID, err
}

// PlayerExistsByID - Check if a player exists by ID
func PlayerExistsByID(pid int) bool {
	cur, err := playersCollection.Distinct(ctx, "_id", bson.M{"_id": pid})
	if err != nil || len(cur) == 0 {
		return false
	}
	return true
}

// AddPlayerToDB - Add player to DB
func AddPlayerToDB(pid int, realm string) error {
	var player DBPlayerPofile
	player.ID = pid
	player.Realm = realm
	_, err := playersCollection.InsertOne(ctx, player)
	return err
}

// UserByDiscordID - Get existing user by discordID
func UserByDiscordID(did int) (user UserData, err error) {
	err = userDataCollection.FindOne(ctx, bson.M{"_id": did}).Decode(&user)
	return user, err
}

// UserByPlayerID - Get existing user by playerID
func UserByPlayerID(pid int) (user UserData, err error) {
	var opts options.FindOneOptions
	opts.SetSort(bson.M{"verified_expiration": -1})

	err = userDataCollection.FindOne(ctx, bson.M{"verified_id": pid}, &opts).Decode(&user)
	return user, err
}

// UpdateUser - Update existing user with upsert
func UpdateUser(newData UserData, upsert bool) error {
	opts := options.Update().SetUpsert(upsert)
	_, err := userDataCollection.UpdateOne(ctx, bson.M{"_id": newData.ID}, bson.M{"$set": newData}, opts)
	return err
}

// RemoveOldLogins - Remove all existing logins linked to pid
func RemoveOldLogins(pid int) error {
	filter := bson.M{"verified_id": pid}
	cur, err := userDataCollection.Find(ctx, filter)
	if err != nil {
		return err
	}
	for cur.Next(ctx) {
		var u UserData
		err := cur.Decode(&u)
		if err != nil {
			return err
		}

		u.VerifiedID = 0
		u.VerifiedExpiration = time.Time{}

		_, err = userDataCollection.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": u})
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove user by DiscordID/WG_player_id
// TBD

// DeleteIntent - Add new intent to DB
func DeleteIntent(intentID string) {
	intentsCollection.DeleteOne(ctx, bson.M{"_id": intentID})
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

// GetBanData - Get existing ban data
func GetBanData(userID int, days int) (data BanData, err error) {
	// Make a filter by user id and ban timestamp
	timestamp := time.Now().Add(time.Duration(-days) * 24 * time.Hour)
	filter := bson.M{"user_id": userID, "timestamp": bson.M{"$gt": timestamp}}

	err = bansCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}

// GetBansCount - Get existing ban data
func GetBansCount(userID int, days int) (bans int, err error) {
	// Make a filter by user id and ban timestamp
	timestamp := time.Now().Add(time.Duration(-days) * 24 * time.Hour)
	filter := bson.M{"user_id": userID, "timestamp": bson.M{"$gt": timestamp}}

	cur, err := bansCollection.Find(ctx, filter)
	return cur.RemainingBatchLength(), err
}

// BanCheck - Check if a user is banned
func BanCheck(userID int) (data BanData, err error) {
	// Make a filter by user id and ban expiration
	filter := bson.M{"user_id": userID, "expiration": bson.M{"$gt": time.Now()}}

	err = bansCollection.FindOne(ctx, filter).Decode(&data)
	return data, err
}

// AddBanData - Add new ban entry
func AddBanData(data BanData) (err error) {
	// Insert ban object
	_, err = bansCollection.InsertOne(ctx, data)
	return err
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

// NewPaymentIntent - Add new intent to DB
func NewPaymentIntent(intent PayPalPaymentIntent) error {
	_, err := paymentsCollection.InsertOne(ctx, intent)
	return err
}

// GetPaymentIntent - Get intent
func GetPaymentIntent(intentID string) (intent PayPalPaymentIntent, err error) {
	err = paymentsCollection.FindOne(ctx, bson.M{"_id": intentID}).Decode(&intent)
	return intent, err
}

// GetPaymentIntentBySubID - Get intent by Subscription ID
func GetPaymentIntentBySubID(subID string) (intent PayPalPaymentIntent, err error) {
	err = paymentsCollection.FindOne(ctx, bson.M{"data.sub_id": subID}).Decode(&intent)
	return intent, err
}

// UpdatePaymentIntent - Get intent
func UpdatePaymentIntent(intent PayPalPaymentIntent) (err error) {
	opts := options.Update().SetUpsert(false)
	_, err = paymentsCollection.UpdateOne(ctx, bson.M{"_id": intent.IntentID}, bson.M{"$set": bson.M{"data": intent.Data, "last_update": time.Now()}}, opts)
	return err
}
