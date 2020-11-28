package paypal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cufee/am-api/config"
	"github.com/cufee/am-api/intents"
	db "github.com/cufee/am-api/mongodbapi"
	"github.com/gofiber/fiber/v2"
	"github.com/plutov/paypal/v3"
)

func makePayPalClient(debug bool) (pc *paypal.Client, err error) {
	// Load cliendID and secretID
	clientID, err := loadToken("paypal/auth/client_id.dat")
	if err != nil {
		return nil, err
	}

	secretID, err := loadToken("paypal/auth/secret.dat")
	if err != nil {
		return nil, err
	}

	// Create a client instance
	pc, err = paypal.NewClient(clientID, secretID, paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}
	if debug {
		pc.SetLog(os.Stdout) // Set log to terminal stdout
	}

	_, err = pc.GetAccessToken()
	return pc, err
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

	var paymentData db.PayPalPaymentIntentData
	paymentData.UserID = userData.ID
	paymentData.PlanID = monthlyRegularPlan.PlanID
	paymentIntent, err := intents.CreatePaymentIntent(paymentData)

	// Create new subscription struct
	var newSub paypal.SubscriptionBase
	var appCtx paypal.ApplicationContext
	newSub.PlanID = monthlyRegularPlan.PlanID
	// Set redirect URL
	appCtx.ReturnURL = config.PayPalSuccessRedirectURL + paymentIntent.IntentID
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

	// Update intent
	paymentIntent.Data.PatchLink = patchLink
	paymentIntent.Data.SubID = subRes.ID
	paymentIntent.Data.Status = subRes.SubscriptionStatus
	err = db.UpdatePaymentIntent(paymentIntent)

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

	case "BILLING.PLAN.DEACTIVATED":
		subDeactivated(event)

	case "BILLING.SUBSCRIPTION.CANCELLED":
		subCancelled(event)

	case "BILLING.SUBSCRIPTION.SUSPENDED":
		subSuspended(event)

	case "BILLING.SUBSCRIPTION.PAYMENT.FAILED":
		subPaymentFaield(event)

	case "BILLING.SUBSCRIPTION.EXPIRED":
		subExpired(event)

	default:
		// Undandled event
		log.Print(fmt.Sprintf("Unknown event type: %s", event.EventType))
	}
	return nil
}

func subActivated(data webhookEvent) {
	// Subscription activated
	// Add premium status to user
	//
	log.Print(data.EventType, data.Links[1])
}

func subDeactivated(data webhookEvent) {
	// Subscription deactivated
	// Suspend premium account
	//
	log.Print(data.EventType, data.Links[1])
}

func subCancelled(data webhookEvent) {
	// Subscription cancelled
	// Suspend premium
	//
	log.Print(data.EventType, data.Links[1])
}

func subSuspended(data webhookEvent) {
	// Subscription suspended
	// Suspend premium, notify owner
	//
	log.Print(data.EventType, data.Links[1])
}

func subPaymentFaield(data webhookEvent) {
	// Subscription payment failed
	// Notify user, set expiration to time.Now() + 25hr
	//
	log.Print(data.EventType, data.Links[1])
}

func subExpired(data webhookEvent) {
	// Subscription expired
	// Send a notification.
	//
	log.Print(data.EventType, data.Links[1])
}
