package auth

import sdk "github.com/friendly-social/golang-sdk"

// LoginMsg signalizes that user logged in with new credentials.
type LoginMsg struct {
	User *sdk.Authorization
}

// WIP
