package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// GenerateTimestamp returns current timestamp in milliseconds
func GENERATE_TIMESTAMP_UNIX_MILLI() string {
	return fmt.Sprint(time.Now().UnixMilli())
}

// GenerateRequestSignature generates a hash signature for Credit Saison API
func GenerateRequestSignature(secretKey, method, endpoint, timestamp string, queryParams url.Values, payload interface{}) (string, error) {
	// 1. Add signedDate into query params
	if queryParams == nil {
		queryParams = url.Values{}
	}
	queryParams.Set("signedDate", timestamp)

	// 2. Sort query parameters
	var sortedParams []string
	for key := range queryParams {
		for _, value := range queryParams[key] {
			sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", key, value))
		}
	}
	sort.Strings(sortedParams)
	queryString := strings.Join(sortedParams, "&")

	// 3. Body hash (H1)
	var h1 string
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal payload: %w", err)
		}
		hash := sha256.Sum256(bodyBytes)
		h1 = strings.ToUpper(hex.EncodeToString(hash[:]))
	}

	// 4. Canonical request
	canonicalRequest := method + "\n" + endpoint + "\n" + queryString
	if h1 != "" {
		canonicalRequest += "\n" + h1
	}

	// 5. Hash canonical request (H2)
	c1Hash := sha256.Sum256([]byte(canonicalRequest))
	h2 := strings.ToUpper(hex.EncodeToString(c1Hash[:]))

	// 6. HMAC-SHA256 with SecretKey
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(h2))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 7. Base64 encode
	finalSig := base64.StdEncoding.EncodeToString([]byte(signature))

	return finalSig, nil
}

func main() {
	secretKey := "J3K4N6P7Q9SATBVDWEXGZH2J4M"
	apiKey := "3gk7y53DRc8W6ALfOAw865MCa1pzIgYZ31ySYz2o"
	username := "moneycontrol"
	baseURL := "https://api.int2.creditsaison.in"

	endpoint := "/api/v2/partner/redirection/aadhaar"

	payload := map[string]string{
		"partnerLoanId":   "mclRobocop9148598756433314",
		"loanProductCode": "MCL",
		"callbackUrl":     "www.google.com",
	}

	// Generate timestamp separately
	timestamp := GENERATE_TIMESTAMP_UNIX_MILLI()

	// Generate signature
	sig, err := GenerateRequestSignature(secretKey, "POST", endpoint, timestamp, nil, payload)
	if err != nil {
		panic(err)
	}

	// Build full URL with signedDate
	fullURL := fmt.Sprintf("%s%s?signedDate=%s", baseURL, endpoint, timestamp)

	// Marshal payload
	bodyBytes, _ := json.Marshal(payload)

	// Create HTTP request
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		panic(err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("username", username)
	req.Header.Set("signature", sig)

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("🚀 Request URL:", fullURL)
	fmt.Println("🔑 Signature :", sig)
	fmt.Println("🕒 Timestamp :", timestamp)
	fmt.Println("📌 Status   :", resp.Status)
	fmt.Println("📋 Response :", string(respBody))
}
