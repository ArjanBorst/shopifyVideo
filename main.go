package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var googleStorageUrl = "https://shopify-video-production-core-originals.storage.googleapis.com"

var stagedUploadsRes StagedUploadsCreateRes

var productId string
var productAlt string
var graphQlUrl string

var filePath string
var filenname string
var shopifyWebsite string

func main() {
	graphQlUrl = "admin/api/graphql.json"

	shopifyWebsite = "https://www-succubus-com.myshopify.com/"
	GetMediaFromShopify()

	println("Starting REST API endpoints")

	http.HandleFunc("/video/add", HTTPSAddVideo)
	http.ListenAndServe(":8081", nil)
}

func HTTPSAddVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: HTTPSAddVideo")

	shopifyWebsite = r.URL.Query().Get("shopifywebsite")
	productId = r.URL.Query().Get("productid")
	productAlt = r.URL.Query().Get("alt")
	filePath = r.URL.Query().Get("filepath")
	filenname = r.URL.Query().Get("filename")
	shopifyAccessToken = r.URL.Query().Get("shopifyaccesstoken")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	PrepareForUpload()
	UploadToGooleStorage()
	AddVideoToShopifyProduct()
}

func GetMediaFromShopify() {

	//resourceURL := stagedUploadsRes.Data.StagedUploadsCreate.StagedTargets[0].ResourceURL
	//"query": `mutation createProductMedia($id: ID! $media:[CreateMediaInput!]!) { productCreateMedia(productId: $id, media: $media) { media { mediaErrors { code details message } } product { id }mediaUserErrors { code field message } } }`,
	payload := map[string]interface{}{
		"query": `{
			node(id: "gid://shopify/Product/4092299345971") {
			  ...on Product {
				id
				media(first: 10) {
				  edges {
					node {
					  mediaContentType
					  alt
					  ...mediaFieldsByType
					}
				  }
				}
			  }
			}
		  }
		  fragment mediaFieldsByType on Media {
			...on ExternalVideo {
			  id
			  host
				  originUrl
			}
			...on MediaImage {
			  image {
				url
			  }
			}
			...on Model3d {
			  sources {
				url
				mimeType
				format
				filesize
			  }
			}
			...on Video {
			  sources {
				url
				mimeType
				format
				height
				width
			  }
			}
		  }`,
	}

	resp, err := makeGraphQLRequestToShopify(shopifyWebsite+graphQlUrl, payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	getMediaFromShopifyRes := GetMediaFromShopifyRes{}

	err = json.Unmarshal([]byte(resp), &getMediaFromShopifyRes)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}
}

// https://shopify.dev/docs/api/admin-graphql/2023-07/mutations/stagedUploadsCreate
func PrepareForUpload() {
	// Define the GraphQL payload
	size, err := GetFileSize(filePath + filenname)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	payload := map[string]interface{}{
		"query": `mutation stagedUploadsCreate($input: [StagedUploadInput!]!) { stagedUploadsCreate(input: $input) { stagedTargets { url resourceUrl parameters { name value } } } }`,
		"variables": map[string]interface{}{
			"input": []map[string]string{
				{
					"filename": filePath + filenname,
					"mimeType": "video/mp4",
					"fileSize": fmt.Sprintf("%d", size), // convert size (int64) to string
					"resource": "VIDEO",
				},
			},
		},
	}

	resp, err := makeGraphQLRequestToShopify(shopifyWebsite+graphQlUrl, payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = json.Unmarshal([]byte(resp), &stagedUploadsRes)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}
}

func AddVideoToShopifyProduct() {

	resourceURL := stagedUploadsRes.Data.StagedUploadsCreate.StagedTargets[0].ResourceURL

	payload := map[string]interface{}{
		"query": `mutation createProductMedia($id: ID! $media:[CreateMediaInput!]!) { productCreateMedia(productId: $id, media: $media) { media { mediaErrors { code details message } } product { id }mediaUserErrors { code field message } } }`,
		"variables": map[string]interface{}{
			"id": fmt.Sprintf("gid://shopify/Product/%s", productId),
			"media": []map[string]string{
				{
					"originalSource":   resourceURL,
					"alt":              productAlt,
					"mediaContentType": "VIDEO",
				},
			},
		},
	}

	resp, err := makeGraphQLRequestToShopify(shopifyWebsite+graphQlUrl, payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	addMediaToShopifyRes := StagedUploadsCreateRes{}

	err = json.Unmarshal([]byte(resp), &addMediaToShopifyRes)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}
}

func UploadToGooleStorage() {

	stagedTarget := stagedUploadsRes.Data.StagedUploadsCreate.StagedTargets[0]

	googleAccessId := findParameterValue(stagedTarget.Parameters, "GoogleAccessId")
	key := findParameterValue(stagedTarget.Parameters, "key")
	policy := findParameterValue(stagedTarget.Parameters, "policy")
	signature := findParameterValue(stagedTarget.Parameters, "signature")

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	fields := map[string]string{
		"GoogleAccessId": googleAccessId,
		"key":            key,
		"policy":         policy,
		"signature":      signature,
	}

	for key, value := range fields {
		_ = multiPartWriter.WriteField(key, value)
	}

	file, err := os.Open(filePath + filenname)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	fileWriter, err := multiPartWriter.CreateFormFile("file", filenname)
	if err != nil {
		log.Fatalf("Error adding file to form: %v", err)
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalf("Error writing file to form: %v", err)
	}

	// Close the multipart writer to finalize the body
	err = multiPartWriter.Close()
	if err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", googleStorageUrl, &requestBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response: %d", resp.StatusCode)
	}
}
