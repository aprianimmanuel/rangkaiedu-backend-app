package handlers_test

import (
	"testing"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/testutils"
)

func TestLoginHandler_SocialAuthUserPasswordLogin(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	email := "social@example.com"
	googleID := "google123456"

	// Create social user
	testutils.CreateSocialUser(t, db, email, googleID, "")

	// Test the logic directly by checking if user is social auth
	isSocialAuth := testutils.IsSocialAuthUser(t, db, email)
	if !isSocialAuth {
		t.Error("User should be identified as social auth user")
	}

	// Clean up
	testutils.CleanupUser(t, db, email)
}

func TestLoginHandler_FacebookAuthUserPasswordLogin(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	email := "facebook@example.com"
	facebookID := "facebook123456"

	// Create Facebook user
	testutils.CreateSocialUser(t, db, email, "", facebookID)

	// Test the logic directly by checking if user is social auth
	isSocialAuth := testutils.IsSocialAuthUser(t, db, email)
	if !isSocialAuth {
		t.Error("User should be identified as social auth user")
	}

	// Clean up
	testutils.CleanupUser(t, db, email)
}

func TestLoginHandler_NormalUserPasswordLogin(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	email := "normal@example.com"
	pass := "StrongPass1!"

	// Create normal user with password
	testutils.CreateNormalUser(t, db, email, pass)

	// Test the logic directly by checking if user is social auth
	isSocialAuth := testutils.IsSocialAuthUser(t, db, email)
	if isSocialAuth {
		t.Error("User should NOT be identified as social auth user")
	}

	// Clean up
	testutils.CleanupUser(t, db, email)
}