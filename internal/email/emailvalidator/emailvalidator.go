package emailvalidator

import (
	"fmt"
	"github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.NewVerifier()
)

// VerifyEmail verifies the syntax of the email address
func VerifyEmail(email string) (*emailverifier.Result, error) {
	// Perform email verification
	ret, err := verifier.Verify(email)
	if err != nil {
		return nil, err
	}

	// Check if email address syntax is valid
	if !ret.Syntax.Valid {
		return nil, fmt.Errorf("email address syntax is invalid")
	}

	// Return the verification result
	return ret, nil
}

// VerifyWithSMTPCheck verifies the email address via SMTP
func VerifyWithSMTPCheck(domain, username string) (*emailverifier.SMTP, error) {
	// Instantiate a verifier with SMTP check enabled
	verifier := emailverifier.NewVerifier().EnableSMTPCheck()

	// Perform email verification via SMTP
	ret, err := verifier.CheckSMTP(domain, username)
	if err != nil {
		return nil, err
	}

	// Return the verification result
	return ret, nil
}

// VerifyWithSOCKS verifies the email address using SOCKS5 proxy
func VerifyWithSOCKS(domain, username string) (*emailverifier.SMTP, error) {
	// Instantiate a verifier with SOCKS proxy enabled
	verifier := emailverifier.NewVerifier().EnableSMTPCheck().Proxy("socks5://user:password@127.0.0.1:1080?timeout=5s")

	// Perform email verification using SOCKS proxy
	ret, err := verifier.CheckSMTP(domain, username)
	if err != nil {
		return nil, err
	}

	// Return the verification result
	return ret, nil
}

// Add more verification methods if needed
