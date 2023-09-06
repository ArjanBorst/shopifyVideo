package main

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StagedUploadsCreateRes struct {
	Data struct {
		StagedUploadsCreate struct {
			StagedTargets []struct {
				URL         string      `json:"url"`
				ResourceURL string      `json:"resourceUrl"`
				Parameters  []Parameter `json:"parameters"`
			} `json:"stagedTargets"`
		} `json:"stagedUploadsCreate"`
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
