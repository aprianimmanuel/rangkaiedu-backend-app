// Package sms provides utilities for sending SMS messages, specifically OTP codes,
// for the Rangkai Edu authentication system using the Twilio Go SDK.

package sms

import (
	"fmt"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendOTPSMS sends an SMS containing the OTP code to the specified phone number.
func SendOTPSMS(cfg *config.Config, to, otp string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TWILIOAccountSID,
		Password: cfg.TWILIOAuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(cfg.TWILIOSenderPhone)
	params.SetBody(fmt.Sprintf("Your Rangkai Edu OTP code is: %s. This code expires in 10 minutes. Do not share it with anyone. If you didn't request this, please ignore this message.", otp))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send OTP SMS: %w", err)
	}

	if resp.ErrorCode != nil {
		return fmt.Errorf("Twilio error: %s", *resp.ErrorMessage)
	}

	return nil
}