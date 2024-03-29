package paypal

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cufee/am-api/config"
	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
	"github.com/plutov/paypal/v3"
)

func makePayPalClient(debug bool) (pc *paypal.Client, err error) {
	// // Load cliendID and secretID
	// clientID, err := loadToken("paypal/auth/client_id.dat")
	// if err != nil {
	// 	return nil, err
	// }

	// secretID, err := loadToken("paypal/auth/secret.dat")
	// if err != nil {
	// 	return nil, err
	// }

	// Create a client instance
	// pc, err = paypal.NewClient(clientID, secretID, paypal.APIBaseLive)
	// if err != nil {
	// 	return nil, err
	// }
	// if debug {
	// 	pc.SetLog(os.Stdout) // Set log to terminal stdout
	// }

	// _, err = pc.GetAccessToken()
	// return pc, err
	return nil, nil
}

// HandleNewSub - Handle new subscription request
func HandleNewSub(ctx *fiber.Ctx) error {
	pc, err := makePayPalClient(false)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	discordID, err := strconv.Atoi(ctx.Params("discordID"))
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user data and create new intent
	userData, err := db.UserByDiscordID(discordID)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if user already has an active subscription
	if userData.HasPremiumSub {
		return ctx.JSON(fiber.Map{
			"error": "already subscribed",
		})
	}

	var paymentData db.PayPalPaymentIntentData
	paymentData.UserID = userData.ID
	paymentData.PlanID = monthlyRegularPlan.PlanID

	// Create new subscription struct
	var newSub paypal.SubscriptionBase
	var appCtx paypal.ApplicationContext
	newSub.PlanID = monthlyRegularPlan.PlanID
	// Set redirect URL
	appCtx.ReturnURL = config.PayPalSuccessRedirectURL
	newSub.ApplicationContext = &appCtx

	// Get response
	subRes, err := pc.CreateSubscription(newSub)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check subscription status
	switch subRes.SubscriptionStatus {
	case paypal.SubscriptionStatusApprovalPending:
		// Pending approval
		log.Print("subscription is pending approval")
	case paypal.SubscriptionStatusApproved:
		// Approved by buyer
		log.Print("subscription was approved")
	case paypal.SubscriptionStatusActive:
		// Already exists?
		log.Print("subscription is active")
	default:
		return ctx.JSON(fiber.Map{
			"error": fmt.Sprintf("subscription status: %s", subRes.SubscriptionStatus),
		})
	}

	// Find payment link
	var paymentLink string
	var patchLink string
	for _, link := range subRes.Links {
		if link.Rel == "approve" && strings.Contains(link.Href, "?ba_token") {
			paymentLink = link.Href
		}
		if link.Rel == "edit" && link.Method == "PATCH" {
			patchLink = link.Href
		}
	}

	// Check and return the payment link or error
	if paymentLink == "" {
		return ctx.JSON(fiber.Map{
			"error": "got invalid payment link",
		})
	}

	// Create intent
	paymentData.PatchLink = patchLink
	paymentData.SubID = subRes.ID
	paymentData.Status = string(subRes.SubscriptionStatus)
	_, err = intents.CreatePaymentIntent(paymentData)

	return ctx.JSON(fiber.Map{
		"payment_link": paymentLink,
		"error":        err,
	})
}

// HandlePaymentEvent - handle paypal events
func HandlePaymentEvent(ctx *fiber.Ctx) error {
	var event webhookEvent
	err := ctx.BodyParser(&event)
	if err != nil {
		log.Print(err)
		return nil
	}

	// Select appropriate handler func
	switch event.EventType {
	case "BILLING.SUBSCRIPTION.ACTIVATED":
		subActivated(event)

	case "BILLING.SUBSCRIPTION.RE-ACTIVATED":
		subActivated(event)

	case "BILLING.SUBSCRIPTION.RENEWED":
		subActivated(event)

	case "BILLING.PLAN.DEACTIVATED":
		subDeactivated(event)

	// Need better code, not needed atm
	// case "BILLING.SUBSCRIPTION.CANCELLED":
	// 	subDeactivated(event)

	case "BILLING.SUBSCRIPTION.SUSPENDED":
		subDeactivated(event)

	// No code to handle this properly - might not even need anything here
	// case "BILLING.SUBSCRIPTION.PAYMENT.FAILED":
	// 	// subPaymentFaield(event)
	// 	subDeactivated(event)

	case "BILLING.SUBSCRIPTION.EXPIRED":
		subDeactivated(event)

	default:
		// Undandled event
		log.Print(fmt.Sprintf("Unknown event type: %s", event.EventType))
	}
	return nil
}

func subActivated(data webhookEvent) {
	// Subscription activated or renewed
	// Add premium status to user

	// Get intent
	paymentIntent, err := db.GetPaymentIntentBySubID(data.Resource.ID)
	if err != nil {
		log.Print(fmt.Sprintf("error getting payment intent for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}

	// Get user data
	userData, err := db.UserByDiscordID(paymentIntent.Data.UserID)
	if userData.PremiumExpiration.After(time.Now()) {
		userData.ExcessPremiumMin = int(userData.PremiumExpiration.Sub(time.Now()) / time.Minute)
	}
	userData.HasPremiumSub = true

	// Parse next renewal date
	nextPaymentDate, err := time.Parse(time.RFC3339, data.Resource.BillingInfo.NextBillingTime)
	if err != nil {
		log.Print(fmt.Sprintf("error parsing next billing time for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}

	// Add premium and commit
	userData.PremiumExpiration = nextPaymentDate
	err = db.UpdateUser(userData, false)
	if err != nil {
		log.Print(fmt.Sprintf("error updating user for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}

	// Update intent
	paymentIntent.Data.Status = data.Resource.Status
	err = db.UpdatePaymentIntent(paymentIntent)
	if err != nil {
		log.Print(fmt.Sprintf("error updating payment intent for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}
}

func subDeactivated(data webhookEvent) {
	// Subscription deactivated
	// Suspend premium account

	// Get intent
	paymentIntent, err := db.GetPaymentIntentBySubID(data.Resource.ID)
	if err != nil {
		log.Print(fmt.Sprintf("error getting payment intent for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}

	// Get user data
	userData, err := db.UserByDiscordID(paymentIntent.Data.UserID)
	userData.PremiumExpiration = time.Now().Add(time.Duration(userData.ExcessPremiumMin) * time.Minute)
	userData.ExcessPremiumMin = 1 // Using one bc this field has omitempty
	userData.HasPremiumSub = false

	// Update user
	err = db.UpdateUser(userData, false)
	if err != nil {
		log.Print(fmt.Sprintf("error updating user for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}

	// Update intent
	paymentIntent.Data.Status = data.Resource.Status
	err = db.UpdatePaymentIntent(paymentIntent)
	if err != nil {
		log.Print(fmt.Sprintf("error updating payment intent for sub_id %s: %s", data.Resource.ID, err.Error()))
		return
	}
}

func subPaymentFaield(data webhookEvent) {
	// Subscription payment failed
	// Notify user, set expiration to time.Now() + 25hr
	//
	log.Print(data.EventType, data.Links[1])
}
