package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const googleClientID = "23538631454-pnn2huu7gpl5hb0ejrhmgqqr362icma0.apps.googleusercontent.com"

var googleClientSecret string

const deviceCodeUrl = "https://oauth2.googleapis.com/device/code"
const tokenURL = "https://oauth2.googleapis.com/token"

// without this we can techincally loop forever if "authorization pending" keeps returning
// 120 seconds ÷ 5 second interval = 24 attempts
const maxLoginAttempts = 24

// https://developers.google.com/identity/protocols/oauth2/limited-input-device#step-1:-request-device-and-user-codes
// https://developers.google.com/identity/protocols/oauth2/limited-input-device#step-4:-poll-googles-authorization-server

// auth.go
//   ├── const googleClientID
//   ├── struct DeviceCodeResponse
//   ├── struct TokenResponse
//   ├── func commandLogin()
//   ├── func requestDeviceCode()
//   ├── func pollForToken()
//   └── func saveToken()

// http.PostForm(url, data)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
}

func requestDeviceCode() (*DeviceCodeResponse, error) {

	// perform POST request
	resp, err := http.PostForm(deviceCodeUrl, url.Values{
		"client_id": {googleClientID},
		"scope":     {"email profile"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}

	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceCode DeviceCodeResponse
	if err := Unmarshal(body, &deviceCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &deviceCode, nil

}

func pollForToken(ctx context.Context, deviceCode *DeviceCodeResponse) (*TokenResponse, error) {
	ticker := time.NewTicker(time.Duration(deviceCode.Interval) * time.Second)
	defer ticker.Stop()

	for i := 0; i < maxLoginAttempts; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("login canceled")

		case <-ticker.C: // wait for interval before each attempt
			resp, err := http.PostForm(tokenURL, url.Values{
				"client_id":     {googleClientID},
				"device_code":   {deviceCode.DeviceCode},
				"grant_type":    {"urn:ietf:params:oauth:grant-type:device_code"},
				"client_secret": {googleClientSecret},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to request token: %w", err)
			}

			// read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body.Close()

			var token TokenResponse
			if err := Unmarshal(body, &token); err != nil {
				return nil, fmt.Errorf("failed to unmarshal token response: %w", err)
			}
			if token.Error == "authorization_pending" {
				fmt.Println("Waiting for approval...")
				continue
			}
			if token.Error != "" {
				return nil, fmt.Errorf("token error: %s", token.Error)
			}

			return &token, nil
		}
	}

	return nil, fmt.Errorf("polling timed out")
}

func commandLogin(cfg *Config, args []string) error {
	if cfg.Token != nil {
		fmt.Println("You are already logged in!")
		return nil
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	deviceCodeResponse, err := requestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}
	fmt.Printf("Go to: %s\n", deviceCodeResponse.VerificationUrl)
	fmt.Printf("Enter code: %s\n", deviceCodeResponse.UserCode)
	fmt.Printf("Press Ctrl + C to cancel login\n")
	fmt.Println("Waiting for you to approve...")

	pollTokenReponse, err := pollForToken(ctx, deviceCodeResponse)
	cancel()
	if err != nil {
		fmt.Println("\nLogin cancelled!")
		return nil
	}
	if err := saveToken(pollTokenReponse); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	cfg.Token = pollTokenReponse
	fmt.Println("Successfully logged in!")
	return nil
}

func commandLogOut(cfg *Config, args []string) error {
	// fmt.Printf("Debug - Token is: %v\n", cfg.Token)
	if cfg.Token == nil {
		fmt.Println("You are already logged out!")
		return nil
	}
	var dirPath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	dirPath = filepath.Join(homeDir, ".pokedex")

	tokenPath := filepath.Join(dirPath, "token.json")
	if err := os.Remove(tokenPath); err != nil {
		return fmt.Errorf("failed to delete token file: %w", err)
	}
	cfg.Token = nil
	fmt.Println("Successfully logged out!")
	return nil
}

func saveToken(token *TokenResponse) error {
	var dirPath string

	// get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// build the folder path ~/.pokedex
	dirPath = filepath.Join(homeDir, ".pokedex")

	// create the folder if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// build the full file path ~/.pokedex/token.json
	tokenPath := filepath.Join(dirPath, "token.json")

	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	if err := os.WriteFile(tokenPath, tokenData, 0644); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	fmt.Println("Token Successfully Saved!")
	return nil
}

func loadToken(cfg *Config) error {
	var dirPath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	dirPath = filepath.Join(homeDir, ".pokedex")

	tokenPath := filepath.Join(dirPath, "token.json")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no token file found. Login to create one!")
			return nil
		}
		return fmt.Errorf("failed to read token file: %w", err)
	}

	var token TokenResponse

	if err := Unmarshal(data, &token); err != nil {
		return fmt.Errorf("failed to unmarshal token file: %w", err)
	}

	cfg.Token = &token

	fmt.Println("Successfully loaded token!")
	return nil
}
