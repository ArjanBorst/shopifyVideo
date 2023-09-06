package main

type addMediaToShopifyRes struct {
	Data       Data       `json:"data"`
	Extensions Extensions `json:"extensions"`
}

type Data struct {
	ProductCreateMedia ProductCreateMedia `json:"productCreateMedia"`
}

type ProductCreateMedia struct {
	Media           []Media     `json:"media"`
	Product         Product     `json:"product"`
	MediaUserErrors []UserError `json:"mediaUserErrors"`
}

type Media struct {
	MediaErrors []Error `json:"mediaErrors"`
}

type Product struct {
	ID string `json:"id"`
}

type Error struct {
}

type UserError struct {
}

type Extensions struct {
	Cost Cost `json:"cost"`
}

type Cost struct {
	RequestedQueryCost float64        `json:"requestedQueryCost"`
	ActualQueryCost    float64        `json:"actualQueryCost"`
	ThrottleStatus     ThrottleStatus `json:"throttleStatus"`
}

type ThrottleStatus struct {
	MaximumAvailable   float64 `json:"maximumAvailable"`
	CurrentlyAvailable int     `json:"currentlyAvailable"`
	RestoreRate        float64 `json:"restoreRate"`
}
