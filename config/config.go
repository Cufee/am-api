package config

import "os"

// WgAPIAppID - WG Application ID for Wargaming API
var WgAPIAppID string = os.Getenv("WG_APP_ID")

// WgBaseRedirectURL - Redirect base URL for login
var WgBaseRedirectURL string = os.Getenv("WG_BASE_REDIRECT_URL") // "https://thisdomain/users/v1/login/r/"

// MongoURI - URI for connecting to MongoDB
var MongoURI string = os.Getenv("MONGO_CONN_STRING")

// NSFWAPIKey -
var NSFWAPIKey string = os.Getenv("NSFW_API_KEY")

// NSFWAPIURL -
var NSFWAPIURL string = os.Getenv("NSFW_API_URL") // https://api.deepai.org/api/nsfw-detector

// CloudinaryUploadURL -
var CloudinaryUploadURL string = os.Getenv("CLOUDINARY_UPLOAD_URL") // https://api.cloudinary.com/v1_1/vkodev/image/upload

// CloudinaryAPISecret -
var CloudinaryAPISecret string = os.Getenv("CLOUDINARY_API_SECRET")

// CloudinaryAPIKey -
var CloudinaryAPIKey string = os.Getenv("CLOUDINARY_API_KEY")

// PayPalSuccessRedirectURL - Redirect URL for successful payment with PayPal
var PayPalSuccessRedirectURL string = os.Getenv("PAYPAL_SUCCESS_REDIRECT_URL") // https://legacy.amth.one/users/v1/payments/redirect

// ReferralLinkBase - Base url for referral links
var ReferralLinkBase string = os.Getenv("REFERRAL_LINK_BASE") // https://legacy.amth.one/users/v1/r

// AllUsersPremium - Enable all users to be set as premium
var AllUsersPremium bool = true

// OutRPSlimit - Outgoing request limiter for Wargaming API
var OutRPSlimit int = 10

// WGProxyURL - Proxy for outgoing requests
var WGProxyURL string = os.Getenv("WG_PROXY_URL") // https://am-wg-proxy-eu.herokuapp.com/get
