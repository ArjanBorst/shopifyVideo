package main

type GetMediaFromShopifyRes struct {
	Data struct {
		Node struct {
			ID    string `json:"id"`
			Media struct {
				Edges []struct {
					Node struct {
						MediaContentType string `json:"mediaContentType"`
						Alt              string `json:"alt"`
						Image            *struct {
							URL string `json:"url"`
						} `json:"image"`
						Sources []struct {
							URL      string `json:"url"`
							MimeType string `json:"mimeType"`
							Format   string `json:"format"`
							Height   int    `json:"height"`
							Width    int    `json:"width"`
						} `json:"sources,omitempty"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"media"`
		} `json:"node"`
	} `json:"data"`
	Extensions struct {
		Cost struct {
			RequestedQueryCost float64 `json:"requestedQueryCost"`
			ActualQueryCost    float64 `json:"actualQueryCost"`
			ThrottleStatus     struct {
				MaximumAvailable   float64 `json:"maximumAvailable"`
				CurrentlyAvailable int     `json:"currentlyAvailable"`
				RestoreRate        float64 `json:"restoreRate"`
			} `json:"throttleStatus"`
		} `json:"cost"`
	} `json:"extensions"`
}
