package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GetFileSize(filePath string) (int64, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Get the file's information
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func findParameterValue(parameters []Parameter, name string) string {
	for _, param := range parameters {
		if param.Name == name {
			return param.Value
		}
	}
	return ""
}

func makeGraphQLRequestToShopify(url string, payload map[string]interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", shopifyAccessToken) // Replace {access_token} with your actual token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
