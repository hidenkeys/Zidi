package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Network mapping
var networkPrefixes = map[string][]string{
	"MTN":     {"0703", "0706", "0803", "0806", "0810", "0813", "0814", "0816", "0903", "0906", "0913", "0916"},
	"Airtel":  {"0701", "0708", "0802", "0808", "0812", "0901", "0902", "0904", "0907", "0912"},
	"Glo":     {"0705", "0805", "0807", "0811", "0815", "0905", "0915"},
	"9mobile": {"0809", "0817", "0818", "0909", "0908"},
}

// Function to get the telecom provider
func GetNetworkProvider(phone string) string {
	if len(phone) < 4 {
		return "Unknown"
	}

	// Extract the first 4 digits
	prefix := phone[:4]

	for network, prefixes := range networkPrefixes {
		for _, p := range prefixes {
			if p == prefix {
				return network
			}
		}
	}

	return "Unknown"
}

// Struct for Bill Payment Request
type BillPaymentRequest struct {
	Country    string `json:"country"`
	Customer   string `json:"customer"`
	Amount     int    `json:"amount"`
	Recurrence string `json:"recurrence"`
	Type       string `json:"type"`
	Reference  string `json:"reference"`
	BillerCode string `json:"biller_code"`
}

// Function to get the biller code for each network
func getBillerCode(network string) (string, error) {
	billerCodes := map[string]string{
		"MTN":     "BIL099", // Replace with actual MTN biller code
		"Airtel":  "BIL100", // Replace with actual Airtel biller code
		"Glo":     "BIL101", // Replace with actual Glo biller code
		"9mobile": "BIL102", // Replace with actual 9mobile biller code
	}

	code, exists := billerCodes[network]
	if !exists {
		return "", fmt.Errorf("unsupported network provider: %s", network)
	}

	return code, nil
}

// Function to send airtime
func SendAirtime(phone, network string, amount int) error {
	billerCode, err := getBillerCode(network)
	if err != nil {
		return err
	}

	url := "https://api.flutterwave.com/v3/bills"
	apiKey := "FLWSECK_TEST-5206ff84c58e85a3717b29f72f507376-X"
	//if apiKey == "" {
	//	return fmt.Errorf("FLW_SECRET_KEY not set in environment variables")
	//}

	// Generate a unique reference using timestamp
	referenceID := fmt.Sprintf("airtime_%s_%d", phone, time.Now().Unix())

	requestBody, _ := json.Marshal(BillPaymentRequest{
		Country:    "NG",
		Customer:   phone,
		Amount:     amount,
		Recurrence: "ONCE",
		Type:       "AIRTIME",
		Reference:  referenceID,
		BillerCode: billerCode,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	// Check if response contains "status" and is "success"
	if status, ok := response["status"].(string); !ok || status != "success" {
		return fmt.Errorf("failed to send airtime: %v", response)
	}

	log.Printf("✅ Airtime sent successfully! Reference: %s\n", referenceID)
	return nil
}

func main() {
	phoneNumber := "08156572209"
	network := GetNetworkProvider(phoneNumber)
	if network == "Unknown" {
		fmt.Println("Invalid phone number. Could not determine network.")
		return
	}

	amount := 500 // Amount to send
	fmt.Printf("Sending %d NGN airtime to %s (%s Network)...\n", amount, phoneNumber, network)

	err := SendAirtime(phoneNumber, network, amount)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
}
